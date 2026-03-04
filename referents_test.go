package genius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_ReferentsService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(&Config{Token: "test", BaseURL: server.URL + "/"})

	tests := []struct {
		client               *Client
		name                 string
		opts                 *ReferentsOptions
		handler              http.HandlerFunc
		wantReferentID       int
		wantReferentPath     string
		wantReferentFeatured *bool
		wantReferentSongID   *int
		wantRangeContent     string
		wantRangeStart       *string
		wantErr              bool
	}{
		{
			name: "Success (Referent as Main Response)",
			opts: &ReferentsOptions{WebPageID: 10611376, TextFormatOptions: TextFormatOptions{TextFormat: FormatPlain}},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/referents" {
					t.Errorf("Expected path /referents, got %s", r.URL.Path)
				}
				if r.URL.RawQuery != "text_format=plain&web_page_id=10611376" {
					t.Errorf("Expected query text_format=plain&web_page_id=10611376, got %s", r.URL.RawQuery)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}
				fmt.Fprint(w, `{
					"response": {
						"referents": [
							{
								"id": 11828417,
								"path": "/11828417/docs.genius.com",
								"featured": false,
								"song_id": null,
								"range": {
									"content": "5848473",
									"start": "/div[2]/div[2]/div[22]"
								}
							}
						]
					}
				}`)
			},
			wantReferentID:       11828417,
			wantReferentPath:     "/11828417/docs.genius.com",
			wantReferentFeatured: func() *bool { b := false; return &b }(),
			wantReferentSongID:   nil,
			wantRangeContent:     "5848473",
			wantRangeStart:       func() *string { s := "/div[2]/div[2]/div[22]"; return &s }(),
			wantErr:              false,
		},
		{
			name: "Success (Referent as Additional Data)",
			opts: &ReferentsOptions{SongID: 245054},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.RawQuery != "song_id=245054" {
					t.Errorf("Expected query song_id=245054, got %s", r.URL.RawQuery)
				}
				fmt.Fprint(w, `{
					"response": {
						"referents": [
							{
								"id": 2359339,
								"path": "/2359339/Vance-joy-riptide",
								"featured": null,
								"song_id": 245054,
								"range": {
									"content": "I wanna be your left-hand man"
								}
							}
						]
					}
				}`)
			},
			wantReferentID:       2359339,
			wantReferentPath:     "/2359339/Vance-joy-riptide",
			wantReferentFeatured: nil,
			wantReferentSongID:   func() *int { i := 245054; return &i }(),
			wantRangeContent:     "I wanna be your left-hand man",
			wantRangeStart:       nil,
			wantErr:              false,
		},
		{
			name: "API Error",
			opts: &ReferentsOptions{SongID: 1},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"meta": {"status": 404, "message": "Not Found"}}`)
			},
			wantErr: true,
		},
		{
			name: "Invalid JSON",
			opts: &ReferentsOptions{SongID: 1},
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
			opts:    &ReferentsOptions{SongID: 1},
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

			referents, _, err := testClient.Referents.Get(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Referents.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(referents) == 0 {
					t.Fatal("Expected referents, got empty slice")
				}
				ref := referents[0]

				if ref.ID != tt.wantReferentID {
					t.Errorf("Expected referent ID %d, got %d", tt.wantReferentID, ref.ID)
				}
				if ref.Path != tt.wantReferentPath {
					t.Errorf("Expected referent path '%s', got %s", tt.wantReferentPath, ref.Path)
				}

				if tt.wantReferentFeatured == nil {
					if ref.Featured != nil {
						t.Errorf("Expected featured nil, got %t", *ref.Featured)
					}
				} else {
					if ref.Featured == nil || *ref.Featured != *tt.wantReferentFeatured {
						t.Errorf("Expected featured '%t', got %v", *tt.wantReferentFeatured, ref.Featured)
					}
				}

				if tt.wantReferentSongID == nil {
					if ref.SongID != nil {
						t.Errorf("Expected song_id nil, got %d", *ref.SongID)
					}
				} else {
					if ref.SongID == nil || *ref.SongID != *tt.wantReferentSongID {
						t.Errorf("Expected song_id '%d', got %v", *tt.wantReferentSongID, ref.SongID)
					}
				}

				if ref.Range.Content != tt.wantRangeContent {
					t.Errorf("Expected range content '%s', got %s", tt.wantRangeContent, ref.Range.Content)
				}

				if tt.wantRangeStart == nil {
					if ref.Range.Start != nil {
						t.Errorf("Expected range start nil, got %s", *ref.Range.Start)
					}
				} else {
					if ref.Range.Start == nil || *ref.Range.Start != *tt.wantRangeStart {
						t.Errorf("Expected range start '%s', got %v", *tt.wantRangeStart, ref.Range.Start)
					}
				}
			}
		})
	}
}
