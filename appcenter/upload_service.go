package appcenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

type releaseUploadsResponse struct {
	UploadID  string `json:"upload_id"`
	UploadURL string `json:"upload_url"`
}

type ReleaseUploadPayload struct {
	ReleaseID    int    `json:"release_id"`
	BuildVersion string `json:"build_version"`
	BuildNumber  string `json:"build_number"`
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
	fmt.Println("\tRequesting upload")
	var result releaseUploadsResponse
	resp, err := s.releaseUploadsRequest(r, &result)
	if err != nil {
		return nil, err
	}

	// b, err := ioutil.ReadAll(resp.Body)
	// log.Println(string(b))

	if resp.StatusError != nil {
		return nil, fmt.Errorf("%v %s", resp.StatusError.StatusCode,
			resp.StatusError.Code)
	}

	fmt.Println("\tUpload requested :", result.UploadID)

	return &result, nil
}

// Do start the upload request witht the provided parameters
func (s *UploadService) Do(r UploadRequest) error {
	fmt.Println("\tBeginning upload")
	uploadResponse, err := s.doUploadRequest(r)
	if err != nil {
		return err
	}

	// Opening file
	file, err := os.Open(r.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Uploading file
	request, err := getBody(uploadResponse.UploadURL, r.FilePath, file)
	if err != nil {
		return err
	}

	// Upload body request
	resp, err := s.client.client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Panicln(resp.Status)
		return nil
	}

	fmt.Println("\tUpload completed")

	return nil
}

func getBody(url string, fileName string, fileReader io.Reader) (*http.Request, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("ipa", fileName)
		if err != nil {
			return
		}

		if _, err = io.Copy(part, fileReader); err != nil {
			return
		}
	}()

	req, err := http.NewRequest("POST", url, r)
	req.Header.Set("Content-Type", m.FormDataContentType())

	return req, err
}

func (s *UploadService) uploadFile(handle io.Reader, uploadURL string, filePath string) error {
	fmt.Println(("Uploading file...."))
	_, err := s.uploadFileRequest(uploadURL,
		map[string]string{},
		"ipa",
		filePath,
		handle)

	if err != nil {
		return err
	}

	return nil
}

func (s *UploadService) uploadFileRequest(
	uri string,
	params map[string]string,
	paramName string,
	path string,
	handler io.Reader) (*http.Response, error) {
	// file *os.File) (*http.Response, error) {

	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	errs := make(chan error, 1)

	go func() {
		defer w.Close()
		defer m.Close()
		defer close(errs)
		part, err := m.CreateFormFile(paramName, path)
		if err != nil {
			errs <- err
			return
		}

		if _, err = io.Copy(part, handler); err != nil {
			errs <- err
			return
		}
	}()

	return http.Post(uri, m.FormDataContentType(), r)
}
