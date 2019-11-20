package appcenter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	//	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	token    string = "AABBCCDDEE"
	uploadID string = "123-456-789"
)

var request = UploadRequest{
	APIToken:  token,
	AppName:   "app-name",
	OwnerName: "owner",
}

var se404 = StatusError{
	Code:       "Not Found",
	StatusCode: 404,
	Message:    "Not found. Context ID: e49d008f-f9c1-4b4e-82b6-e89dc8279d65",
}

type handlerFunc func(t *testing.T, w http.ResponseWriter, r *http.Request)

func handleSuccessFullReleaseUpload(t *testing.T, w http.ResponseWriter, r *http.Request) {
	validateMethod(t, r, http.MethodPost)
	validateHeader(t, r, "X-API-Token", token)

	resp := releaseUploadsResponse{uploadID, "http://" + r.Host + "/upload/file"}
	json, _ := json.Marshal(resp)
	w.Write(json)
}

func handleFailure404(t *testing.T, w http.ResponseWriter, r *http.Request) {
	validateMethod(t, r, http.MethodPost)
	validateHeader(t, r, "X-API-Token", token)

	b, _ := json.Marshal(se404)
	w.WriteHeader(se404.StatusCode)
	w.Write(b)
}

func handlePath(t *testing.T, path string, hf handlerFunc) {
	mux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			hf(t, w, r)
		})
}

func setupServer(t *testing.T, hf handlerFunc) {
	handlePath(t,
		fmt.Sprintf("/apps/%s/%s/release_uploads",
			request.OwnerName,
			request.AppName),
		hf)
}

func TestUploadRequestReleaseSuccess(t *testing.T) {
	// setup:
	openServer()
	setupServer(t, handleSuccessFullReleaseUpload)
	defer closeServer()

	// when:
	var response releaseUploadsResponse
	resp, err := testClient.Upload.releaseUploadsRequest(request, &response)

	// then:
	assert.Nil(t, resp.StatusError)
	assert.Nil(t, err)
}

func TestUploadRequestShouldHandleFailure(t *testing.T) {
	// setup:
	openServer()
	setupServer(t, handleFailure404)
	defer closeServer()

	// when:
	var response releaseUploadsResponse
	resp, _ := testClient.Upload.releaseUploadsRequest(request, &response)

	// then:
	assert.EqualValues(t, resp.StatusError, &se404)
}

func TestUploadDo(t *testing.T) {
	// setup:
	openServer()
	setupServer(t, handleSuccessFullReleaseUpload)
	defer closeServer()

	// when:
	err := testClient.Upload.Do(request)
	log.Println(err)
}

func TestUploadShouldFailInCaseOfErrorDuringUploadRequest(t *testing.T) {
	fakePayload := "fake-data-payload"
	t.Run("Test multipart creation", func(t *testing.T) {
		req, err := getBody("file.ipa", "ipa", strings.NewReader(fakePayload))
		assert.Nil(t, err)

		_, params, err := mime.ParseMediaType(req.Header.Get("Content-type"))

		t.Run("Part should be populated", func(t *testing.T) {
			mr := multipart.NewReader(req.Body, params["boundary"])

			p, err := mr.NextPart()
			assert.Nil(t, err)

			b, err := ioutil.ReadAll(p)
			assert.Nil(t, err)

			assert.EqualValues(t, string(b), fakePayload)
		})
	})
}
