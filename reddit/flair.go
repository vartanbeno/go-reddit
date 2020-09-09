package reddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
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

// ConfigureFlairRequest represents a request to configure a subreddit's flair settings.
// Not setting an attribute can have unexpected side effects.
type ConfigureFlairRequest struct {
	// Enable user flair in the subreddit.
	UserFlairEnabled *bool `url:"flair_enabled,omitempty"`
	// One of: left, right.
	UserFlairPosition string `url:"flair_position,omitempty"`
	// Allow users to assign their own flair.
	UserFlairSelfAssignEnabled *bool `url:"flair_self_assign_enabled,omitempty"`
	// One of: none, left, right.
	PostFlairPosition string `url:"link_flair_position,omitempty"`
	// Allow submitters to assign their own post flair.
	PostFlairSelfAssignEnabled *bool `url:"link_flair_self_assign_enabled,omitempty"`
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

	root := new(struct {
		UserFlairs []*FlairSummary `json:"users"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UserFlairs, resp, nil
}

// Configure the subreddit's flair settings.
func (s *FlairService) Configure(ctx context.Context, subreddit string, request *ConfigureFlairRequest) (*Response, error) {
	if request == nil {
		return nil, errors.New("request: cannot be nil")
	}

	path := fmt.Sprintf("r/%s/api/flairconfig", subreddit)

	form, err := query.Values(request)
	if err != nil {
		return nil, err
	}
	form.Set("api_type", "json")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Enable your flair in the subreddit.
func (s *FlairService) Enable(ctx context.Context, subreddit string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/setflairenabled", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("flair_enabled", "true")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Disable your flair in the subreddit.
func (s *FlairService) Disable(ctx context.Context, subreddit string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/setflairenabled", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("flair_enabled", "false")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
