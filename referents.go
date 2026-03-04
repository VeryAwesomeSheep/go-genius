package genius

import (
	"context"
	"net/http"
)

type ReferentsService Service

// Referent represents data about the referent.
type Referent struct {
	AnnotatorID    int           `json:"annotator_id"`
	AnnotatorLogin string        `json:"annotator_login"`
	APIPath        string        `json:"api_path"`
	Classification string        `json:"classification"`
	Featured       *bool         `json:"featured"`
	Fragment       string        `json:"fragment"`
	ID             int           `json:"id"`
	IsDescription  bool          `json:"is_description"`
	Path           string        `json:"path"`
	Range          Range         `json:"range"`
	SongID         *int          `json:"song_id"`
	URL            string        `json:"url"`
	Annotatable    *Annotatable  `json:"annotatable"`
	Annotations    []*Annotation `json:"annotations"` // Available only via ReferentsService
}

type Range struct {
	Start       *string `json:"start"`
	StartOffset *string `json:"startOffset"`
	End         *string `json:"end"`
	EndOffset   *string `json:"endOffset"`
	Before      *string `json:"before"`
	After       *string `json:"after"`
	Content     string  `json:"content"`
}

type Annotatable struct {
	APIPath  string  `json:"api_path"`
	Context  *string `json:"context"`
	ID       int     `json:"id"`
	ImageURL string  `json:"image_url"`
	Title    string  `json:"title"`
	Type     string  `json:"type"`
	URL      string  `json:"url"`
}

type ReferentsOptions struct {
	CreatedByID int `url:"created_by_id,omitempty"`
	SongID      int `url:"song_id,omitempty"`
	WebPageID   int `url:"web_page_id,omitempty"`
	TextFormatOptions
	PagingOptions
}

// Get returns list of referents by song_id or web_page_id.
func (s *ReferentsService) Get(ctx context.Context, opts *ReferentsOptions) ([]*Referent, *http.Response, error) {
	u := "referents"

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(u)
	if err != nil {
		return nil, nil, err
	}

	var r Response[struct {
		Referents []*Referent `json:"referents"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r.Response.Referents, resp, nil
}
