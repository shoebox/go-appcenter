package appcenter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testClient *Client
	path       = "/test/fake/path"
	mux        *http.ServeMux
	server     *httptest.Server
)

func openServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	url, _ := url.Parse(server.URL)

	testClient = NewClient("test-key")
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
