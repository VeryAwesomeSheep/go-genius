package genius

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	version          = "v0.1.0"
	defaultUserAgent = "go-genius" + "/" + version
	defaultBaseURL   = "https://api.genius.com/"
)

type Client struct {
	http      *http.Client
	baseURL   *url.URL
	userAgent string
	token     string

	common Service // reuse a single Client copy for all services

	Songs *SongsService
}

type Service struct {
	client *Client
}

// Creates a new Genius API client with access token
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	c := &Client{}

	c.http = &http.Client{}
	c.baseURL, _ = url.Parse(defaultBaseURL)
	c.userAgent = defaultUserAgent
	c.token = token

	// Create services
	c.common.client = c
	c.Songs = (*SongsService)(&c.common)

	return c, nil
}

func (c *Client) NewRequest(method, url string, body any) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL)
	}

	u, err := c.baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if body != nil {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}

		bodyReader = buf
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))

	return req, nil
}

type Response[T any] struct {
	Meta struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Response T `json:"response"`
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.http.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF caused by empty response body
		}

		if decErr != nil {
			err = decErr
			return resp, err
		}
	}

	return resp, nil
}

type ErrorResponse struct {
	Response *http.Response `json:"-"`
	Message  string         `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Errorf("Status: %v, Message: %v", r.Response.StatusCode, r.Message).Error()
}

func CheckResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &Response[any]{}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, errorResponse)
	}

	return &ErrorResponse{r, errorResponse.Meta.Message}
}
