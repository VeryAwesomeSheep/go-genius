package genius

import (
	"context"
	"fmt"
	"net/http"
)

type AnnotationsService Service

// Annotation represents data about an annotation.
type Annotation struct {
	APIPath          string             `json:"api_path"`
	Body             TextBody           `json:"body"`
	CommentCount     int                `json:"comment_count"`
	HasVoters        bool               `json:"has_voters"`
	ID               int                `json:"id"`
	Pinned           bool               `json:"pinned"`
	ShareURL         string             `json:"share_url"`
	State            string             `json:"state"`
	URL              string             `json:"url"`
	Verified         bool               `json:"verified"`
	VotesTotal       int                `json:"votes_total"`
	Authors          []AnnotationAuthor `json:"authors"`
	CosignedBy       []*User            `json:"cosigned_by"`
	RejectionComment *string            `json:"rejection_comment"`
	VerifiedBy       *User              `json:"verified_by"`
}

type AnnotationAuthor struct {
	Attribution float32 `json:"attribution"`
	User        User    `json:"user"`
}

type AnnotationsOptions struct {
	TextFormatOptions
}

// Get returns data for an annotation by its ID.
func (s *AnnotationsService) Get(ctx context.Context, id int, opts *AnnotationsOptions) (*Annotation, *Referent, *http.Response, error) {
	u := fmt.Sprintf("annotations/%d", id)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, nil, err
	}

	req, err := s.client.NewRequest(u)
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
