package appcenter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	apiKey   string = "AABBCCDDEE"
	uploadID string = "123-456-789"
)

var request = UploadRequest{
	AppName:   "app-name",
	OwnerName: "owner",
	FilePath:  "test",
}

var se404 = StatusError{
	Code:       "Not Found",
	StatusCode: 404,
	Message:    "Not found. Context ID: e49d008f-f9c1-4b4e-82b6-e89dc8279d65",
}

type handlerFunc func(t *testing.T, w http.ResponseWriter, r *http.Request)

func handleSuccessFullReleaseUpload(t *testing.T, w http.ResponseWriter, r *http.Request) {
	validateMethod(t, r, http.MethodPost)
	validateHeader(t, r, "X-API-Token", apiKey)

	resp := releaseUploadsResponse{uploadID, "http://" + r.Host + "/upload/file"}
	json, err := json.Marshal(resp)
	assert.Nil(t, err)

	c, err := w.Write(json)
	assert.GreaterOrEqual(t, c, 0)
	assert.Nil(t, err)
}

func handleFailure404(t *testing.T, w http.ResponseWriter, r *http.Request) {
	validateMethod(t, r, http.MethodPost)
	validateHeader(t, r, "X-API-Token", apiKey)

	b, _ := json.Marshal(se404)
	w.WriteHeader(se404.StatusCode)

	c, err := w.Write(b)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, c, 0)
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
	openServer(apiKey)
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
	openServer(apiKey)
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
	openServer(apiKey)
	setupServer(t, handleSuccessFullReleaseUpload)
	defer closeServer()

	// when:
	err := testClient.Upload.Do(request)
	fmt.Println(err)
}

func TestUploadShouldFailInCaseOfErrorDuringUploadRequest(t *testing.T) {
	fakePayload := "fake-data-payload"
	t.Run("Test multipart creation", func(t *testing.T) {
		req, err := getBody("file.ipa", "ipa", strings.NewReader(fakePayload))
		assert.Nil(t, err)

		_, params, err := mime.ParseMediaType(req.Header.Get("Content-type"))
		assert.NoError(t, err)

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

func TestValdationBuildVersionArgument(t *testing.T) {
	r := UploadRequest{}

	t.Run("When the build verson number is missing", func(t *testing.T) {
		testCases := []struct {
			ext          string
			buildVersion string
			buildNumber  string
			err          bool
		}{
			{"zip", "", "", true},
			{"zip", "1.2.3", "", false},
			{"zip", "", "1", true},

			{"msi", "", "", true},
			{"msi", "1.2.3", "", false},
			{"msi", "", "1", true},

			{"apk", "", "", false},
			{"apk", "1.2.3", "", false},
			{"apk", "1.2.3", "1", false},
			{"apk", "", "1", false},

			{"pkg", "", "", true},
			{"pkg", "1.2.3", "", true},
			{"pkg", "", "1", true},
			{"pkg", "1.2.3", "1", false},

			{"dmg", "", "", true},
			{"dmg", "1.2.3", "", true},
			{"dmg", "", "1", true},
			{"dmg", "1.2.3", "1", false},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("For ext: %v BuildVersion: %v BuildNumber: %v",
				tc.buildVersion, tc.buildNumber, tc.ext), func(t *testing.T) {

				r.Option.BuildVersion = tc.buildVersion
				r.Option.BuildNumber = tc.buildNumber

				r.FilePath = fmt.Sprintf("toto.%v", tc.ext)
				err := validateRequestBuildVersion(r)
				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func handleUploadFailure(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Run("Method should be POST", func(t *testing.T) {
		validateMethod(t, r, http.MethodPost)
	})
	t.Run("API Should be present", func(t *testing.T) {
		validateHeader(t, r, "X-API-Token", apiKey)
	})

	b, _ := json.Marshal(se404)

	w.WriteHeader(http.StatusNotFound)

	c, err := w.Write(b)
	assert.GreaterOrEqual(t, c, 1)
	assert.NoError(t, err)
}

func TestShouldHandleErrorAfterUpload(t *testing.T) {
	// setup:
	openServer(apiKey)
	setupServer(t, handleUploadFailure)
	defer closeServer()

	// when:
	t.Run("When doing uploading request", func(t *testing.T) {
		var response releaseUploadsResponse
		resp, err := testClient.Upload.releaseUploadsRequest(request, &response)

		t.Run("Should report error", func(t *testing.T) {
			assert.NotNil(t, resp.StatusError)
			assert.Nil(t, err)
		})
	})
}
