package genius

import (
	"context"
	"fmt"
	"net/http"
)

type SongsService Service

type Song struct {
	APIPath            string               `json:"api_path"`
	AppleMusicID       *string              `json:"apple_music_id"`
	ID                 int                  `json:"id"`
	Language           string               `json:"language"`
	PyongsCount        *int                 `json:"pyongs_count"`
	RecordingLocation  *string              `json:"recording_location"`
	ReleaseDate        *string              `json:"release_date"`
	SongArtImageURL    string               `json:"song_art_image_url"`
	Stats              *Stats               `json:"stats"`
	Title              string               `json:"title"`
	URL                string               `json:"url"`
	Album              *Album               `json:"album"`
	CustomPerformances []*CustomPerformance `json:"custom_performances"`
	FeaturedArtists    []*Artist            `json:"featured_artists"`
	Media              []*Media             `json:"media"`
	PrimaryArtists     []*Artist            `json:"primary_artists"`
	ProducerArtists    []*Artist            `json:"producer_artists"`
	SongRelationships  []*SongRelationship  `json:"song_relationships"`
	TranslationSongs   []*TranslationSong   `json:"translation_songs"`
	WriterArtists      []*Artist            `json:"writer_artists"`
}

type Stats struct {
	Contributors int  `json:"contributors"`
	Hot          bool `json:"hot"`
	PageViews    *int `json:"pageviews"`
}

type Album struct {
	APIPath               string    `json:"api_path"`
	CoverArtURL           string    `json:"cover_art_url"`
	FullTitle             string    `json:"full_title"`
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	ReleaseDateForDisplay *string   `json:"release_date_for_display"`
	URL                   string    `json:"url"`
	PrimaryArtists        []*Artist `json:"primary_artists"`
}

type CustomPerformance struct {
	Label   string    `json:"label"`
	Artists []*Artist `json:"artists"`
}

type Media struct {
	Provider    string  `json:"provider"`
	Start       *int    `json:"start"`
	Type        string  `json:"type"`
	URL         string  `json:"url"`
	NativeURI   *string `json:"native_uri"`
	Attribution *string `json:"attribution"`
}

type SongRelationship struct {
	RelationshipType string                   `json:"relationship_type"`
	Type             string                   `json:"type"`
	URL              *string                  `json:"url"`
	Songs            []*SongRelationshipsSong `json:"songs"`
}

type SongRelationshipsSong struct {
	APIPath               string                 `json:"api_path"`
	ID                    int                    `json:"id"`
	PrimaryArtistNames    string                 `json:"primary_artist_names"`
	PyongsCount           *int                   `json:"pyongs_count"`
	RelationshipsIndexURL string                 `json:"relationships_index_url"`
	ReleaseDateComponents *ReleaseDateComponents `json:"release_date_components"`
	SongArtImageURL       string                 `json:"song_art_image_url"`
	Title                 string                 `json:"title"`
	URL                   *string                `json:"url"`
	FeaturedArtists       []*Artist              `json:"featured_artists"`
	PrimaryArtists        []*Artist              `json:"primary_artists"`
}

type ReleaseDateComponents struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type TranslationSong struct {
	APIPath  string `json:"api_path"`
	ID       int    `json:"id"`
	Language string `json:"language"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

func (s *SongsService) Get(ctx context.Context, id int) (*Song, *http.Response, error) {
	u := fmt.Sprintf("songs/%d", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r Response[struct {
		Song *Song `json:"song"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r.Response.Song, resp, nil
}
