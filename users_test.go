package genius

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_UsersService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, _ := NewClient(&Config{Token: "test", BaseURL: server.URL + "/"})

	tests := []struct {
		client              *Client
		name                string
		id                  int
		opts                *UserOptions
		handler             http.HandlerFunc
		wantUserID          int
		wantUserName        string
		wantUserAboutMe     string
		wantUserAvatarURL   string
		wantUserCustomImage *string
		wantUserRoles       []string
		wantUserArtistName  string
		wantErr             bool
	}{
		{
			name: "Success (Artist User)",
			id:   1234,
			opts: &UserOptions{TextFormatOptions: TextFormatOptions{TextFormat: FormatPlain}},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/users/1234" {
					t.Errorf("Expected path /users/1234, got %s", r.URL.Path)
				}
				if r.URL.RawQuery != "text_format=plain" {
					t.Errorf("Expected query text_format=plain, got %s", r.URL.RawQuery)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}
				fmt.Fprint(w, `{
					"response": {
						"user": {
							"id": 1234,
							"name": "Test User",
							"role_for_display": "verified_artist",
							"roles_for_display": ["verified_artist", "contributor"],
							"about_me": {
								"plain": "I am a test user."
							},
							"avatar": {
								"tiny": {"url": "http://tiny.url", "bounding_box": {"width": 16, "height": 16}},
								"thumb": {"url": "http://thumb.url", "bounding_box": {"width": 32, "height": 32}},
								"small": {"url": "http://small.url", "bounding_box": {"width": 100, "height": 100}},
								"medium": {"url": "http://medium.url", "bounding_box": {"width": 300, "height": 300}}
							},
							"custom_header_image_url": "http://custom.url",
							"artist": {
								"name": "The Artist"
							}
						}
					}
				}`)
			},
			wantUserID:          1234,
			wantUserName:        "Test User",
			wantUserAboutMe:     "I am a test user.",
			wantUserAvatarURL:   "http://medium.url",
			wantUserCustomImage: func() *string { s := "http://custom.url"; return &s }(),
			wantUserRoles:       []string{"verified_artist", "contributor"},
			wantUserArtistName:  "The Artist",
			wantErr:             false,
		},
		{
			name: "Success (Standard User)",
			id:   5678,
			opts: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{
					"response": {
						"user": {
							"id": 5678,
							"name": "Standard User",
							"role_for_display": "contributor",
							"roles_for_display": ["contributor"],
							"about_me": {
								"plain": ""
							},
							"avatar": {
								"tiny": {"url": "http://tiny.url", "bounding_box": {"width": 16, "height": 16}},
								"thumb": {"url": "http://thumb.url", "bounding_box": {"width": 32, "height": 32}},
								"small": {"url": "http://small.url", "bounding_box": {"width": 100, "height": 100}},
								"medium": {"url": "http://medium.url", "bounding_box": {"width": 300, "height": 300}}
							},
							"custom_header_image_url": null
						}
					}
				}`)
			},
			wantUserID:          5678,
			wantUserName:        "Standard User",
			wantUserAboutMe:     "",
			wantUserAvatarURL:   "http://medium.url",
			wantUserCustomImage: nil,
			wantUserRoles:       []string{"contributor"},
			wantUserArtistName:  "",
			wantErr:             false,
		},
		{
			name: "API Error",
			id:   1234,
			opts: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"meta": {"status": 404, "message": "Not Found"}}`)
			},
			wantErr: true,
		},
		{
			name: "Invalid JSON",
			id:   1234,
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
			id:      1234,
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

			user, _, err := testClient.Users.Get(context.Background(), tt.id, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Fatal("Expected user, got nil")
				}
				if user.ID != tt.wantUserID {
					t.Errorf("Expected user ID %d, got %d", tt.wantUserID, user.ID)
				}
				if user.Name != tt.wantUserName {
					t.Errorf("Expected user name '%s', got %s", tt.wantUserName, user.Name)
				}
				if len(user.RolesForDisplay) != len(tt.wantUserRoles) {
					t.Errorf("Expected %d roles, got %d", len(tt.wantUserRoles), len(user.RolesForDisplay))
				} else {
					for i, v := range user.RolesForDisplay {
						if v != tt.wantUserRoles[i] {
							t.Errorf("Expected role[%d] '%s', got %s", i, tt.wantUserRoles[i], v)
						}
					}
				}
				if user.AboutMe.Plain == nil || *user.AboutMe.Plain != tt.wantUserAboutMe {
					t.Errorf("Expected user about_me plain '%s', got %v", tt.wantUserAboutMe, user.AboutMe.Plain)
				}
				if user.Avatar.Medium.URL != tt.wantUserAvatarURL {
					t.Errorf("Expected user avatar medium url '%s', got %s", tt.wantUserAvatarURL, user.Avatar.Medium.URL)
				}

				if tt.wantUserCustomImage == nil {
					if user.CustomHeaderImageURL != nil {
						t.Errorf("Expected custom_header_image_url nil, got %v", *user.CustomHeaderImageURL)
					}
				} else {
					if user.CustomHeaderImageURL == nil || *user.CustomHeaderImageURL != *tt.wantUserCustomImage {
						t.Errorf("Expected custom_header_image_url '%s', got %v", *tt.wantUserCustomImage, user.CustomHeaderImageURL)
					}
				}

				if tt.wantUserArtistName == "" {
					if user.Artist != nil {
						t.Errorf("Expected artist to be nil, got %+v", user.Artist)
					}
				} else {
					if user.Artist == nil || user.Artist.Name != tt.wantUserArtistName {
						t.Errorf("Expected artist name '%s', got %v", tt.wantUserArtistName, user.Artist)
					}
				}
			}
		})
	}
}
