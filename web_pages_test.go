package genius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_WebPagesService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(&Config{Token: "test", BaseURL: server.URL + "/"})

	tests := []struct {
		client               *Client
		name                 string
		opts                 *WebPagesOptions
		handler              http.HandlerFunc
		wantWebPageID        *int
		wantWebPageDomain    string
		wantWebPageAPIPath   *string
		wantAnnotationCount  int
		wantErr              bool
	}{
		{
			name: "Success (Existing Web Page)",
			opts: &WebPagesOptions{RawAnnotatableURL: "https://docs.genius.com"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/web_pages/lookup" {
					t.Errorf("Expected path /web_pages/lookup, got %s", r.URL.Path)
				}
				if r.URL.RawQuery != "raw_annotatable_url=https%3A%2F%2Fdocs.genius.com" {
					t.Errorf("Expected query raw_annotatable_url=https%%3A%%2F%%2Fdocs.genius.com, got %s", r.URL.RawQuery)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}
				fmt.Fprint(w, `{
					"response": {
						"web_page": {
							"api_path": "/web_pages/10347",
							"domain": "docs.genius.com",
							"id": 10347,
							"normalized_url": "//docs.genius.com",
							"share_url": "http://genius.it/docs.genius.com",
							"title": "Genius API",
							"url": "https://genius.com/docs.genius.com",
							"annotation_count": 26
						}
					}
				}`)
			},
			wantWebPageID:        func() *int { i := 10347; return &i }(),
			wantWebPageDomain:    "docs.genius.com",
			wantWebPageAPIPath:   func() *string { s := "/web_pages/10347"; return &s }(),
			wantAnnotationCount:  26,
			wantErr:              false,
		},
		{
			name: "Success (Non-existent Web Page)",
			opts: &WebPagesOptions{RawAnnotatableURL: "https://youtube.com"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{
					"response": {
						"web_page": {
							"api_path": null,
							"domain": "youtube.com",
							"id": null,
							"normalized_url": "//youtube.com",
							"share_url": "http://genius.it/youtube.com",
							"title": "youtube.com",
							"url": "https://genius.com/youtube.com",
							"annotation_count": 0
						}
					}
				}`)
			},
			wantWebPageID:        nil,
			wantWebPageDomain:    "youtube.com",
			wantWebPageAPIPath:   nil,
			wantAnnotationCount:  0,
			wantErr:              false,
		},
		{
			name: "API Error",
			opts: &WebPagesOptions{RawAnnotatableURL: "https://error.com"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"meta": {"status": 404, "message": "Not Found"}}`)
			},
			wantErr: true,
		},
		{
			name: "Invalid JSON",
			opts: &WebPagesOptions{RawAnnotatableURL: "https://invalid.com"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `invalid-json`)
			},
			wantErr: true,
		},
		{
			name: "NewRequest Error",
			client: func() *Client {
				c, _ := NewClient(&Config{Token: "test"})
				c.baseURL, _ = url.Parse("https://api.genius.com/bad")
				return c
			}(),
			opts:    &WebPagesOptions{RawAnnotatableURL: "https://docs.genius.com"},
			handler: func(w http.ResponseWriter, r *http.Request) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.Config.Handler = http.HandlerFunc(tt.handler)

			testClient := client
			if tt.client != nil {
				testClient = tt.client
			}

			webPage, _, err := testClient.WebPages.Get(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WebPages.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if webPage == nil {
					t.Fatal("Expected web_page, got nil")
				}

				if tt.wantWebPageID == nil {
					if webPage.ID != nil {
						t.Errorf("Expected web_page ID nil, got %d", *webPage.ID)
					}
				} else {
					if webPage.ID == nil || *webPage.ID != *tt.wantWebPageID {
						t.Errorf("Expected web_page ID '%d', got %v", *tt.wantWebPageID, webPage.ID)
					}
				}

				if webPage.Domain != tt.wantWebPageDomain {
					t.Errorf("Expected web_page domain '%s', got %s", tt.wantWebPageDomain, webPage.Domain)
				}

				if tt.wantWebPageAPIPath == nil {
					if webPage.APIPath != nil {
						t.Errorf("Expected web_page api_path nil, got %s", *webPage.APIPath)
					}
				} else {
					if webPage.APIPath == nil || *webPage.APIPath != *tt.wantWebPageAPIPath {
						t.Errorf("Expected web_page api_path '%s', got %v", *tt.wantWebPageAPIPath, webPage.APIPath)
					}
				}

				if webPage.AnnotationCount != tt.wantAnnotationCount {
					t.Errorf("Expected web_page annotation_count %d, got %d", tt.wantAnnotationCount, webPage.AnnotationCount)
				}
			}
		})
	}
}
