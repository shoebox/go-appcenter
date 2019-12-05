package appcenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadService definition
type UploadService struct {
	client *Client
}

// UploadRequest wrap the required arguments for the upload specifications
type UploadRequest struct {
	//APIToken  string
	AppName   string
	OwnerName string
	FilePath  string
	Option    ReleaseUploadPayload
}

// ReleaseUploadPayload wrap optional informations about the release
type ReleaseUploadPayload struct {
	ReleaseID    int    `json:"release_id,omitempty"`
	BuildVersion string `json:"build_version,omitempty"`
	BuildNumber  string `json:"build_number,omitempty"`
}

type releaseUploadsResponse struct {
	UploadID  string `json:"upload_id"`
	UploadURL string `json:"upload_url"`
}

type releaseUploadStatus struct {
	Status string `json:"status"`
}

type patchReleaseUploadResponse struct {
	ReleaseID  string `json:"release_id"`
	ReleaseURL string `json:"release_url"`
}

// Do start the upload request witht the provided parameters
func (s *UploadService) Do(r UploadRequest) error {
	err := validateRequest(r)
	if err != nil {
		return err
	}

	// Request Upload "slot"
	fmt.Println("\t", "Beginning upload")
	uploadResponse, err := s.doUploadRequest(r)
	if err != nil {
		return err
	}

	// Opening file
	reader, err := os.Open(r.FilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Do the file upload
	err = s.doFileUpload(r, uploadResponse.UploadURL, reader)
	if err != nil {
		return err
	}

	// Commit the release
	return s.releaseCommit(r, uploadResponse)
}

func (s *UploadService) doFileUpload(r UploadRequest, uploadURL string, reader io.Reader) error {
	fmt.Println("\t", "Starting upload")

	// Create multipart request  body
	multipart, pr, err := getBody(uploadURL, r.FilePath, reader)
	if err != nil {
		return err
	}

	// Create the request
	req, err := http.NewRequest("POST", uploadURL, pr)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", multipart.FormDataContentType())

	// Upload body request
	resp, err := s.client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Upload failed : %s", resp.Status)
	}

	fmt.Println("\t", "Upload completed")
	return nil
}

func (s *UploadService) releaseUploadsRequest(r UploadRequest, res *releaseUploadsResponse) (*Response, error) {

	b, err := json.Marshal(r.Option)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/apps/%s/%s/release_uploads",
			s.client.BaseURL,
			r.OwnerName,
			r.AppName),
		bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	req.Header.Add("X-API-Token", s.client.APIKey)
	req.Header.Add("Content-Type", "application/json")

	return s.client.do(req, &res)
}

func (s *UploadService) doUploadRequest(r UploadRequest) (*releaseUploadsResponse, error) {
	fmt.Println("\t", "Requesting upload")
	var result releaseUploadsResponse
	resp, err := s.releaseUploadsRequest(r, &result)
	if err != nil {
		return nil, err
	}

	if resp.Response.StatusCode < http.StatusOK || resp.Response.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	if resp.StatusError != nil {
		return nil, fmt.Errorf("%v %s", resp.StatusError.StatusCode,
			resp.StatusError.Code)
	}

	fmt.Println("\t", "Upload requested :", result.UploadID)

	return &result, nil
}

func validateRequest(r UploadRequest) error {
	//
	err := validateSource(r)
	if err != nil {
		return err
	}

	//
	err = validateRequestBuildVersion(r)
	if err != nil {
		return err
	}
	return nil
}

func validateSource(r UploadRequest) error {
	_, err := os.Stat(r.FilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("File `%v` does not exsts", r.FilePath)
	}
	return err
}

func validateRequestBuildVersion(r UploadRequest) error {
	ext := filepath.Ext(r.FilePath)

	if r.Option.BuildNumber == "" || r.Option.BuildVersion == "" {
		if ext == ".pkg" || ext == ".dmg" {
			return fmt.Errorf("'--build_version' and '--build_number' parameters "+
				"must be specified to upload file of extension %v", ext)
		}
	}

	if r.Option.BuildVersion == "" {
		if ext == ".zip" || ext == ".msi" {
			return fmt.Errorf("'--build_version' parameter must be "+
				"specified to upload fle of extension '%v'", ext)
		}
	}

	return nil
}

func getBody(url string, fileName string, fileReader io.Reader) (*multipart.Writer, io.Reader, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		var err error

		defer func() {
			if err != nil {
				err := pw.CloseWithError(err)
				if err != nil {
					log.Panicln(err)
				}
			} else {
				pw.Close()
			}
		}()

		partWriter, err := writer.CreateFormFile("ipa", fileName)
		if err != nil {
			return
		}

		_, err = io.Copy(partWriter, fileReader)
		if err != nil {
			return
		}

		err = writer.Close()
	}()

	return writer, pr, nil
}

func (s *UploadService) releaseCommit(r UploadRequest, u *releaseUploadsResponse) error {
	// Create request
	req, err := s.createReleaeCommitRequest(r, u)
	if err != nil {
		return err
	}

	// Emit the request
	resp, err := s.client.client.Do(req)
	if err != nil {
		return err
	}

	// handle request response
	if resp.StatusCode != http.StatusOK {
		error := checkError(resp)
		return fmt.Errorf("Failed : [%v] %v %v", error.Code, error.StatusCode, error.Message)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Invalid json format for body response %v", string(body))
	}

	response := &patchReleaseUploadResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}

	fmt.Println("\t", "Release commited", response.ReleaseID, response.ReleaseURL)

	return err
}

func (s *UploadService) createReleaeCommitRequest(r UploadRequest, u *releaseUploadsResponse) (*http.Request, error) {
	// The json payload
	b, err := json.Marshal(releaseUploadStatus{Status: "committed"})
	if err != nil {
		return nil, err
	}

	// Releasing the upload
	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s/apps/%s/%s/release_uploads/%s",
			s.client.BaseURL,
			r.OwnerName,
			r.AppName,
			u.UploadID),
		bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	req.Header.Add("X-API-Token", s.client.APIKey)
	req.Header.Add("Content-Type", "application/json")

	return req, err
}
