package geddit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserService handles communication with the user
// related methods of the Reddit API
type UserService interface {
	Get(ctx context.Context, username string) (*User, *Response, error)
	GetMultipleByID(ctx context.Context, ids ...string) (map[string]*UserShort, *Response, error)
	UsernameAvailable(ctx context.Context, username string) (bool, *Response, error)

	GetComments(ctx context.Context, sort sort, opts *ListOptions) (*CommentList, *Response, error)
	GetCommentsOf(ctx context.Context, username string, sort sort, opts *ListOptions) (*CommentList, *Response, error)
}

// UserServiceOp implements the UserService interface
type UserServiceOp struct {
	client *Client
}

var _ UserService = &UserServiceOp{}

type userRoot struct {
	Kind *string `json:"kind,omitempty"`
	Data *User   `json:"data,omitempty"`
}

// User represents a Reddit user
type User struct {
	ID         string  `json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`

	LinkKarma    int `json:"link_karma"`
	CommentKarma int `json:"comment_karma"`

	IsFriend         bool `json:"is_friend"`
	IsEmployee       bool `json:"is_employee"`
	HasVerifiedEmail bool `json:"has_verified_email"`
	NSFW             bool `json:"over_18"`
}

// UserShort represents a Reddit user, but contains fewer pieces of information
// It is returned from the GET /api/user_data_by_account_ids endpoint
type UserShort struct {
	Name       string  `json:"name,omitempty"`
	CreatedUTC float64 `json:"created_utc"`

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

// GetComments returns a list of the client's comments
func (s *UserServiceOp) GetComments(ctx context.Context, sort sort, opts *ListOptions) (*CommentList, *Response, error) {
	return s.GetCommentsOf(ctx, s.client.Username, sort, opts)
}

// GetCommentsOf returns a list of the user's comments
func (s *UserServiceOp) GetCommentsOf(ctx context.Context, username string, sort sort, opts *ListOptions) (*CommentList, *Response, error) {
	path := fmt.Sprintf("user/%s/comments?sort=%s", username, sorts[sort])
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(commentRootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.Data == nil {
		return nil, resp, nil
	}

	cl := new(CommentList)
	var comments []Comment

	for _, child := range root.Data.Roots {
		comments = append(comments, *child.Data)
	}
	cl.Comments = comments
	cl.After = root.Data.After
	cl.Before = root.Data.Before

	return cl, resp, nil
}

// UsernameAvailable checks whether a username is available for registration
// If a valid username is provided, this endpoint returns a body with just "true" or "false"
// If an invalid username is provided, it returns the following:
/*
{
	"json": {
		"errors": [
			[
				"BAD_USERNAME",
				"invalid username",
				"user"
			]
		]
	}
}
*/
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
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return false, resp, fmt.Errorf("must provide username conforming to Reddit's criteria; username provided: %q", username)
	}
	if err != nil {
		return false, resp, err
	}

	return *root, resp, nil
}

// Friend creates or updates a "friend" relationship
// Request body contains JSON data with:
//   name: existing Reddit username
//   note: a string no longer than 300 characters
func (s *UserServiceOp) Friend(ctx context.Context, username string, note string) (interface{}, *Response, error) {
	type request struct {
		Name string `url:"name"`
		Note string `url:"note"`
	}

	path := fmt.Sprintf("api/v1/me/friends/%s", username)
	body := request{Name: username, Note: note}

	_, err := s.client.NewRequest(http.MethodPut, path, body)
	if err != nil {
		return false, nil, err
	}

	// todo: requires gold
	return nil, nil, nil
}
