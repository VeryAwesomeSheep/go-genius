package genius

import (
	"context"
	"fmt"
	"net/http"
)

type ArtistsService Service

type Artist struct {
	AlternateNames []string     `json:"alternate_names"` // Available only via ArtistsService
	APIPath        string       `json:"api_path"`
	ID             int          `json:"id"`
	ImageURL       string       `json:"image_url"`
	IsVerified     bool         `json:"is_verified"`
	Name           string       `json:"name"`
	SocialLinks    *SocialLinks `json:"social_links"` // Available only via ArtistsService
	URL            string       `json:"url"`
	FollowersCount int          `json:"followers_count"` // Available only via ArtistsService
}

type SocialLinks struct {
	Twitter   *string `json:"twitter"`
	Facebook  *string `json:"facebook"`
	Instagram *string `json:"instagram"`
}

func (s *ArtistsService) Get(ctx context.Context, id int) (*Artist, *http.Response, error) {
	u := fmt.Sprintf("artists/%d", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r Response[struct {
		Artist *Artist `json:"artist"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r.Response.Artist, resp, nil
}
