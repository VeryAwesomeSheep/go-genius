package genius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_ArtistsService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(&Config{Token: "test", BaseURL: server.URL + "/"})

	tests := []struct {
		client               *Client
		name                 string
		id                   int
		opts                 *ArtistOptions
		handler              http.HandlerFunc
		wantArtistID         int
		wantArtistIsVerified bool
		wantArtistName       string
		wantArtistAltNames   []string
		wantArtistIQ         *int
		wantArtistDesc       string
		wantArtistUserName   string
		wantErr              bool
	}{
		{
			name: "Success (Artist User)",
			id:   123,
			opts: &ArtistOptions{TextFormatOptions: TextFormatOptions{TextFormat: FormatPlain}},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/artists/123" {
					t.Errorf("Expected path /artists/123, got %s", r.URL.Path)
				}
				if r.URL.RawQuery != "text_format=plain" {
					t.Errorf("Expected query text_format=plain, got %s", r.URL.RawQuery)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}
				fmt.Fprint(w, `{
					"response": {
						"artist": {
							"id": 123,
							"is_verified": true,
							"name": "Test Artist",
							"alternate_names": ["Tester A", "T-Art"],
							"iq": 450,
							"description": {"plain": "I am a test artist."},
							"user": {
								"id": 999,
								"name": "Test User"
							}
						}
					}
				}`)
			},
			wantArtistID:         123,
			wantArtistIsVerified: true,
			wantArtistName:       "Test Artist",
			wantArtistAltNames:   []string{"Tester A", "T-Art"},
			wantArtistIQ:         func() *int { i := 450; return &i }(),
			wantArtistDesc:       "I am a test artist.",
			wantArtistUserName:   "Test User",
			wantErr:              false,
		},
		{
			name: "Success (Standard Artist without User or IQ)",
			id:   456,
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{
					"response": {
						"artist": {
							"id": 456,
							"is_verified": false,
							"name": "Standard Artist",
							"alternate_names": [],
							"iq": null,
							"description": {"plain": ""},
							"user": null
						}
					}
				}`)
			},
			wantArtistID:         456,
			wantArtistIsVerified: false,
			wantArtistName:       "Standard Artist",
			wantArtistAltNames:   []string{},
			wantArtistIQ:         nil,
			wantArtistDesc:       "",
			wantArtistUserName:   "",
			wantErr:              false,
		},
		{
			name: "API Error",
			id:   123,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"meta": {"status": 404, "message": "Not Found"}}`)
			},
			wantErr: true,
		},
		{
			name: "Invalid JSON",
			id:   123,
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

			artist, _, err := testClient.Artists.Get(context.Background(), tt.id, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Artists.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if artist == nil {
					t.Fatal("Expected artist, got nil")
				}
				if artist.ID != tt.wantArtistID {
					t.Errorf("Expected artist ID %d, got %d", tt.wantArtistID, artist.ID)
				}
				if artist.IsVerified != tt.wantArtistIsVerified {
					t.Errorf("Expected artist is_verified %t, got %t", tt.wantArtistIsVerified, artist.IsVerified)
				}
				if artist.Name != tt.wantArtistName {
					t.Errorf("Expected artist name '%s', got %s", tt.wantArtistName, artist.Name)
				}
				if len(artist.AlternateNames) != len(tt.wantArtistAltNames) {
					t.Errorf("Expected %d alternate names, got %d", len(tt.wantArtistAltNames), len(artist.AlternateNames))
				} else {
					for i, v := range artist.AlternateNames {
						if v != tt.wantArtistAltNames[i] {
							t.Errorf("Expected alternate_name[%d] '%s', got %s", i, tt.wantArtistAltNames[i], v)
						}
					}
				}
				if artist.Description.Plain == nil || *artist.Description.Plain != tt.wantArtistDesc {
					t.Errorf("Expected artist description plain '%s', got %v", tt.wantArtistDesc, artist.Description.Plain)
				}

				if tt.wantArtistIQ == nil {
					if artist.IQ != nil {
						t.Errorf("Expected artist IQ nil, got %d", *artist.IQ)
					}
				} else {
					if artist.IQ == nil || *artist.IQ != *tt.wantArtistIQ {
						t.Errorf("Expected artist IQ '%d', got %v", *tt.wantArtistIQ, artist.IQ)
					}
				}

				if tt.wantArtistUserName == "" {
					if artist.User != nil {
						t.Errorf("Expected user to be nil, got %+v", artist.User)
					}
				} else {
					if artist.User == nil || artist.User.Name != tt.wantArtistUserName {
						t.Errorf("Expected user name '%s', got %v", tt.wantArtistUserName, artist.User)
					}
				}
			}
		})
	}
}

func Test_ArtistsService_GetSongs(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(&Config{Token: "test", BaseURL: server.URL + "/"})

	tests := []struct {
		client         *Client
		name           string
		id             int
		opts           *ArtistSongsOptions
		handler        http.HandlerFunc
		wantSongID     int
		wantSongTitle  string
		wantSongPyongs *int
		wantErr        bool
	}{
		{
			name: "Success",
			id:   123,
			opts: &ArtistSongsOptions{Sort: SortPopularity, PagingOptions: PagingOptions{Page: 2}},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/artists/123/songs" {
					t.Errorf("Expected path /artists/123/songs, got %s", r.URL.Path)
				}
				if r.URL.RawQuery != "page=2&sort=popularity" {
					t.Errorf("Expected query page=2&sort=popularity, got %s", r.URL.RawQuery)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}
				fmt.Fprint(w, `{
					"response": {
						"songs": [
							{
								"id": 1,
								"title": "Song 1",
								"pyongs_count": 15
							}
						],
						"next_page": 3
					}
				}`)
			},
			wantSongID:     1,
			wantSongTitle:  "Song 1",
			wantSongPyongs: func() *int { i := 15; return &i }(),
			wantErr:        false,
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
				c.baseURL, _ = url.Parse("https://api.genius.com/bad")
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

			songs, _, err := testClient.Artists.GetSongs(context.Background(), tt.id, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Artists.GetSongs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(songs) == 0 {
					t.Fatal("Expected songs, got empty slice")
				}
				song := songs[0]
				if song.ID != tt.wantSongID {
					t.Errorf("Expected song ID %d, got %d", tt.wantSongID, song.ID)
				}
				if song.Title != tt.wantSongTitle {
					t.Errorf("Expected song title '%s', got %s", tt.wantSongTitle, song.Title)
				}
				if tt.wantSongPyongs == nil {
					if song.PyongsCount != nil {
						t.Errorf("Expected song pyongs count nil, got %d", *song.PyongsCount)
					}
				} else {
					if song.PyongsCount == nil || *song.PyongsCount != *tt.wantSongPyongs {
						t.Errorf("Expected song pyongs count '%d', got %v", *tt.wantSongPyongs, song.PyongsCount)
					}
				}
			}
		})
	}
}
