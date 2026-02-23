package genius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func Test_DefaultConfig(t *testing.T) {
	t.Setenv("GENIUS_ACCESS_TOKEN", "test-token-from-env")

	config := DefaultConfig()

	if config.BaseURL != defaultBaseURL {
		t.Errorf("DefaultConfig() BaseURL = %v, want %v", config.BaseURL, defaultBaseURL)
	}

	if config.Token != "test-token-from-env" {
		t.Errorf("DefaultConfig() Token = %v, want %v", config.Token, "test-token-from-env")
	}

	if config.UserAgent != defaultUserAgent {
		t.Errorf("DefaultConfig() UserAgent = %v, want %v", config.UserAgent, defaultUserAgent)
	}

	if config.Version != version {
		t.Errorf("DefaultConfig() Version = %v, want %v", config.Version, version)
	}

	if config.HTTP == nil {
		t.Error("DefaultConfig() HTTP client is nil")
	} else if config.HTTP.Timeout != 5*time.Second {
		t.Errorf("DefaultConfig() HTTP timeout = %v, want 5s", config.HTTP.Timeout)
	}
}

func Test_NewClient(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		envToken  string
		wantURL   string
		wantToken string
		wantHTTP  *http.Client
		wantUA    string
		wantErr   bool
	}{
		{
			name:      "Default (nil) config",
			cfg:       nil,
			envToken:  "env-token",
			wantURL:   defaultBaseURL,
			wantToken: "env-token",
			wantUA:    defaultUserAgent,
			wantErr:   false,
		},
		{
			name: "Custom BaseURL, Token, HTTP, UserAgent and Version",
			cfg: &Config{
				BaseURL:   "https://example.com/",
				Token:     "custom-token",
				HTTP:      &http.Client{Timeout: 10 * time.Second},
				UserAgent: "custom-ua",
				Version:   "v2",
			},
			envToken:  "env-token",
			wantURL:   "https://example.com/",
			wantToken: "custom-token",
			wantHTTP:  &http.Client{Timeout: 10 * time.Second},
			wantUA:    "custom-ua/v2",
			wantErr:   false,
		},
		{
			name: "BaseURL missing trailing slash error",
			cfg: &Config{
				BaseURL: "https://example.com",
				Token:   "token",
			},
			envToken: "env-token",
			wantErr:  true,
		},
		{
			name: "Missing token error",
			cfg: &Config{
				Token: "",
			},
			envToken: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GENIUS_ACCESS_TOKEN", tt.envToken)

			c, err := NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if c.baseURL.String() != tt.wantURL {
					t.Errorf("NewClient() baseURL = %v, want %v", c.baseURL, tt.wantURL)
				}
				if c.token != tt.wantToken {
					t.Errorf("NewClient() token = %v, want %v", c.token, tt.wantToken)
				}
				if tt.wantHTTP != nil && c.http.Timeout != tt.wantHTTP.Timeout {
					t.Errorf("NewClient() http timeout = %v, want %v", c.http.Timeout, tt.wantHTTP.Timeout)
				}
				if c.userAgent != tt.wantUA {
					t.Errorf("NewClient() userAgent = %v, want %v", c.userAgent, tt.wantUA)
				}
				if c.Annotations == nil || c.Referents == nil || c.Artists == nil ||
					c.Songs == nil || c.WebPages == nil || c.Search == nil {
					t.Error("NewClient() services were not correctly initialized")
				}
			}
		})
	}
}

func Test_NewRequest(t *testing.T) {
	t.Parallel()

	c, _ := NewClient(&Config{Token: "test-token", UserAgent: "test-ua"})

	tests := []struct {
		name    string
		path    string
		wantURL string
		wantErr bool
	}{
		{
			name:    "Basic GET",
			path:    "songs/123",
			wantURL: "https://api.genius.com/songs/123",
			wantErr: false,
		},
		{
			name:    "Invalid path",
			path:    "%%",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req, err := c.NewRequest(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if req.Method != http.MethodGet {
					t.Errorf("NewRequest() method = %v, want %v", req.Method, http.MethodGet)
				}
				if req.URL.String() != tt.wantURL {
					t.Errorf("NewRequest() URL = %v, want %v", req.URL.String(), tt.wantURL)
				}

				// Check headers
				if req.Header.Get("Accept") != "application/json" {
					t.Errorf("NewRequest() Accept = %v, want application/json", req.Header.Get("Accept"))
				}
				if req.Header.Get("User-Agent") != "test-ua" {
					t.Errorf("NewRequest() User-Agent = %v, want test-ua", req.Header.Get("User-Agent"))
				}
				if req.Header.Get("Authorization") != "Bearer test-token" {
					t.Errorf("NewRequest() Authorization = %v, want Bearer test-token", req.Header.Get("Authorization"))
				}
			}
		})
	}

	t.Run("Invalid baseURL", func(t *testing.T) {
		t.Parallel()
		badURL, _ := url.Parse("https://api.genius.com/api") // No trailing slash
		client := &Client{baseURL: badURL}
		_, err := client.NewRequest("test")
		if err == nil {
			t.Error("NewRequest() expected error for baseURL without trailing slash")
		}
	})
}

func Test_addOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       string
		opts    any
		want    string
		wantErr bool
	}{
		{
			name: "Nil options",
			s:    "http://example.com",
			opts: nil,
			want: "http://example.com",
		},
		{
			name: "PagingOptions",
			s:    "http://example.com",
			opts: &PagingOptions{Page: 1, PerPage: 10},
			want: "http://example.com?page=1&per_page=10",
		},
		{
			name: "TextFormatOptions",
			s:    "http://example.com",
			opts: &TextFormatOptions{TextFormat: FormatPlain},
			want: "http://example.com?text_format=plain",
		},
		{
			name: "TextFormatOptions multiple",
			s:    "http://example.com",
			opts: &TextFormatOptions{TextFormat("dom,html,plain")},
			want: "http://example.com?text_format=dom%2Chtml%2Cplain",
		},
		{
			name: "AnnotationsOptions",
			s:    "http://example.com",
			opts: &AnnotationsOptions{TextFormatOptions: TextFormatOptions{TextFormat: FormatHTML}},
			want: "http://example.com?text_format=html",
		},
		{
			name: "ArtistSongsOptions",
			s:    "http://example.com",
			opts: &ArtistSongsOptions{
				Sort:          SortPopularity,
				PagingOptions: PagingOptions{Page: 2, PerPage: 20},
			},
			want: "http://example.com?page=2&per_page=20&sort=popularity",
		},
		{
			name: "ReferentsOptions (songID)",
			s:    "http://example.com",
			opts: &ReferentsOptions{
				CreatedByID:       123,
				SongID:            456,
				TextFormatOptions: TextFormatOptions{TextFormat: FormatDom},
				PagingOptions:     PagingOptions{Page: 5, PerPage: 10},
			},
			want: "http://example.com?created_by_id=123&page=5&per_page=10&song_id=456&text_format=dom",
		},
		{
			name: "ReferentsOptions (webPageID)",
			s:    "http://example.com",
			opts: &ReferentsOptions{
				CreatedByID:       123,
				WebPageID:         456,
				TextFormatOptions: TextFormatOptions{TextFormat: FormatPlain},
				PagingOptions:     PagingOptions{Page: 1, PerPage: 1},
			},
			want: "http://example.com?created_by_id=123&page=1&per_page=1&web_page_id=456&text_format=plain",
		},
		{
			name: "SearchOptions",
			s:    "http://example.com",
			opts: &SearchOptions{PagingOptions: PagingOptions{Page: 1, PerPage: 10}},
			want: "http://example.com?page=1&per_page=10",
		},
		{
			name: "WebPagesOptions",
			s:    "http://example.com",
			opts: &WebPagesOptions{
				RawAnnotatableURL: "http://raw.com/hi",
				CanonicalURL:      "http://canon.net",
				OgURL:             "https://og.org/yes"},
			want: "http://example.com?canonical_url=http%3A%2F%2Fcanon.net&og_url=https%3A%2F%2Fog.org%2Fyes&raw_annotatable_url=http%3A%2F%2Fraw.com%2Fhi",
		},
		{
			name: "OmitEmpty behavior",
			s:    "http://example.com",
			opts: &PagingOptions{Page: 0, PerPage: 0}, // Opts are 0, should be omitted if omitempty is working
			want: "http://example.com",
		},
		{
			name: "Boolean and pointers",
			s:    "http://example.com",
			opts: &struct {
				Active  bool    `url:"active"`
				Pointer *string `url:"ptr,omitempty"`
				NilPtr  *string `url:"nil_ptr,omitempty"`
				Ignored string  `url:"-"`
				NoTag   string
			}{
				Active:  true,
				Pointer: func() *string { s := "hello"; return &s }(),
				NilPtr:  nil,
				Ignored: "hidden",
				NoTag:   "no-tag",
			},
			want: "http://example.com?active=true&ptr=hello",
		},
		{
			name:    "Invalid base URL",
			s:       "%%",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "Not a struct",
			s:    "http://example.com",
			opts: "not-a-struct",
			want: "http://example.com",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := addOptions(tt.s, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("addOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			gotURL, _ := url.Parse(got)
			wantURL, _ := url.Parse(tt.want)

			if gotURL.Host != wantURL.Host || gotURL.Path != wantURL.Path {
				t.Errorf("addOptions() = %v, want %v", got, tt.want)
			}

			if gotURL.Query().Encode() != wantURL.Query().Encode() {
				t.Errorf("addOptions() query = %v, want %v", gotURL.Query().Encode(), wantURL.Query().Encode())
			}
		})
	}
}

func Test_Do(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	c, _ := NewClient(&Config{Token: "test", BaseURL: server.URL + "/"})

	type testResponse struct {
		Field string `json:"field"`
	}

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		ctx        context.Context
		v          any
		wantResult string
		wantErr    bool
	}{
		{
			name: "Success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"field": "value"}`)
			},
			ctx:        context.Background(),
			v:          &testResponse{},
			wantResult: "value",
			wantErr:    false,
		},
		{
			name: "API Error (404)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"meta": {"message": "not found"}}`)
			},
			ctx:     context.Background(),
			v:       &testResponse{},
			wantErr: true,
		},
		{
			name: "Invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `invalid-json`)
			},
			ctx:     context.Background(),
			v:       &testResponse{},
			wantErr: true,
		},
		{
			name: "Empty body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			ctx:     context.Background(),
			v:       &testResponse{},
			wantErr: false,
		},
		{
			name: "Context cancellation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond)
				fmt.Fprint(w, `{"field": "value"}`)
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			v:       &testResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update the handler for each test case
			// Using a single mux and re-registering the root handler
			server.Config.Handler = http.HandlerFunc(tt.handler)

			req, _ := c.NewRequest("/")
			_, err := c.Do(tt.ctx, req, tt.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.v != nil {
				if res, ok := tt.v.(*testResponse); ok && res.Field != tt.wantResult {
					t.Errorf("Do() result = %v, want %v", res.Field, tt.wantResult)
				}
			}
		})
	}
}

func Test_Error(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusNotFound}
	err := &ErrorResponse{Response: resp, Message: "not found"}
	want := "Status: 404, Message: not found"
	if err.Error() != want {
		t.Errorf("ErrorResponse.Error() = %v, want %v", err.Error(), want)
	}
}
