package genius

import (
	"context"
	"net/http"
)

type WebPagesService Service

type WebPage struct {
	APIPath         *string `json:"api_path"`
	Domain          string  `json:"domain"`
	ID              *int    `json:"id"`
	NormalizedURL   string  `json:"normalized_url"`
	ShareURL        string  `json:"share_url"`
	Title           string  `json:"title"`
	URL             string  `json:"url"`
	AnnotationCount int     `json:"annotation_count"`
}

type WebPagesOptions struct {
	RawAnnotatableURL string `url:"raw_annotatable_url,omitempty"`
	CanonicalURL      string `url:"canonical_url,omitempty"`
	OgURL             string `url:"og_url,omitempty"`
}

func (s *WebPagesService) Get(ctx context.Context, opts *WebPagesOptions) (*WebPage, *http.Response, error) {
	u := "web_pages/lookup"

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r Response[struct {
		WebPage *WebPage `json:"web_page"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r.Response.WebPage, resp, nil
}
