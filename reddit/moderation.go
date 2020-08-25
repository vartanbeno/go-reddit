package reddit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
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

// ModPermissions are the different permissions moderators have or don't have on a subreddit.
// Read about them here: https://mods.reddithelp.com/hc/en-us/articles/360009381491-User-Management-moderators-and-permissions
type ModPermissions struct {
	All          bool
	Access       bool
	ChatConfig   bool
	ChatOperator bool
	Config       bool
	Flair        bool
	Mail         bool
	Posts        bool
	Wiki         bool
}

func (p *ModPermissions) String() (s string) {
	if p == nil {
		return "+all"
	}

	if p.All {
		s += "+"
	} else {
		s += "-"
	}
	s += "all,"

	if p.Access {
		s += "+"
	} else {
		s += "-"
	}
	s += "access,"

	if p.ChatConfig {
		s += "+"
	} else {
		s += "-"
	}
	s += "chat_config,"

	if p.ChatOperator {
		s += "+"
	} else {
		s += "-"
	}
	s += "chat_operator,"

	if p.Config {
		s += "+"
	} else {
		s += "-"
	}
	s += "config,"

	if p.Flair {
		s += "+"
	} else {
		s += "-"
	}
	s += "flair,"

	if p.Mail {
		s += "+"
	} else {
		s += "-"
	}
	s += "mail,"

	if p.Posts {
		s += "+"
	} else {
		s += "-"
	}
	s += "posts,"

	if p.Wiki {
		s += "+"
	} else {
		s += "-"
	}
	s += "wiki"

	return
}

// Invite a user to become a moderator of the subreddit.
// If permissions is nil, all permissions will be granted.
func (s *ModerationService) Invite(ctx context.Context, subreddit string, username string, permissions *ModPermissions) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/friend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "moderator_invite")
	form.Set("permissions", permissions.String())

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Uninvite a user from becoming a moderator of the subreddit.
func (s *ModerationService) Uninvite(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/unfriend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "moderator_invite")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// BanConfig configures the ban of the user being banned.
type BanConfig struct {
	Reason string `url:"reason,omitempty"`
	// Not visible to the user being banned.
	ModNote string `url:"note,omitempty"`
	// How long the ban will last. 0-999. Leave nil for permanent.
	Days *int `url:"duration,omitempty"`
	// Note to include in the ban message to the user.
	Message string `url:"ban_message,omitempty"`
}

// Ban a user from the subreddit.
func (s *ModerationService) Ban(ctx context.Context, subreddit string, username string, config *BanConfig) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/friend", subreddit)

	form, err := query.Values(config)
	if err != nil {
		return nil, err
	}

	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "banned")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unban a user from the subreddit.
func (s *ModerationService) Unban(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/unfriend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "banned")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// BanWiki a user from contributing to the subreddit wiki.
func (s *ModerationService) BanWiki(ctx context.Context, subreddit string, username string, config *BanConfig) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/friend", subreddit)

	form, err := query.Values(config)
	if err != nil {
		return nil, err
	}

	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "wikibanned")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UnbanWiki a user from contributing to the subreddit wiki.
func (s *ModerationService) UnbanWiki(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/unfriend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "wikibanned")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Mute a user in the subreddit.
func (s *ModerationService) Mute(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/friend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "muted")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unmute a user in the subreddit.
func (s *ModerationService) Unmute(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/unfriend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "muted")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ApproveUser adds a user as an approved user to the subreddit.
func (s *ModerationService) ApproveUser(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/friend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "contributor")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UnapproveUser removes a user as an approved user to the subreddit.
func (s *ModerationService) UnapproveUser(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/unfriend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "contributor")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ApproveUserWiki adds a user as an approved wiki contributor in the subreddit.
func (s *ModerationService) ApproveUserWiki(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/friend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "wikicontributor")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UnapproveUserWiki removes a user as an approved wiki contributor in the subreddit.
func (s *ModerationService) UnapproveUserWiki(ctx context.Context, subreddit string, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/unfriend", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)
	form.Set("type", "wikicontributor")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
