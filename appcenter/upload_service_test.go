package appcenter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func handleFailure404AndCheckMethod(t *testing.T,
	w http.ResponseWriter,
	r *http.Request,
	method string) {
	validateMethod(t, r, method)
	validateHeader(t, r, "X-API-Token", apiKey)

	b, _ := json.Marshal(se404)
	w.WriteHeader(se404.StatusCode)

	c, err := w.Write(b)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, c, 0)

}
func handleFailure404(t *testing.T,
	w http.ResponseWriter,
	r *http.Request) {
	handleFailure404AndCheckMethod(t, w, r, http.MethodPost)
}

func handlePath(t *testing.T, path string, hf handlerFunc) {
	mux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			hf(t, w, r)
		})
}

func setupServer(t *testing.T, hf handlerFunc) {
	setupServerWithPath(t, hf,
		fmt.Sprintf("/apps/%s/%s/release_uploads",
			request.OwnerName,
			request.AppName))
}

func writeBody(t *testing.T, w http.ResponseWriter, resp interface{}) {
	json, err := json.Marshal(resp)
	assert.Nil(t, err)

	c, err := w.Write(json)
	assert.GreaterOrEqual(t, c, 0)
	assert.Nil(t, err)
}

func setupServerWithPath(t *testing.T, hf handlerFunc, path string) {
	handlePath(t, path, hf)
}

func serverTest(t *testing.T, hf handlerFunc) {
	openServer(apiKey)
	setupServer(t, hf)
}

func TestUploadsRequest(t *testing.T) {
	t.Run("Should handle success", func(t *testing.T) {
		serverTest(t, handleSuccessFullReleaseUpload)
		defer closeServer()

		// when:
		var response releaseUploadsResponse
		resp, err := testClient.Upload.releaseUploadsRequest(request, &response)

		// then:
		assert.Nil(t, resp.StatusError)
		assert.Nil(t, err)
	})

	t.Run("Should handle failure", func(t *testing.T) {
		serverTest(t, handleFailure404)
		defer closeServer()

		// when:
		var response releaseUploadsResponse
		resp, err := testClient.Upload.releaseUploadsRequest(request, &response)
		fmt.Println(err)

		// then:
		assert.EqualValues(t, resp.StatusError, &se404)
	})

	t.Run("Should handle request create error", func(t *testing.T) {
		var response releaseUploadsResponse
		testClient.BaseURL = nil
		resp, err := testClient.Upload.releaseUploadsRequest(request, &response)

		assert.Nil(t, resp)
		assert.EqualError(t, err, "Post %3Cnil%3E/apps/owner/app-name/"+
			"release_uploads: unsupported protocol scheme \"\"")
	})
}

func TestUploadDo(t *testing.T) {
	// setup:
	serverTest(t, handleSuccessFullReleaseUpload)
	defer closeServer()

	// when:
	uploadID, err := testClient.Upload.Do(request)
	assert.Nil(t, uploadID)
	assert.NotNil(t, err)
}

func TestUploadShouldFailInCaseOfErrorDuringUploadRequest(t *testing.T) {
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

/*
func TestShouldHandleErrorAfterUpload(t *testing.T) {
	// setup:
	openServer(apiKey)
	setupServer(t, handleUploadFailure)
	defer closeServer()

	// when:
	t.Run("When doing uploading request", func(t *testing.T) {
		var response releaseUploadsResponse

		// when
		resp, err := testClient.Upload.releaseUploadsRequest(request, &response)

		// then:
		assert.NotNil(t, resp.StatusError)
		assert.Nil(t, err)
	})
}
 */

func TestShouldRequestCommitInProperFormat(t *testing.T) {
	// setup:
	openServer(apiKey)
	setupServer(t, handleCommitRequestSuccess)
	defer closeServer()

	// when:
	resp, err := testClient.Upload.createReleaseCommitRequest(request,
		&releaseUploadsResponse{UploadID: "123-456-789"})

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func handleCommitRequestSuccess(t *testing.T, w http.ResponseWriter, r *http.Request) {
	validateMethod(t, r, http.MethodPost)
	validateHeader(t, r, "X-API-Token", apiKey)

	writeBody(t, w, patchReleaseUploadResponse{
		ReleaseID:  "123",
		ReleaseURL: "http://test.com/test",
	})
}

func TestShouldReleaseTheCommit(t *testing.T) {
	// setup
	openServer(apiKey)
	defer closeServer()

	path := fmt.Sprintf("/apps/%s/%s/release_uploads/%v",
		request.OwnerName,
		request.AppName,
		100)

	cb := func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		validateMethod(t, r, http.MethodPatch)
		validateHeader(t, r, "X-API-Token", apiKey)

		writeBody(t, w, patchReleaseUploadResponse{
			ReleaseID:  "123",
			ReleaseURL: "http://test.com/test",
		})
	}

	setupServerWithPath(t, cb, path)

	t.Run("When committing the release the release_upload endpoint should be invoked",
		func(t *testing.T) {
			resp := releaseUploadsResponse{UploadID: "100"}
			commitID, err := testClient.Upload.releaseCommit(request, &resp)
			assert.Nil(t, err)
			assert.NotNil(t, commitID)
		})
}

func TestErrorShouldBeHandleWhenTryingToReleaseTheCommit(t *testing.T) {
	openServer(apiKey)
	defer closeServer()

	path := fmt.Sprintf("/apps/%s/%s/release_uploads/%v",
		request.OwnerName,
		request.AppName,
		100)

	cb := func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		handleFailure404AndCheckMethod(t, w, r, http.MethodPatch)
	}

	setupServerWithPath(t, cb, path)

	resp := releaseUploadsResponse{UploadID: "100"}
	commitID, err := testClient.Upload.releaseCommit(request, &resp)

	assert.NotNil(t, err)
	assert.Nil(t, commitID)
	assert.EqualError(t, err, "Failed : [Not Found] 404 Not found. Context ID: e49d008f-f9c1-4b4e-82b6-e89dc8279d65")
}

type zeroErrReader struct {
	err error
}

func (r zeroErrReader) Read(p []byte) (int, error) {
	return copy(p, []byte{0}), r.err
}

func TestGetBody(t *testing.T) {

	t.Run("In case of invalid io.reader flow", func(t *testing.T) {
		errorMessage := "Reading error"
		reader := zeroErrReader{err: errors.New(errorMessage)}

		w, r, err := getBody("toto", "file", reader)
		assert.Nil(t, w)
		assert.Nil(t, r)
		assert.EqualError(t, err, errorMessage)
	})

	fakePayload := "fake-data-payload"

	t.Run("Test multipart creation", func(t *testing.T) {
		mw, r, err := getBody("file.ipa", "ipa", strings.NewReader(fakePayload))
		assert.Nil(t, err)

		reader := multipart.NewReader(r, mw.Boundary())

		t.Run("We should expect part 1", func(t *testing.T) {
			part, err := reader.NextPart()
			if part == nil || err != nil {
				t.Error("Expected part1")
				return
			}

			t.Run("And should contain the payload", func(t *testing.T) {
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, part); err != nil {
					t.Errorf("part 1 copy: %v", err)
				}
				assert.Equal(t, buf.String(), fakePayload)
			})
		})

		t.Run("And no more part further", func(t *testing.T) {
			_, err := reader.NextPart()
			assert.Equal(t, err, io.EOF)
		})
	})
}

//
func handleUploadRequestSuccess(t *testing.T, w http.ResponseWriter, r *http.Request) {
	validateMethod(t, r, http.MethodPost)
	validateHeader(t, r, "X-API-Token", apiKey)
}

/*
func TestDoUploadRequest(t *testing.T) {
	t.Run("Should handle release uploadd request errors", func(t *testing.T) {

	})

	t.Run("Should handle success", func(t *testing.T) {
		serverTest(t, handleUploadRequestSuccess)
		defer closeServer()

		// when:
		resp, err := testClient.Upload.doUploadRequest(request)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("Should handle netowking errors", func(t *testing.T) {
		serverTest(t, handleUploadRequestSuccess)
		defer closeServer()

		request := UploadRequest{
			AppName:   "FakeAppName",
			OwnerName: request.OwnerName,
			FilePath:  "test/file.apk",
		}

		resp, err := testClient.Upload.doUploadRequest(request)
		assert.EqualError(t, err, "HTTP request failed: 404 Not Found")
		assert.Nil(t, resp)
	})

	t.Run("Should handle status error", func(t *testing.T) {
		serverTest(t, handleFailure404)
		defer closeServer()

		resp, err := testClient.Upload.doUploadRequest(request)
		assert.EqualError(t, err, "HTTP request failed: 404 Not Found")
		assert.Nil(t, resp)
	})
}
 */
