package genius

import (
	"context"
	"fmt"
	"net/http"
)

type ArtistsService Service

// Artist represents data about an artist.
type Artist struct {
	AlternateNames        []string     `json:"alternate_names"` // Available only via ArtistsService
	APIPath               string       `json:"api_path"`
	Description           TextBody     `json:"description"`
	HeaderImageURL        string       `json:"header_image_url"`
	ID                    int          `json:"id"`
	ImageURL              string       `json:"image_url"`
	IsMemeVerified        bool         `json:"is_meme_verified"`
	IsVerified            bool         `json:"is_verified"`
	Name                  string       `json:"name"`
	SocialLinks           *SocialLinks `json:"social_links"` // Available only via ArtistsService
	TranslationArtist     bool         `json:"translation_artist"`
	URL                   string       `json:"url"`
	FollowersCount        int          `json:"followers_count"` // Available only via ArtistsService
	IQ                    *int         `json:"iq"`
	DescriptionAnnotation Referent     `json:"description_annotation"`
	User                  *User        `json:"user"`
}

type SocialLinks struct {
	Twitter   *string `json:"twitter"`
	Facebook  *string `json:"facebook"`
	Instagram *string `json:"instagram"`
}

type ArtistSongsSort string

const (
	SortTitle      ArtistSongsSort = "title"
	SortPopularity ArtistSongsSort = "popularity"
)

type ArtistOptions struct {
	TextFormatOptions
}

type ArtistSongsOptions struct {
	Sort ArtistSongsSort `url:"sort,omitempty"`
	PagingOptions
}

// Get returns data of an artist by its ID.
func (s *ArtistsService) Get(ctx context.Context, id int, opts *ArtistOptions) (*Artist, *http.Response, error) {
	u := fmt.Sprintf("artists/%d", id)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(u)
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

// GetSongs returns paginated data of all songs for an artist.
func (s *ArtistsService) GetSongs(ctx context.Context, id int, opts *ArtistSongsOptions) ([]*SongRelationshipsSong, *http.Response, error) {
	u := fmt.Sprintf("artists/%d/songs", id)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(u)
	if err != nil {
		return nil, nil, err
	}

	var r Response[struct {
		Songs    []*SongRelationshipsSong `json:"songs"`
		NextPage *int                     `json:"next_page"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r.Response.Songs, resp, nil
}
