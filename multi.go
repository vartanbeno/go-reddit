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

// MultiService handles communication with the multireddit
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api#section_multis
type MultiService struct {
	client *Client
}

type multiRoot struct {
	Kind string `json:"kind,omitempty"`
	Data *Multi `json:"data,omitempty"`
}

// Multi is a multireddit, i.e. a customizable group of subreddits.
// Users can create multis for custom navigation, instead of browsing
// one subreddit or all subreddits at a time.
type Multi struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	// Format: user/{username}/m/{multiname}
	Path        string         `json:"path,omitempty"`
	Description string         `json:"description_md,omitempty"`
	Subreddits  SubredditNames `json:"subreddits"`
	CopiedFrom  *string        `json:"copied_from"`

	Owner   string     `json:"owner,omitempty"`
	OwnerID string     `json:"owner_id,omitempty"`
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
	type subreddit struct {
		Name string `json:"name"`
	}
	var subreddits []subreddit

	err := json.Unmarshal(data, &subreddits)
	if err != nil {
		return err
	}

	for _, sr := range subreddits {
		*n = append(*n, sr.Name)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n *SubredditNames) MarshalJSON() ([]byte, error) {
	type subreddit struct {
		Name string `json:"name"`
	}
	var subreddits []subreddit

	for _, name := range *n {
		subreddits = append(subreddits, subreddit{name})
	}

	return json.Marshal(subreddits)
}

// MultiCopyRequest represents a request to copy a multireddit.
type MultiCopyRequest struct {
	FromPath string `url:"from"`
	ToPath   string `url:"to"`
	// Raw markdown text.
	Description string `url:"description_md,omitempty"`
	// No longer than 50 characters.
	DisplayName string `url:"display_name,omitempty"`
}

// MultiCreateOrUpdateRequest represents a request to create/update a multireddit.
type MultiCreateOrUpdateRequest struct {
	// For updates, this is the display name, i.e. the header of the multi.
	// Not part of the path necessarily.
	Name        string         `json:"display_name,omitempty"`
	Description string         `json:"description_md,omitempty"`
	Subreddits  SubredditNames `json:"subreddits,omitempty"`
	// One of: private, public, hidden
	Visibility string `json:"visibility,omitempty"`
}

// Form parameterizes the fields and returns the form.
func (r *MultiCreateOrUpdateRequest) Form() url.Values {
	byteValue, _ := json.Marshal(r)
	form := url.Values{}
	form.Set("model", string(byteValue))
	return form
}

type rootMultiDescription struct {
	Data struct {
		Body string `json:"body_md"`
	} `json:"data"`
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
// Or, if the user is you, all of your multireddits.
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
		return nil, nil, errors.New("copyRequest: cannot be nil")
	}

	path := "api/multi/copy"
	form, err := query.Values(copyRequest)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
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
func (s *MultiService) Create(ctx context.Context, createRequest *MultiCreateOrUpdateRequest) (*Multi, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("createRequest: cannot be nil")
	}

	path := "api/multi"

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, createRequest.Form())
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

// Update updates a multireddit.
// If the multireddit does not exist, it will be created.
func (s *MultiService) Update(ctx context.Context, multiPath string, updateRequest *MultiCreateOrUpdateRequest) (*Multi, *Response, error) {
	if updateRequest == nil {
		return nil, nil, errors.New("updateRequest: cannot be nil")
	}

	path := fmt.Sprintf("api/multi/%s", multiPath)

	req, err := s.client.NewRequestWithForm(http.MethodPut, path, updateRequest.Form())
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

// GetDescription gets a multireddit's description.
func (s *MultiService) GetDescription(ctx context.Context, multiPath string) (string, *Response, error) {
	path := fmt.Sprintf("api/multi/%s/description", multiPath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", nil, err
	}

	root := new(rootMultiDescription)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}

	return root.Data.Body, resp, nil
}

// UpdateDescription updates a multireddit's description.
func (s *MultiService) UpdateDescription(ctx context.Context, multiPath string, description string) (string, *Response, error) {
	path := fmt.Sprintf("api/multi/%s/description", multiPath)

	form := url.Values{}
	form.Set("model", fmt.Sprintf(`{"body_md":"%s"}`, description))

	req, err := s.client.NewRequestWithForm(http.MethodPut, path, form)
	if err != nil {
		return "", nil, err
	}

	root := new(rootMultiDescription)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}

	return root.Data.Body, resp, nil
}

// AddSubreddit adds a subreddit to a multireddit.
func (s *MultiService) AddSubreddit(ctx context.Context, multiPath string, subreddit string) (*Response, error) {
	path := fmt.Sprintf("api/multi/%s/r/%s", multiPath, subreddit)

	form := url.Values{}
	form.Set("model", fmt.Sprintf(`{"name":"%s"}`, subreddit))

	req, err := s.client.NewRequestWithForm(http.MethodPut, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteSubreddit removes a subreddit from a multireddit.
func (s *MultiService) DeleteSubreddit(ctx context.Context, multiPath string, subreddit string) (*Response, error) {
	path := fmt.Sprintf("api/multi/%s/r/%s", multiPath, subreddit)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
