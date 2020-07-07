package geddit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// MultiService handles communication with the multireddit
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api#section_multis
type MultiService service

type multiRoot struct {
	Kind string `json:"kind,omitempty"`
	Data *Multi `json:"data,omitempty"`
}

// Multi is a multireddit, i.e. a customizable group of subreddits.
// Users can create multis for custom navigation, instead of browsing
// one subreddit or all subreddits at a time.
type Multi struct {
	Name        string         `json:"name,omitempty"`
	DisplayName string         `json:"display_name,omitempty"`
	Path        string         `json:"path,omitempty"`
	Description string         `json:"description_md,omitempty"`
	Subreddits  SubredditNames `json:"subreddits"`
	CopedFrom   *string        `json:"copied_from"`

	Owner   string     `json:"owner,omitempty"`
	OwnerID string     `json:"ownerID,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	NumberOfSubscribers int    `json:"num_subscribers"`
	Visibility          string `json:"visibility,omitempty"`
	Subscribed          bool   `json:"is_subscriber"`
	Favorite            bool   `json:"is_favorited"`
	CanEdit             bool   `json:"can_edit"`
	NSFW                bool   `json:"over_18"`
}

// SubredditNames is a list of subreddit names.
type SubredditNames []string

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *SubredditNames) UnmarshalJSON(data []byte) error {
	var subreddits []map[string]string

	err := json.Unmarshal(data, &subreddits)
	if err != nil {
		return err
	}

	for _, subreddit := range subreddits {
		name, ok := subreddit["name"]
		if !ok {
			continue
		}
		*n = append(*n, name)
	}

	return nil
}

// MultiCopyRequest represents a request to copy a multireddit.
type MultiCopyRequest struct {
	From string
	To   string
	// Raw markdown text.
	Description string
	// No longer than 50 characters.
	DisplayName string
}

// Form parameterizes the fields and returns the form.
func (r *MultiCopyRequest) Form() url.Values {
	form := url.Values{}
	form.Set("from", r.From)
	form.Set("to", r.To)
	form.Set("description_md", r.Description)
	form.Set("display_name", r.DisplayName)
	return form
}

// MultiCreateRequest represents a request to create a multireddit.
type MultiCreateRequest struct {
	Description string   `json:"description_md,omitempty"`
	DisplayName string   `json:"display_name,omitempty"`
	Subreddits  []string `json:"subreddits"`
	Visibility  string   `json:"visibility,omitempty"`
}

// Form parameterizes the fields and returns the form.
func (r *MultiCreateRequest) Form() url.Values {
	byteValue, _ := json.Marshal(r)
	form := url.Values{}
	form.Set("model", string(byteValue))
	return form
}

// Get gets information about the multireddit from its url path.
func (s *MultiService) Get(ctx context.Context, multiPath string) (*Multi, *Response, error) {
	path := fmt.Sprintf("api/multi/%s", multiPath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(multiRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, nil
}

// Mine returns your multireddits.
func (s *MultiService) Mine(ctx context.Context) ([]Multi, *Response, error) {
	path := "api/multi/mine"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []multiRoot
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	multis := make([]Multi, 0)
	for _, multi := range root {
		multis = append(multis, *multi.Data)
	}

	return multis, resp, nil
}

// Of returns the user's public multireddits.
func (s *MultiService) Of(ctx context.Context, username string) ([]Multi, *Response, error) {
	path := fmt.Sprintf("api/multi/user/%s", username)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []multiRoot
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	multis := make([]Multi, 0)
	for _, multi := range root {
		multis = append(multis, *multi.Data)
	}

	return multis, resp, nil
}

// Copy copies a multireddit.
func (s *MultiService) Copy(ctx context.Context, copyRequest *MultiCopyRequest) (*Multi, *Response, error) {
	if copyRequest == nil {
		return nil, nil, errors.New("copyRequest cannot be nil")
	}

	path := fmt.Sprintf("api/multi/copy")

	req, err := s.client.NewPostForm(path, copyRequest.Form())
	if err != nil {
		return nil, nil, err
	}

	root := new(multiRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, nil
}

// Create creates a multireddit.
func (s *MultiService) Create(ctx context.Context, createRequest *MultiCreateRequest) (*Multi, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("createRequest cannot be nil")
	}

	path := fmt.Sprintf("api/multi/copy")

	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(multiRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, nil
}

// Delete deletes a multireddit.
func (s *MultiService) Delete(ctx context.Context, multiPath string) (*Response, error) {
	path := fmt.Sprintf("api/multi/%s", multiPath)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
