package genius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_AnnotationsService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(&Config{Token: "test", BaseURL: server.URL + "/"})

	tests := []struct {
		client                 *Client
		name                   string
		id                     int
		opts                   *AnnotationsOptions
		handler                http.HandlerFunc
		wantAnnID              int
		wantAnnHasVoters       bool
		wantAnnShareURL        string
		wantAnnBodyPlain       string
		wantAnnAuthorName      string
		wantAnnCosignerName    string
		wantAnnVerifierName    string
		wantAnnRejectionReason *string
		wantErr                bool
	}{
		{
			name: "Success (Standard Annotation)",
			id:   123,
			opts: &AnnotationsOptions{TextFormatOptions: TextFormatOptions{TextFormat: FormatPlain}},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/annotations/123" {
					t.Errorf("Expected path /annotations/123, got %s", r.URL.Path)
				}
				if r.URL.RawQuery != "text_format=plain" {
					t.Errorf("Expected query text_format=plain, got %s", r.URL.RawQuery)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}
				fmt.Fprint(w, `{
					"response": {
						"annotation": {
							"id": 123,
							"has_voters": true,
							"share_url": "https://genius.com/123",
							"body": {"plain": "test body"},
							"authors": [{"attribution": 1.0, "user": {"name": "AuthorName"}}],
							"cosigned_by": [],
							"verified_by": null,
							"rejection_comment": null
						},
						"referent": {"id": 456}
					}
				}`)
			},
			wantAnnID:              123,
			wantAnnHasVoters:       true,
			wantAnnShareURL:        "https://genius.com/123",
			wantAnnBodyPlain:       "test body",
			wantAnnAuthorName:      "AuthorName",
			wantAnnCosignerName:    "",
			wantAnnVerifierName:    "",
			wantAnnRejectionReason: nil,
			wantErr:                false,
		},
		{
			name: "Success (Cosigned Annotation)",
			id:   124,
			opts: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{
					"response": {
						"annotation": {
							"id": 124,
							"has_voters": false,
							"share_url": "https://genius.com/124",
							"body": {"plain": "cosigned body"},
							"authors": [],
							"cosigned_by": [{"name": "CosignerName"}],
							"verified_by": null,
							"rejection_comment": "needs work"
						},
						"referent": {"id": 456}
					}
				}`)
			},
			wantAnnID:              124,
			wantAnnHasVoters:       false,
			wantAnnShareURL:        "https://genius.com/124",
			wantAnnBodyPlain:       "cosigned body",
			wantAnnAuthorName:      "",
			wantAnnCosignerName:    "CosignerName",
			wantAnnVerifierName:    "",
			wantAnnRejectionReason: func() *string { s := "needs work"; return &s }(),
			wantErr:                false,
		},
		{
			name: "Success (Verified Annotation)",
			id:   125,
			opts: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{
					"response": {
						"annotation": {
							"id": 125,
							"has_voters": true,
							"share_url": "https://genius.com/125",
							"body": {"plain": "verified body"},
							"authors": [],
							"cosigned_by": [],
							"verified_by": {"name": "VerifierName"},
							"rejection_comment": null
						},
						"referent": {"id": 456}
					}
				}`)
			},
			wantAnnID:              125,
			wantAnnHasVoters:       true,
			wantAnnShareURL:        "https://genius.com/125",
			wantAnnBodyPlain:       "verified body",
			wantAnnAuthorName:      "",
			wantAnnCosignerName:    "",
			wantAnnVerifierName:    "VerifierName",
			wantAnnRejectionReason: nil,
			wantErr:                false,
		},
		{
			name: "API Error",
			id:   123,
			opts: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"meta": {"status": 404, "message": "Not Found"}}`)
			},
			wantErr: true,
		},
		{
			name: "Invalid JSON",
			id:   123,
			opts: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `invalid-json`)
			},
			wantErr: true,
		},
		{
			name: "NewRequest Error",
			client: func() *Client {
				c, _ := NewClient(&Config{Token: "test"})
				c.baseURL, _ = url.Parse("https://api.genius.com/bad") // missing trailing slash
				return c
			}(),
			id:      123,
			opts:    nil,
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

			ann, _, _, err := testClient.Annotations.Get(context.Background(), tt.id, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Annotations.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if ann == nil {
					t.Fatal("Expected annotation, got nil")
				}
				if ann.ID != tt.wantAnnID {
					t.Errorf("Expected annotation ID %d, got %d", tt.wantAnnID, ann.ID)
				}
				if ann.HasVoters != tt.wantAnnHasVoters {
					t.Errorf("Expected annotation has_voters %t, got %t", tt.wantAnnHasVoters, ann.HasVoters)
				}
				if ann.ShareURL != tt.wantAnnShareURL {
					t.Errorf("Expected annotation share_url '%s', got %s", tt.wantAnnShareURL, ann.ShareURL)
				}
				if ann.Body.Plain == nil || *ann.Body.Plain != tt.wantAnnBodyPlain {
					t.Errorf("Expected annotation body plain '%s', got %v", tt.wantAnnBodyPlain, ann.Body.Plain)
				}

				if tt.wantAnnAuthorName != "" {
					if len(ann.Authors) == 0 || ann.Authors[0].User.Name != tt.wantAnnAuthorName {
						t.Errorf("Expected author name '%s', got %v", tt.wantAnnAuthorName, ann.Authors)
					}
				}

				if tt.wantAnnCosignerName != "" {
					if len(ann.CosignedBy) == 0 || ann.CosignedBy[0].Name != tt.wantAnnCosignerName {
						t.Errorf("Expected cosigner name '%s', got %v", tt.wantAnnCosignerName, ann.CosignedBy)
					}
				}

				if tt.wantAnnVerifierName == "" {
					if ann.VerifiedBy != nil {
						t.Errorf("Expected verified_by to be nil, got %+v", ann.VerifiedBy)
					}
				} else {
					if ann.VerifiedBy == nil || ann.VerifiedBy.Name != tt.wantAnnVerifierName {
						t.Errorf("Expected verifier name '%s', got %v", tt.wantAnnVerifierName, ann.VerifiedBy)
					}
				}

				if tt.wantAnnRejectionReason == nil {
					if ann.RejectionComment != nil {
						t.Errorf("Expected rejection_comment nil, got %v", *ann.RejectionComment)
					}
				} else {
					if ann.RejectionComment == nil || *ann.RejectionComment != *tt.wantAnnRejectionReason {
						t.Errorf("Expected rejection_comment '%s', got %v", *tt.wantAnnRejectionReason, ann.RejectionComment)
					}
				}
			}
		})
	}
}
