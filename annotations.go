package genius

import (
	"context"
	"fmt"
	"net/http"
)

type AnnotationsService Service

type Annotation struct {
	APIPath      string   `json:"api_path"`
	Body         TextBody `json:"body"`
	CommentCount int      `json:"comment_count"`
	ID           int      `json:"id"`
	State        string   `json:"state"`
	URL          string   `json:"url"`
	Verified     bool     `json:"verified"`
	VotesTotal   int      `json:"votes_total"`
}

type AnnotationsOptions struct {
	TextFormatOptions
}

func (s *AnnotationsService) Get(ctx context.Context, id int, opts *AnnotationsOptions) (*Annotation, *Referent, *http.Response, error) {
	u := fmt.Sprintf("annotations/%d", id)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var r Response[struct {
		Annotation *Annotation `json:"annotation"`
		Referent   *Referent   `json:"referent"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, nil, resp, err
	}

	return r.Response.Annotation, r.Response.Referent, resp, nil
}
