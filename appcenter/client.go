package appcenter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
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

	Config struct {
		OwnerName string
		AppName   string
	}
}

// NewClient create a new instance of the client for the provided APIKey
func NewClient(APIKey string) *Client {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		log.Err(err)
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

// ResponseError errors
type ResponseError struct {
	Message string `json:"message"`
}

// StatusError is the generic reponse body in case of error from AppCenter
type StatusError struct {
	Code       string `json:"Code"`
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"Message"`

	Err *struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	} `json:"error,omitempty"`
}

func (se StatusError) Error() string {
	if se.Err != nil {
		return fmt.Sprintf("Error Code: '%v' Message: '%v' ", se.Err.Code, se.Err.Message)
	}
	return fmt.Sprintf("Error Code: '%v' StatusCode: '%v' Message: '%v'", se.Code, se.StatusCode, se.Message)
}

func checkError(r *http.Response) *StatusError {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &StatusError{}
	if err := json.NewDecoder(r.Body).Decode(errorResponse); err != nil {
		return &StatusError{Message: "Failed to decode response body"}
	}

	return errorResponse
}

func (c *Client) applyTokenToRequest(req *http.Request) *http.Request {
	req.Header.Add("X-API-Token", c.APIKey)
	return req
}

func (c *Client) simpleRequest(ctx context.Context, method string, url string, body interface{}, responseBody interface{}) (*Response, error) {
	var b io.Reader
	if r, ok := body.(io.Reader); ok {
		b = r
	}

	req, err := http.NewRequestWithContext(ctx, method, url, b)
	if err != nil {
		return nil, err
	}

	return c.do(req, &responseBody)
}

func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		//nolint:errcheck
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := &Response{
		Response:    resp,
		StatusError: checkError(resp),
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}

			log.Debug().Str("Body", string(body)).Msg("Response")

			err = json.Unmarshal(body, &v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

// NewAPIRequest is a helper method to do request to AppCenter OpenAPI endpoints
func (c *Client) NewAPIRequest(
	ctx context.Context,
	method string,
	path string,
	requestBody interface{},
	responseBody interface{},
) error {
	body := new(bytes.Buffer)
	if requestBody != nil {
		err := json.NewEncoder(body).Encode(requestBody)
		if err != nil {
			return err
		}
	}

	// Create Request
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s/apps/%s/%s/%s",
			BaseURL,
			c.Config.OwnerName,
			c.Config.AppName,
			path), body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	log.Debug().Str("URL", req.URL.String()).Msg("API Request")
	if err != nil {
		return err
	}

	resp, err := c.do(c.applyTokenToRequest(req), &responseBody)
	if err != nil {
		return err
	}

	if resp.StatusError != nil {
		return resp.StatusError
	}

	return nil
}
