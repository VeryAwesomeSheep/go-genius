package genius

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
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

	common      Service // reuse a single Client copy for all services
	Annotations *AnnotationsService
	Artists     *ArtistsService
	Songs       *SongsService
	Search      *SearchService
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
	c.Annotations = (*AnnotationsService)(&c.common)
	c.Artists = (*ArtistsService)(&c.common)
	c.Songs = (*SongsService)(&c.common)
	c.Search = (*SearchService)(&c.common)

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

func addOptions(s string, opts any) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	q := u.Query()

	// Helper function for recursive parsing of struct fields
	var parseStruct func(reflect.Value)
	parseStruct = func(v reflect.Value) {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return
			}
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return
		}

		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			// Handle embedded structs (ex. PagingOptions)
			if field.Anonymous {
				parseStruct(value)
				continue
			}

			// Get the "url" tag
			tag := field.Tag.Get("url")
			if tag == "" || tag == "-" {
				continue
			}

			// Parse tag options
			parts := strings.Split(tag, ",")
			key := parts[0]
			omitEmpty := len(parts) > 1 && parts[1] == "omitempty"

			// Check for zero values if omitempty is set
			if omitEmpty && value.IsZero() {
				continue
			}

			// Convert value to string and add to query
			switch value.Kind() {
			case reflect.String:
				q.Set(key, value.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				q.Set(key, strconv.FormatInt(value.Int(), 10))
			case reflect.Bool:
				q.Set(key, strconv.FormatBool(value.Bool()))
			}
		}
	}

	parseStruct(v)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

type PagingOptions struct {
	PerPage int `url:"per_page,omitempty"`
	Page    int `url:"page,omitempty"`
}

type TextFormat string

const (
	FormatDom   TextFormat = "dom"
	FormatPlain TextFormat = "plain"
	FormatHTML  TextFormat = "html"
)

type TextFormatOptions struct {
	TextFormat TextFormat `url:"text_format,omitempty"`
}

type TextBody struct {
	Dom   any    `json:"dom"`   // Populated by default or if ?text_format=dom is used
	Plain string `json:"plain"` // Only populated if ?text_format=plain is used
	HTML  string `json:"html"`  // Only populated if ?text_format=html is used
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
