package appcenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// UploadService definition
type UploadService struct {
	client *Client
}

// UploadRequest wrap the required arguments for the upload specifications
type UploadRequest struct {
	//APIToken  string
	AppName    string
	OwnerName  string
	FilePath   string
	Distribute DistributionPayload
	Option     ReleaseUploadPayload
}

type DistributionPayload struct {
	GroupName string
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
func (s *UploadService) Do(r UploadRequest) (*string, error) {
	err := validateRequest(r)
	if err != nil {
		return nil, err
	}

	// Request Upload "slot"
	fmt.Println("\tBeginning upload", s.client.APIKey)
	uploadResponse, err := s.doUploadRequest(r)
	if err != nil {
		return nil, err
	}

	// Opening file
	reader, err := os.Open(r.FilePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Do the file upload
	err = s.doFileUpload(r, uploadResponse.UploadURL, reader)
	if err != nil {
		return nil, err
	}

	// Commit the release
	return s.releaseCommit(r, uploadResponse)
}

func (s *UploadService) doFileUpload(r UploadRequest, uploadURL string, reader io.Reader) error {
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

	req.Header.Add("Content-Type", multipart.FormDataContentType())

	// Upload body request
	resp, err := s.client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Upload failed : %s %s", err, string(b))
		}
		return fmt.Errorf("Upload failed : %s", string(b))
	}

	fmt.Println("\tUpload completed")
	return nil
}

func (s *UploadService) releaseUploadsRequest(r UploadRequest,
	res *releaseUploadsResponse) (*Response, error) {

	req, err := newRequestWithPayload("POST",
		fmt.Sprintf("%s/apps/%s/%s/release_uploads",
			s.client.BaseURL,
			r.OwnerName,
			r.AppName),
		r.Option)

	req.Header.Add("X-API-Token", s.client.APIKey)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	return s.client.do(req, &res)
}

func (s *UploadService) doUploadRequest(r UploadRequest) (*releaseUploadsResponse, error) {
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

	fmt.Println(color.CyanString("\tUpload Requested"))
	fmt.Println("\t\tUpload ID \t:", result.UploadID)
	fmt.Println("\t\tUpload URL\t:", result.UploadURL)

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
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("ipa", fileName)
	if err != nil {
		return nil, nil, err
	}

	if _, err = io.Copy(fw, fileReader); err != nil {
		return nil, nil, err
	}

	w.Close()

	return w, &b, err
}

func (s *UploadService) releaseCommit(r UploadRequest, u *releaseUploadsResponse) (*string, error) {
	color.Green("\n\tCommitting release")
	// Create request
	req, err := s.createReleaseCommitRequest(r, u)
	if err != nil {
		return nil, err
	}

	// Emit the request
	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, err
	}

	// handle request response
	if resp.StatusCode != http.StatusOK {
		error := checkError(resp)
		return nil, fmt.Errorf("Failed : [%v] %v %v",
			error.Code,
			error.StatusCode,
			error.Message)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Invalid json format for body response %v", string(body))
	}

	response := &patchReleaseUploadResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	fmt.Println("\t\t Release ID  : ", response.ReleaseID)
	fmt.Println("\t\t Release URL : ", response.ReleaseURL)
	color.Green("\tRelease Committed")

	return &response.ReleaseID, nil
}

func (s *UploadService) createReleaseCommitRequest(r UploadRequest, u *releaseUploadsResponse) (*http.Request, error) {
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
