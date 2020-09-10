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

// FlairChoice is a choice of flair when selecting one for yourself or for a post.
type FlairChoice struct {
	TemplateID string `json:"flair_template_id"`
	Text       string `json:"flair_text"`
	Editable   bool   `json:"flair_text_editable"`
	Position   string `json:"flair_position"`
	CSSClass   string `json:"flair_css_class"`
}

// ConfigureFlairRequest represents a request to configure a subreddit's flair settings.
// Not setting an attribute can have unexpected side effects, so assign every one just in case.
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

// FlairTemplateCreateOrUpdateRequest represents a request to create/update a flair template.
// Not setting an attribute can have unexpected side effects, so assign every one just in case.
type FlairTemplateCreateOrUpdateRequest struct {
	// The id of the template. Only provide this if it's an update request.
	// If provided and it's not a valid id, the template will be created.
	ID string `url:"flair_template_id,omitempty"`

	// One of: all, emoji, text.
	AllowableContent string `url:"allowable_content,omitempty"`
	// No longer than 64 characters.
	Text string `url:"text,omitempty"`
	// One of: light, dark.
	TextColor string `url:"text_color,omitempty"`
	// Allow user to edit the text of the flair.
	TextEditable *bool `url:"text_editable,omitempty"`

	ModOnly *bool `url:"mod_only,omitempty"`

	// Between 1 and 10 (inclusive). Default: 10.
	MaxEmojis *int `url:"max_emojis,omitempty"`

	// One of: none, transparent, 6-digit rgb hex color, e.g. #AABBCC.
	BackgroundColor string `url:"background_color,omitempty"`
	CSSClass        string `url:"css_class,omitempty"`
}

// FlairTemplate is a generic flair structure that can users can use next to their username
// or posts in a subreddit.
type FlairTemplate struct {
	ID string `json:"id"`
	// USER_FLAIR (for users) or LINK_FLAIR (for posts).
	Type    string `json:"flairType"`
	ModOnly bool   `json:"modOnly"`

	AllowableContent string              `json:"allowableContent"`
	Text             string              `json:"text"`
	TextType         string              `json:"type"`
	TextColor        string              `json:"textColor"`
	TextEditable     bool                `json:"textEditable"`
	RichText         []map[string]string `json:"richtext"`

	OverrideCSS     bool   `json:"overrideCss"`
	MaxEmojis       int    `json:"maxEmojis"`
	BackgroundColor string `json:"backgroundColor"`
	CSSClass        string `json:"cssClass"`
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

	req, err := s.client.NewRequest(http.MethodPost, path, form)
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

	req, err := s.client.NewRequest(http.MethodPost, path, form)
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

	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UpsertUserTemplate creates a user flair template, or updates it if the request.ID is valid.
// It returns the created/updated flair template.
func (s *FlairService) UpsertUserTemplate(ctx context.Context, subreddit string, request *FlairTemplateCreateOrUpdateRequest) (*FlairTemplate, *Response, error) {
	if request == nil {
		return nil, nil, errors.New("request: cannot be nil")
	}

	path := fmt.Sprintf("r/%s/api/flairtemplate_v2", subreddit)

	form, err := query.Values(request)
	if err != nil {
		return nil, nil, err
	}
	form.Set("api_type", "json")
	form.Set("flair_type", "USER_FLAIR")

	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(FlairTemplate)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// UpsertPostTemplate creates a post flair template, or updates it if the request.ID is valid.
// It returns the created/updated flair template.
func (s *FlairService) UpsertPostTemplate(ctx context.Context, subreddit string, request *FlairTemplateCreateOrUpdateRequest) (*FlairTemplate, *Response, error) {
	if request == nil {
		return nil, nil, errors.New("request: cannot be nil")
	}

	path := fmt.Sprintf("r/%s/api/flairtemplate_v2", subreddit)

	form, err := query.Values(request)
	if err != nil {
		return nil, nil, err
	}
	form.Set("api_type", "json")
	form.Set("flair_type", "LINK_FLAIR")

	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(FlairTemplate)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Delete the flair of the user.
func (s *FlairService) Delete(ctx context.Context, subreddit, username string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/deleteflair", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("name", username)

	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteTemplate deletes the flair template via its id.
func (s *FlairService) DeleteTemplate(ctx context.Context, subreddit, id string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/deleteflairtemplate", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("flair_template_id", id)

	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteAllUserTemplates deletes all user flair templates.
func (s *FlairService) DeleteAllUserTemplates(ctx context.Context, subreddit string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/clearflairtemplates", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("flair_type", "USER_FLAIR")

	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteAllPostTemplates deletes all post flair templates.
func (s *FlairService) DeleteAllPostTemplates(ctx context.Context, subreddit string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/clearflairtemplates", subreddit)

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("flair_type", "LINK_FLAIR")

	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ReorderUserTemplates reorders the user flair templates in the order provided in the slice.
// The order should contain every single flair id of this flair type; omitting any id will result in an error.
func (s *FlairService) ReorderUserTemplates(ctx context.Context, subreddit string, ids []string) (*Response, error) {
	path := fmt.Sprintf("api/v1/%s/flair_template_order/USER_FLAIR", subreddit)
	req, err := s.client.NewJSONRequest(http.MethodPatch, path, ids)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ReorderPostTemplates reorders the post flair templates in the order provided in the slice.
// The order should contain every single flair id of this flair type; omitting any id will result in an error.
func (s *FlairService) ReorderPostTemplates(ctx context.Context, subreddit string, ids []string) (*Response, error) {
	path := fmt.Sprintf("api/v1/%s/flair_template_order/LINK_FLAIR", subreddit)
	req, err := s.client.NewJSONRequest(http.MethodPatch, path, ids)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// Choices returns a list of flairs you can assign to yourself in the subreddit, and your current one.
func (s *FlairService) Choices(ctx context.Context, subreddit string) ([]*FlairChoice, *FlairChoice, *Response, error) {
	return s.ChoicesOf(ctx, subreddit, s.client.Username)
}

// ChoicesOf returns a list of flairs the user can assign to themself in the subreddit, and their current one.
// Unless the user is you, this only works if you're a moderator of the subreddit.
func (s *FlairService) ChoicesOf(ctx context.Context, subreddit, username string) ([]*FlairChoice, *FlairChoice, *Response, error) {
	path := fmt.Sprintf("r/%s/api/flairselector", subreddit)
	form := url.Values{}
	form.Set("name", username)
	return s.choices(ctx, path, form)
}

// ChoicesForPost returns a list of flairs you can assign to an existing post, and the current one assigned to it.
// If the post isn't yours, this only works if you're the moderator of the subreddit it's in.
func (s *FlairService) ChoicesForPost(ctx context.Context, postID string) ([]*FlairChoice, *FlairChoice, *Response, error) {
	path := "api/flairselector"
	form := url.Values{}
	form.Set("link", postID)
	return s.choices(ctx, path, form)
}

// ChoicesForNewPost returns a list of flairs you can assign to a new post in a subreddit.
func (s *FlairService) ChoicesForNewPost(ctx context.Context, subreddit string) ([]*FlairChoice, *Response, error) {
	path := fmt.Sprintf("r/%s/api/flairselector", subreddit)

	form := url.Values{}
	form.Set("is_newlink", "true")

	choices, _, resp, err := s.choices(ctx, path, form)
	return choices, resp, err
}

func (s *FlairService) choices(ctx context.Context, path string, form url.Values) ([]*FlairChoice, *FlairChoice, *Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, path, form)
	if err != nil {
		return nil, nil, nil, err
	}

	root := new(struct {
		Choices []*FlairChoice `json:"choices"`
		Current *FlairChoice   `json:"current"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, resp, err
	}

	return root.Choices, root.Current, resp, nil
}
