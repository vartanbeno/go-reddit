package reddit

import (
	"context"
	"fmt"
	"net/http"
)

// FlairService handles communication with the flair
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_flair
type FlairService struct {
	client *Client
}

// Flair is a tag that can be attached to a user or a post.
type Flair struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`

	Color           string `json:"text_color,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	CSSClass        string `json:"css_class,omitempty"`

	Editable bool `json:"text_editable"`
	ModOnly  bool `json:"mod_only"`
}

// FlairSummary is a condensed version of Flair.
type FlairSummary struct {
	User     string `json:"user,omitempty"`
	Text     string `json:"flair_text,omitempty"`
	CSSClass string `json:"flair_css_class,omitempty"`
}

// GetUserFlairs returns the user flairs from the subreddit.
func (s *FlairService) GetUserFlairs(ctx context.Context, subreddit string) ([]*Flair, *Response, error) {
	path := fmt.Sprintf("r/%s/api/user_flair_v2", subreddit)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var flairs []*Flair
	resp, err := s.client.Do(ctx, req, &flairs)
	if err != nil {
		return nil, resp, err
	}

	return flairs, resp, nil
}

// GetPostFlairs returns the post flairs from the subreddit.
func (s *FlairService) GetPostFlairs(ctx context.Context, subreddit string) ([]*Flair, *Response, error) {
	path := fmt.Sprintf("r/%s/api/link_flair_v2", subreddit)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var flairs []*Flair
	resp, err := s.client.Do(ctx, req, &flairs)
	if err != nil {
		return nil, resp, err
	}

	return flairs, resp, nil
}

// ListUserFlairs returns all flairs of individual users in the subreddit.
func (s *FlairService) ListUserFlairs(ctx context.Context, subreddit string) ([]*FlairSummary, *Response, error) {
	path := fmt.Sprintf("r/%s/api/flairlist", subreddit)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root struct {
		UserFlairs []*FlairSummary `json:"users"`
	}
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.UserFlairs, resp, nil
}
