package geddit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// UserService handles communication with the user
// related methods of the Reddit API
type UserService interface {
	Get(ctx context.Context, username string) (*User, *Response, error)
	GetMultipleByID(ctx context.Context, ids ...string) (map[string]*UserShort, *Response, error)
	UsernameAvailable(ctx context.Context, username string) (bool, *Response, error)

	Overview(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Comments, *Response, error)
	OverviewOf(ctx context.Context, username string, opts ...SearchOptionSetter) (*Links, *Comments, *Response, error)

	Posts(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error)
	PostsOf(ctx context.Context, username string, opts ...SearchOptionSetter) (*Links, *Response, error)

	Comments(ctx context.Context, opts ...SearchOptionSetter) (*Comments, *Response, error)
	CommentsOf(ctx context.Context, username string, opts ...SearchOptionSetter) (*Comments, *Response, error)

	Saved(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Comments, *Response, error)
	Upvoted(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error)
	Downvoted(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error)
	Hidden(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error)
	Gilded(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error)

	Friend(ctx context.Context, username string, note string) (interface{}, *Response, error)
	Unblock(ctx context.Context, username string) (*Response, error)
	Unfriend(ctx context.Context, username string) (*Response, error)
}

// UserServiceOp implements the UserService interface
type UserServiceOp struct {
	client *Client
}

var _ UserService = &UserServiceOp{}

// User represents a Reddit user
type User struct {
	// this is not the full ID, watch out
	ID      string     `json:"id,omitempty"`
	Name    string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	LinkKarma    int `json:"link_karma"`
	CommentKarma int `json:"comment_karma"`

	IsFriend         bool `json:"is_friend"`
	IsEmployee       bool `json:"is_employee"`
	HasVerifiedEmail bool `json:"has_verified_email"`
	NSFW             bool `json:"over_18"`
	IsSuspended      bool `json:"is_suspended"`
}

// UserShort represents a Reddit user, but contains fewer pieces of information
// It is returned from the GET /api/user_data_by_account_ids endpoint
type UserShort struct {
	Name    string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	LinkKarma    int `json:"link_karma"`
	CommentKarma int `json:"comment_karma"`

	NSFW bool `json:"profile_over_18"`
}

// Get returns information about the user
func (s *UserServiceOp) Get(ctx context.Context, username string) (*User, *Response, error) {
	path := fmt.Sprintf("user/%s/about", username)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(userRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, nil
}

// GetMultipleByID returns multiple users from their full IDs
// The response body is a map where the keys are the IDs (if they exist), and the value is the user
func (s *UserServiceOp) GetMultipleByID(ctx context.Context, ids ...string) (map[string]*UserShort, *Response, error) {
	type query struct {
		IDs []string `url:"ids,omitempty,comma"`
	}

	path := "api/user_data_by_account_ids"
	path, err := addOptions(path, query{ids})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(map[string]*UserShort)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, nil
}

// UsernameAvailable checks whether a username is available for registration
// If a valid username is provided, this endpoint returns a body with just "true" or "false"
func (s *UserServiceOp) UsernameAvailable(ctx context.Context, username string) (bool, *Response, error) {
	type query struct {
		User string `url:"user,omitempty"`
	}

	path := "api/username_available"
	path, err := addOptions(path, query{username})
	if err != nil {
		return false, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return false, nil, err
	}

	root := new(bool)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return false, resp, err
	}

	return *root, resp, nil
}

// Overview returns a list of the client's comments and links
func (s *UserServiceOp) Overview(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Comments, *Response, error) {
	return s.OverviewOf(ctx, s.client.Username, opts...)
}

// OverviewOf returns a list of the user's comments and links
func (s *UserServiceOp) OverviewOf(ctx context.Context, username string, opts ...SearchOptionSetter) (*Links, *Comments, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/overview", username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, resp, err
	}

	return root.getLinks(), root.getComments(), resp, nil
}

// Posts returns a list of the client's posts.
func (s *UserServiceOp) Posts(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error) {
	return s.PostsOf(ctx, s.client.Username, opts...)
}

// PostsOf returns a list of the user's posts.
func (s *UserServiceOp) PostsOf(ctx context.Context, username string, opts ...SearchOptionSetter) (*Links, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/submitted", username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// Comments returns a list of the client's comments.
func (s *UserServiceOp) Comments(ctx context.Context, opts ...SearchOptionSetter) (*Comments, *Response, error) {
	return s.CommentsOf(ctx, s.client.Username, opts...)
}

// CommentsOf returns a list of the user's comments.
func (s *UserServiceOp) CommentsOf(ctx context.Context, username string, opts ...SearchOptionSetter) (*Comments, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/comments", username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getComments(), resp, nil
}

// Saved returns a list of the user's saved posts and comments.
func (s *UserServiceOp) Saved(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Comments, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/saved", s.client.Username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, resp, err
	}

	return root.getLinks(), root.getComments(), resp, nil
}

// Upvoted returns a list of the user's upvoted posts.
func (s *UserServiceOp) Upvoted(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/upvoted", s.client.Username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// Downvoted returns a list of the user's downvoted posts.
func (s *UserServiceOp) Downvoted(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/downvoted", s.client.Username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// Hidden returns a list of the user's hidden posts.
func (s *UserServiceOp) Hidden(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/hidden", s.client.Username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// Gilded returns a list of the user's gilded posts.
func (s *UserServiceOp) Gilded(ctx context.Context, opts ...SearchOptionSetter) (*Links, *Response, error) {
	form := newSearchOptions(opts...)

	path := fmt.Sprintf("user/%s/gilded", s.client.Username)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// // Friend creates or updates a "friend" relationship
// // Request body contains JSON data with:
// //   name: existing Reddit username
// //   note: a string no longer than 300 characters
// func (s *UserServiceOp) Friend(ctx context.Context, username string, note string) (interface{}, *Response, error) {
// 	type request struct {
// 		Username string `url:"name"`
// 		Note     string `url:"note"`
// 	}

// 	path := fmt.Sprintf("api/v1/me/friends/%s", username)
// 	body := request{Username: username, Note: note}

// 	_, err := s.client.NewRequest(http.MethodPut, path, body)
// 	if err != nil {
// 		return false, nil, err
// 	}

// 	// todo: requires gold
// 	return nil, nil, nil
// }

// Friend creates or updates a "friend" relationship
// Request body contains JSON data with:
//   name: existing Reddit username
//   note: a string no longer than 300 characters
func (s *UserServiceOp) Friend(ctx context.Context, username string, note string) (interface{}, *Response, error) {
	type request struct {
		Username string `url:"name"`
		Note     string `url:"note"`
	}

	path := fmt.Sprintf("api/v1/me/friends/%s", username)
	body := request{Username: username, Note: note}

	_, err := s.client.NewRequest(http.MethodPut, path, body)
	if err != nil {
		return false, nil, err
	}

	// todo: requires gold
	return nil, nil, nil
}

// Unblock unblocks a user
func (s *UserServiceOp) Unblock(ctx context.Context, username string) (*Response, error) {
	path := "api/unfriend"

	form := url.Values{}
	form.Set("name", username)
	form.Set("type", "enemy")
	form.Set("container", "todo: this should be the current user's full id")

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unfriend unfriends a user
func (s *UserServiceOp) Unfriend(ctx context.Context, username string) (*Response, error) {
	path := fmt.Sprintf("api/v1/me/friends/%s", username)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
