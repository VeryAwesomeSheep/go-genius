package genius

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type SearchService Service

type Hit struct {
	Type   string                `json:"type"`
	Result SongRelationshipsSong `json:"result"`
}

type SearchOptions struct {
	PagingOptions
}

func (s *SearchService) Get(ctx context.Context, query string, opts *SearchOptions) ([]*SongRelationshipsSong, *http.Response, error) {
	u := fmt.Sprintf("search?q=%s", url.QueryEscape(query))

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r Response[struct {
		Hits []*Hit `json:"hits"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	// Official API always returns songs from search, so to simplify
	// struct access unwrap the response into the array of songs
	var songs []*SongRelationshipsSong
	for _, hit := range r.Response.Hits {
		songs = append(songs, &hit.Result)
	}

	return songs, resp, nil
}
