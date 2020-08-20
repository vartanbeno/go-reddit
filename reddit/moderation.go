package reddit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ModerationService handles communication with the moderation
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_moderation
type ModerationService struct {
	client *Client
}

// ModAction is an action executed by a moderator of a subreddit, such
// as inviting another user to be a mod, or setting permissions.
type ModAction struct {
	ID      string     `json:"id,omitempty"`
	Action  string     `json:"action,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	Moderator string `json:"mod,omitempty"`
	// Not the full ID, just the ID36.
	ModeratorID string `json:"mod_id36,omitempty"`

	// The author of whatever the action was produced on, e.g. a user, post, comment, etc.
	TargetAuthor string `json:"target_author,omitempty"`
	// This is the full ID of whatever the target was.
	TargetID        string `json:"target_fullname,omitempty"`
	TargetTitle     string `json:"target_title,omitempty"`
	TargetPermalink string `json:"target_permalink,omitempty"`
	TargetBody      string `json:"target_body,omitempty"`

	Subreddit string `json:"subreddit,omitempty"`
	// Not the full ID, just the ID36.
	SubredditID string `json:"sr_id36,omitempty"`
}

// GetActions gets a list of moderator actions on a subreddit.
func (s *ModerationService) GetActions(ctx context.Context, subreddit string, opts *ListModActionOptions) (*ModActions, *Response, error) {
	path := fmt.Sprintf("r/%s/about/log", subreddit)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	path, err = addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getModActions(), resp, nil
}

// AcceptInvite accepts a pending invite to moderate the specified subreddit.
func (s *ModerationService) AcceptInvite(ctx context.Context, subreddit string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/accept_moderator_invite", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Approve approves a post or comment via its full ID.
func (s *ModerationService) Approve(ctx context.Context, id string) (*Response, error) {
	path := "api/approve"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Remove removes a post, comment or modmail message via its full ID.
func (s *ModerationService) Remove(ctx context.Context, id string) (*Response, error) {
	path := "api/remove"

	form := url.Values{}
	form.Set("id", id)
	form.Set("spam", "false")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RemoveSpam removes a post, comment or modmail message via its full ID and marks it as spam.
func (s *ModerationService) RemoveSpam(ctx context.Context, id string) (*Response, error) {
	path := "api/remove"

	form := url.Values{}
	form.Set("id", id)
	form.Set("spam", "true")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Leave abdicates your moderator status in a subreddit via its full ID.
func (s *ModerationService) Leave(ctx context.Context, subredditID string) (*Response, error) {
	path := "api/leavemoderator"

	form := url.Values{}
	form.Set("id", subredditID)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// LeaveContributor abdicates your approved user status in a subreddit via its full ID.
func (s *ModerationService) LeaveContributor(ctx context.Context, subredditID string) (*Response, error) {
	path := "api/leavecontributor"

	form := url.Values{}
	form.Set("id", subredditID)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Edited gets posts and comments that have been edited recently.
func (s *ModerationService) Edited(ctx context.Context, subreddit string, opts *ListOptions) (*Posts, *Comments, *Response, error) {
	path := fmt.Sprintf("r/%s/about/edited", subreddit)

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, nil, err
	}

	return root.getPosts(), root.getComments(), resp, nil
}

// IgnoreReports prevents reports on a post or comment from causing notifications.
func (s *ModerationService) IgnoreReports(ctx context.Context, id string) (*Response, error) {
	path := "api/ignore_reports"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UnignoreReports allows reports on a post or comment to cause notifications.
func (s *ModerationService) UnignoreReports(ctx context.Context, id string) (*Response, error) {
	path := "api/unignore_reports"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
