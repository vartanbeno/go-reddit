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
// By default, the limit parameter is 25; max is 500.
func (s *ModerationService) GetActions(ctx context.Context, subreddit string, opts *ListOptions) (*ModActions, *Response, error) {
	return s.GetActionsByType(ctx, subreddit, "", opts)
}

// GetActionsByType gets a list of moderator actions of the specified type on a subreddit.
// By default, the limit parameter is 25; max is 500.
// The type should be one of: banuser, unbanuser, spamlink, removelink, approvelink, spamcomment,
// removecomment, approvecomment, addmoderator, showcomment, invitemoderator, uninvitemoderator,
// acceptmoderatorinvite, removemoderator, addcontributor, removecontributor, editsettings, editflair,
// distinguish, marknsfw, wikibanned, wikicontributor, wikiunbanned, wikipagelisted, removewikicontributor,
// wikirevise, wikipermlevel, ignorereports, unignorereports, setpermissions, setsuggestedsort, sticky,
// unsticky, setcontestmode, unsetcontestmode, lock, unlock, muteuser, unmuteuser, createrule, editrule,
// reorderrules, deleterule, spoiler, unspoiler, modmail_enrollment, community_styling, community_widgets,
// markoriginalcontent, collections, events, hidden_award, add_community_topics, remove_community_topics,
// create_scheduled_post, edit_scheduled_post, delete_scheduled_post, submit_scheduled_post, edit_post_requirements,
// invitesubscriber, submit_content_rating_survey.
func (s *ModerationService) GetActionsByType(ctx context.Context, subreddit string, actionType string, opts *ListOptions) (*ModActions, *Response, error) {
	path := fmt.Sprintf("r/%s/about/log", subreddit)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	type params struct {
		Type string `url:"type,omitempty"`
	}
	path, err = addOptions(path, params{actionType})
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
