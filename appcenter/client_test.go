package appcenter

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testClient *Client
	mux        *http.ServeMux
	server     *httptest.Server
)

func openServer(APIKey string) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	url, _ := url.Parse(server.URL)

	setupClient(APIKey, url)
}

func setupClient(APIKey string, url *url.URL) {
	testClient = NewClient(APIKey)
	testClient.BaseURL = url
}

func closeServer() {
	server.Close()
}

type Test struct {
	Test string `json:"test"`
}

func validateBody(t *testing.T, r *http.Request, expected string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}

	if bodyContent := string(b); bodyContent != expected {
		t.Errorf("request Body is %s, expected %s", bodyContent, expected)
	}
}

func validateMethod(t *testing.T, r *http.Request, expected string) {
	if m := r.Method; m != expected {
		t.Errorf("Request method: %v, expected %v", m, expected)
	}
}

func validateHeader(t *testing.T, r *http.Request, header string, expectedValue string) {
	assert.EqualValues(t, r.Header.Get(header), expectedValue)
}

func TestErrorCheckHelper(t *testing.T) {
	r := http.Response{}

	t.Run("Status Code <= 299", func(t *testing.T) {
		r.StatusCode = 201
		assert.Nil(t, checkError(&r))
	})

	t.Run("Status Code >= 300", func(t *testing.T) {
		r.StatusCode = 421

		r.Body = ioutil.NopCloser(strings.NewReader(`{"error":"error message", "errorCode":123}`))
		err := checkError(&r)
		assert.NotNil(t, err)
	})

	t.Run("Status Code >= 300 && Json payload is invalid", func(t *testing.T) {
		r.StatusCode = 421

		r.Body = ioutil.NopCloser(strings.NewReader(`<xml></xml>`))
		err := checkError(&r)

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid JSON body `<xml></xml>`", err.Message)
	})
}

func TestNewRequestWithPayload(t *testing.T) {
	t.Run("A body who cannot be marshalled should throw an error", func(t *testing.T) {
		// Define something who cannot be Marshalled
		body := map[string]interface{}{
			"foo": make(chan int),
		}
		req, err := testClient.NewRequestWithPayload("POST", "http://www.google.com", body)

		assert.Nil(t, req)
		assert.EqualError(t, err, "Error marshalling body "+
			"(Error : json: unsupported type: chan int)")
	})

	t.Run("When request creation fails", func(t *testing.T) {
		req, err := testClient.NewRequestWithPayload("bad method", "", nil)
		assert.Nil(t, req)
		assert.EqualError(t, err, "Error creating request (Error : net/http: "+
			"invalid method \"bad method\")")
	})

	t.Run("Should apply argument", func(t *testing.T) {
		fakeURL := "http://fake-url.com/test"
		req, err := testClient.NewRequestWithPayload("PATCH", fakeURL, []string{})
		assert.Nil(t, err)

		assert.Equal(t, req.Method, "PATCH")
		assert.Equal(t, req.URL.String(), fakeURL)
	})
}

func TestTokenApplication(t *testing.T) {
	// setup
	key := "123-456-XYZ"
	setupClient(key, nil)

	// when:
	head := http.Header{}
	testClient.ApplyTokenToRequest(&head)

	// then:
	assert.EqualValues(t, head.Get("X-API-Token"), key)
}

func TestJsonHeader(t *testing.T) {
	setupClient("", nil)

	// when
	head := http.Header{}
	testClient.RequestContentTypeJSON(&head)

	// then:
	assert.EqualValues(t, head.Get("Content-Type"), "application/json")
}

func TestUnMarshall(t *testing.T) {
	t.Run("Should handle reading errors", func(t *testing.T) {
		// setup:
		setupClient("", nil)
		r := errReader(0)

		// when:
		err := testClient.unmarshal(r, nil)

		// then:
		assert.EqualError(t, err, "Test error")
	})

	t.Run("Should handle JSON decoding errors", func(t *testing.T) {
		// setup:
		setupClient("", nil)
		reader := strings.NewReader("Invalid json")

		// when:
		err := testClient.unmarshal(reader, nil)

		// then:
		assert.EqualError(t, err, "invalid character 'I' looking for beginning of value")
	})
}

type errReader int

func (r errReader) Close() error {
	return nil
}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Test error")
}
