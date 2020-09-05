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

// WikiService handles communication with the wiki
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_wiki
type WikiService struct {
	client *Client
}

// WikiPage is a wiki page in a subreddit.
type WikiPage struct {
	Content   string `json:"content_md,omitempty"`
	Reason    string `json:"reason,omitempty"`
	MayRevise bool   `json:"may_revise"`

	RevisionID   string     `json:"revision_id,omitempty"`
	RevisionDate *Timestamp `json:"revision_date,omitempty"`
	RevisionBy   *User      `json:"revision_by,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *WikiPage) UnmarshalJSON(b []byte) error {
	root := new(struct {
		Content   string `json:"content_md,omitempty"`
		Reason    string `json:"reason,omitempty"`
		MayRevise bool   `json:"may_revise"`

		RevisionID   string     `json:"revision_id,omitempty"`
		RevisionDate *Timestamp `json:"revision_date,omitempty"`
		RevisionBy   thing      `json:"revision_by,omitempty"`
	})

	err := json.Unmarshal(b, root)
	if err != nil {
		return err
	}

	p.Content = root.Content
	p.Reason = root.Reason
	p.MayRevise = root.MayRevise

	p.RevisionID = root.RevisionID
	p.RevisionDate = root.RevisionDate

	if user, ok := root.RevisionBy.User(); ok {
		p.RevisionBy = user
	}

	return nil
}

// WikiPagePermissionLevel defines who can edit a specific wiki page in a subreddit.
type WikiPagePermissionLevel int

const (
	// PermissionSubredditWikiPermissions uses subreddit wiki permissions.
	PermissionSubredditWikiPermissions WikiPagePermissionLevel = iota
	// PermissionApprovedContributorsOnly is only for approved wiki contributors.
	PermissionApprovedContributorsOnly
	// PermissionModeratorsOnly is only for moderators.
	PermissionModeratorsOnly
)

// WikiPageSettings holds the settings for a specific wiki page.
type WikiPageSettings struct {
	PermissionLevel WikiPagePermissionLevel `json:"permlevel"`
	Listed          bool                    `json:"listed"`
	Editors         []*User                 `json:"editors"`
}

// WikiPageSettingsUpdateRequest represents a request to update the visibility and
// permissions of a wiki page.
type WikiPageSettingsUpdateRequest struct {
	// This HAS to be provided no matter what, or else we get a 500 response.
	PermissionLevel WikiPagePermissionLevel `url:"permlevel"`
	Listed          *bool                   `url:"listed,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *WikiPageSettings) UnmarshalJSON(b []byte) error {
	root := new(struct {
		PermissionLevel WikiPagePermissionLevel `json:"permlevel"`
		Listed          bool                    `json:"listed"`
		Things          []thing                 `json:"editors"`
	})

	err := json.Unmarshal(b, root)
	if err != nil {
		return err
	}

	s.PermissionLevel = root.PermissionLevel
	s.Listed = root.Listed

	for _, thing := range root.Things {
		if user, ok := thing.User(); ok {
			s.Editors = append(s.Editors, user)
		}
	}

	return nil
}

// Page gets a wiki page.
func (s *WikiService) Page(ctx context.Context, subreddit, page string) (*WikiPage, *Response, error) {
	path := fmt.Sprintf("r/%s/wiki/%s", subreddit, page)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	wikiPage, _ := root.WikiPage()
	return wikiPage, resp, nil
}

// Pages retrieves a list of wiki pages in the subreddit.
// Returns 403 Forbidden if the wiki is disabled.
func (s *WikiService) Pages(ctx context.Context, subreddit string) ([]string, *Response, error) {
	path := fmt.Sprintf("r/%s/wiki/pages", subreddit)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	wikiPages, _ := root.WikiPages()
	return wikiPages, resp, nil
}

// Settings gets the subreddit's wiki page's settings.
func (s *WikiService) Settings(ctx context.Context, subreddit, page string) (*WikiPageSettings, *Response, error) {
	path := fmt.Sprintf("r/%s/wiki/settings/%s", subreddit, page)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	settings, _ := root.WikiPageSettings()
	return settings, resp, nil
}

// UpdateSettings updates the subreddit's wiki page's settings.
func (s *WikiService) UpdateSettings(ctx context.Context, subreddit, page string, updateRequest *WikiPageSettingsUpdateRequest) (*WikiPageSettings, *Response, error) {
	if updateRequest == nil {
		return nil, nil, errors.New("updateRequest: cannot be nil")
	}

	form, err := query.Values(updateRequest)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("r/%s/wiki/settings/%s", subreddit, page)
	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	settings, _ := root.WikiPageSettings()
	return settings, resp, nil
}

// Allow the user to edit the specified wiki page in the subreddit.
func (s *WikiService) Allow(ctx context.Context, subreddit, page, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/wiki/alloweditor/add", subreddit)

	form := url.Values{}
	form.Set("page", page)
	form.Set("username", username)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Deny the user the ability to edit the specified wiki page in the subreddit.
func (s *WikiService) Deny(ctx context.Context, subreddit, page, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/wiki/alloweditor/del", subreddit)

	form := url.Values{}
	form.Set("page", page)
	form.Set("username", username)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Discussions gets a list of discussions (posts) about the wiki page.
func (s *WikiService) Discussions(ctx context.Context, subreddit, page string, opts *ListOptions) ([]*Post, *Response, error) {
	path := fmt.Sprintf("r/%s/wiki/discussions/%s", subreddit, page)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	listing, _ := root.Listing()
	return listing.Posts(), resp, nil
}
