package genius

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	version          = "v0.1.0"
	defaultUserAgent = "go-genius" + "/" + version
	defaultBaseURL   = "https://api.genius.com/"
)

// Config specifies Client configuration.
type Config struct {
	BaseURL   string
	Token     string
	HTTP      *http.Client
	UserAgent string
	Version   string
}

// DefaultConfig creates a default configuration for Client.
func DefaultConfig() *Config {
	config := &Config{
		BaseURL:   defaultBaseURL,
		Token:     os.Getenv("GENIUS_ACCESS_TOKEN"),
		HTTP:      &http.Client{Timeout: 5 * time.Second},
		UserAgent: defaultUserAgent,
		Version:   version,
	}

	return config
}

// Client is a Genius API client that provides bacis for accessing Genius API.
type Client struct {
	http      *http.Client
	baseURL   *url.URL
	userAgent string
	token     string

	common      Service // reuse a single Client copy for all services
	Annotations *AnnotationsService
	Referents   *ReferentsService
	Artists     *ArtistsService
	Songs       *SongsService
	WebPages    *WebPagesService
	Search      *SearchService
}

// Service is the common service struct that holds a reference to the Client.
type Service struct {
	client *Client
}

// NewClient creates a new Genius API client.
func NewClient(cfg *Config) (*Client, error) {
	config := DefaultConfig()

	// Overwrite default config with user defined values
	if cfg != nil {
		if cfg.BaseURL != "" {
			config.BaseURL = cfg.BaseURL
		}
		if cfg.Token != "" {
			config.Token = cfg.Token
		}
		if cfg.HTTP != nil {
			config.HTTP = cfg.HTTP
		}
		if cfg.UserAgent != "" {
			config.UserAgent = cfg.UserAgent
		}
		if cfg.Version != "" {
			config.UserAgent += "/" + cfg.Version
		}
	}

	if !strings.HasSuffix(config.BaseURL, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", config.BaseURL)
	}
	if config.Token == "" {
		return nil, fmt.Errorf("Missing API token")
	}

	c := &Client{}

	c.http = config.HTTP
	c.baseURL, _ = url.Parse(config.BaseURL)
	c.userAgent = config.UserAgent
	c.token = config.Token

	// Create services
	c.common.client = c
	c.Annotations = (*AnnotationsService)(&c.common)
	c.Referents = (*ReferentsService)(&c.common)
	c.Artists = (*ArtistsService)(&c.common)
	c.Songs = (*SongsService)(&c.common)
	c.WebPages = (*WebPagesService)(&c.common)
	c.Search = (*SearchService)(&c.common)

	return c, nil
}

// NewRequest performs basic API request preparation.
func (c *Client) NewRequest(path string) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL)
	}

	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))

	return req, nil
}

// Response is a Genius API response. This wraps the standard http.Response
// and provides the decoded data.
type Response[T any] struct {
	Meta struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Response T `json:"response"`
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields contain "url" tags.
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

			// Handle pointer options
			if value.Kind() == reflect.Ptr {
				if value.IsNil() {
					continue
				}
				value = value.Elem()
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

// PagingOptions specifies the optional parameters for requests that support offset pagination.
type PagingOptions struct {
	PerPage int `url:"per_page,omitempty"`
	Page    int `url:"page,omitempty"`
}

// TextFormat is an enum that represents all possible formats.
type TextFormat string

const (
	FormatDom   TextFormat = "dom"
	FormatPlain TextFormat = "plain"
	FormatHTML  TextFormat = "html"
)

// TextFormatOptions specifies the optional parameters for requests that support variable text format.
type TextFormatOptions struct {
	TextFormat TextFormat `url:"text_format,omitempty"`
}

// TextBody represents the body of a text-based resource, supporting
// multiple formats (DOM, Plain text, HTML).
type TextBody struct {
	Dom   *any    `json:"dom"`   // Populated by default or if ?text_format=dom is used
	Plain *string `json:"plain"` // Only populated if ?text_format=plain is used
	HTML  *string `json:"html"`  // Only populated if ?text_format=html is used
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error occurred.
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

// ErrorResponse reports an error caused by an API request.
type ErrorResponse struct {
	Response *http.Response `json:"-"`
	Message  string         `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Errorf("Status: %v, Message: %v", r.Response.StatusCode, r.Message).Error()
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
func CheckResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &Response[any]{}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, errorResponse)
	}

	return &ErrorResponse{r, errorResponse.Meta.Message}
}
