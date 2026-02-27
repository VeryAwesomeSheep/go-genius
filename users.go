package genius

import (
	"context"
	"fmt"
	"net/http"
)

type UsersService Service

// User represents data about a user.
type User struct {
	AboutMe                     TextBody   `json:"about_me"`
	APIPath                     string     `json:"api_path"`
	Avatar                      AvatarSet  `json:"avatar"`
	CustomHeaderImageURL        *string    `json:"custom_header_image_url"`
	FollowedUsersCount          int        `json:"followed_users_count"`
	FollowersCount              int        `json:"followers_count"`
	HeaderImageURL              string     `json:"header_image_url"`
	HumanReadableRoleForDisplay string     `json:"human_readable_role_for_display"`
	ID                          int        `json:"id"`
	IQ                          int        `json:"iq"`
	IQForDisplay                string     `json:"iq_for_display"`
	Login                       string     `json:"login"`
	Name                        string     `json:"name"`
	PhotoURL                    string     `json:"photo_url"`
	RoleForDisplay              string     `json:"role_for_display"`
	RolesForDisplay             []string   `json:"roles_for_display"`
	URL                         string     `json:"url"`
	Artist                      *Artist    `json:"artist"`
	Stats                       *UserStats `json:"stats"`
}

type AvatarSet struct {
	Tiny   Avatar `json:"tiny"`
	Thumb  Avatar `json:"thumb"`
	Small  Avatar `json:"small"`
	Medium Avatar `json:"medium"`
}

type Avatar struct {
	URL         string      `json:"url"`
	BoundingBox BoundingBox `json:"bounding_box"`
}

type BoundingBox struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type UserStats struct {
	AnnotationsCount    int `json:"annotations_count"`
	AnswersCount        int `json:"answers_count"`
	CommentsCount       int `json:"comments_count"`
	ForumPostsCount     int `json:"forum_posts_count"`
	PyongsCount         int `json:"pyongs_count"`
	QuestionsCount      int `json:"questions_count"`
	TranscriptionsCount int `json:"transcriptions_count"`
}

type UserOptions struct {
	TextFormatOptions
}

// Get returns data for a user by its ID.
func (s *UsersService) Get(ctx context.Context, id int, opts *UserOptions) (*User, *http.Response, error) {
	u := fmt.Sprintf("users/%d", id)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(u)
	if err != nil {
		return nil, nil, err
	}

	var r Response[struct {
		User *User `json:"user"`
	}]

	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r.Response.User, resp, nil
}
