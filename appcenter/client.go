package appcenter

import (
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
	Code       string `json:Code`
	StatusCode int    `json:StatusCode`
	Message    string `json:Message`
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
		errorResponse.Message = string(body)
	}

	return errorResponse
}

func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := &Response{
		Response:    resp,
		StatusError: checkError(resp),
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			body, err := ioutil.ReadAll(response.Body)
			err = json.Unmarshal(body, &v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}
