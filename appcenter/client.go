package appcenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	// BaseURL AppCenter base URL for API calls
	BaseURL = "https://api.appcenter.ms/v0.1"
)

var httpClient = &http.Client{}

// Client structure
type Client struct {
	client *http.Client

	BaseURL *url.URL

	APIKey string

	Upload *UploadService

	Distribute *DistributeService
}

// NewClient create a new instance of the client for the provided APIKey
func NewClient(APIKey string) *Client {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		log.Panic(err)
	}

	c := &Client{APIKey: APIKey}
	c.BaseURL = baseURL
	c.client = httpClient
	c.Distribute = &DistributeService{client: c}
	c.Upload = &UploadService{client: c}
	return c
}

// Response of request
type Response struct {
	*http.Response
	*StatusError
}

// AppCenterError errors
type AppCenterError struct {
	Message string `json:"message"`
}

// StatusError is the generic reponse body in case of error from AppCenter
type StatusError struct {
	Code       string `json:"Code"`
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"Message"`
}

func checkError(r *http.Response) *StatusError {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &StatusError{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, errorResponse)

	if err != nil {
		fmt.Println("> failed to parse response body as JSON")
		// failed to unmarhsal API Error, use body as Message
		errorResponse.Message = fmt.Sprintf("Invalid JSON body `%v`", string(body))
	}

	return errorResponse
}

func (c *Client) NewRequestWithPayload(method string,
	url string,
	body interface{}) (*http.Request, error) {

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling body (Error : %v)", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		return nil, fmt.Errorf("Error creating request (Error : %v)", err)
	}

	return req, err
}

func (c *Client) ApplyTokenToRequest(h *http.Header) {
	h.Add("X-API-Token", c.APIKey)
}

func (c *Client) RequestContentTypeJSON(h *http.Header) {
	h.Add("Content-Type", "application/json")
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	// Invoke HTTP client
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Close body
	defer resp.Body.Close()

	// Compose result
	response := &Response{
		Response:    resp,
		StatusError: checkError(resp),
	}

	// Try to unmarshal
	if v != nil {
		c.unmarshal(resp.Body, &v)
	}

	return response, err
}

func (c *Client) unmarshal(reader io.Reader, v *interface{}) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
