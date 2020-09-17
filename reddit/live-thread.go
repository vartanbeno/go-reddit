package reddit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// LiveThreadService handles communication with the live thread
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_live
type LiveThreadService struct {
	client *Client
}

// LiveThread is a thread on Reddit that provides real-time updates.
type LiveThread struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Resources   string `json:"resources,omitempty"`

	State             string `json:"state,omitempty"`
	ViewerCount       int    `json:"viewer_count"`
	ViewerCountFuzzed bool   `json:"viewer_count_fuzzed"`

	// Empty when a list thread has ended.
	WebSocketURL string `json:"websocket_url,omitempty"`

	Announcement bool `json:"is_announcement"`
	NSFW         bool `json:"nsfw"`
}

// LiveThreadCreateRequest represents a request to create a live thread.
type LiveThreadCreateRequest struct {
	// No longer than 120 characters.
	Title       string `url:"title"`
	Description string `url:"description,omitempty"`
	Resources   string `url:"resources,omitempty"`
	NSFW        bool   `url:"nsfw,omitempty"`
}

// LiveThreadContributor is a user that can contribute to a live thread.
type LiveThreadContributor struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// LiveThreadContributors is a list of users that can contribute to a live thread.
type LiveThreadContributors struct {
	Current []*LiveThreadContributor `json:"current_contributors"`
	// This is only filled if you are a contributor in the live thread with the "manage" permission.
	Invited []*LiveThreadContributor `json:"invited_contributors,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *LiveThreadContributors) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.New("no bytes to unmarshal")
	}

	// neat trick taken from:
	// https://www.calhoun.io/how-to-parse-json-that-varies-between-an-array-or-a-single-item-with-go
	switch b[0] {
	case '{':
		return c.unmarshalSingle(b)
	case '[':
		return c.unmarshalMany(b)
	}

	// This shouldn't really happen as the standard library seems to strip
	// whitespace from the bytes being passed in, but just in case let's guess at
	// multiple tags and fall back to a single one if that doesn't work.
	err := c.unmarshalSingle(b)
	if err != nil {
		return c.unmarshalMany(b)
	}

	return nil
}

func (c *LiveThreadContributors) unmarshalSingle(b []byte) error {
	root := new(struct {
		Data struct {
			Children []*LiveThreadContributor `json:"children"`
		} `json:"data"`
	})

	err := json.Unmarshal(b, &root)
	if err != nil {
		return err
	}

	c.Current = root.Data.Children
	return nil
}

func (c *LiveThreadContributors) unmarshalMany(b []byte) error {
	var root [2]struct {
		Data struct {
			Children []*LiveThreadContributor `json:"children"`
		} `json:"data"`
	}

	err := json.Unmarshal(b, &root)
	if err != nil {
		return err
	}

	c.Current = root[0].Data.Children
	c.Invited = root[1].Data.Children
	return nil
}

// Get information about a live thread.
func (s *LiveThreadService) Get(ctx context.Context, id string) (*LiveThread, *Response, error) {
	path := fmt.Sprintf("live/%s/about", id)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	t, _ := root.LiveThread()
	return t, resp, nil
}

// Create a live thread and get its id.
func (s *LiveThreadService) Create(ctx context.Context, request *LiveThreadCreateRequest) (string, *Response, error) {
	if request == nil {
		return "", nil, errors.New("*LiveThreadCreateRequest: cannot be nil")
	}

	form, err := query.Values(request)
	if err != nil {
		return "", nil, err
	}
	form.Set("api_type", "json")

	path := "api/live/create"
	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return "", nil, err
	}

	root := new(struct {
		JSON struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		} `json:"json"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}

	return root.JSON.Data.ID, resp, nil
}

// Contributors gets a list of users that are contributors to the live thread.
// If you are a contributor and you have the "manage" permission (to manage contributors), you
// also get a list of invited contributors that haven't yet accepted/refused their invitation.
func (s *LiveThreadService) Contributors(ctx context.Context, id string) (*LiveThreadContributors, *Response, error) {
	path := fmt.Sprintf("live/%s/contributors", id)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(LiveThreadContributors)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Accept a pending invite to contribute to the live thread.
func (s *LiveThreadService) Accept(ctx context.Context, id string) (*Response, error) {
	form := url.Values{}
	form.Set("api_type", "json")

	path := fmt.Sprintf("api/live/%s/accept_contributor_invite", id)
	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Leave the live thread by abdicating your status as contributor.
// todo: test as the author who leaves the thread.
func (s *LiveThreadService) Leave(ctx context.Context, id string) (*Response, error) {
	form := url.Values{}
	form.Set("api_type", "json")

	path := fmt.Sprintf("api/live/%s/leave_contributor", id)
	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
