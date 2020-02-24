package appcenter
import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) ApplyTokenToRequest(h *http.Header) {
	m.Called(h)
}

func (m *MockClient) Do(req *http.Request, v interface{}) (*Response, error) {
	args := m.Called(req, v)
	return args.Get(0).(*Response), args.Error(1)

}

func (m *MockClient) NewServiceRequest(method string, path string, body interface{}) (*http.Request, error) {
	args := m.Called(method, path, body)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockClient) NewRequestWithPayload(method string, url string, body interface{}) (*http.Request, error) {
	args := m.Called(method, url, body)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockClient) RequestContentTypeJSON(h *http.Header) {
	m.Called(h)
}
