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

	Overview(ctx context.Context, opts *ListOptions) (*CommentsLinks, *Response, error)
	OverviewOf(ctx context.Context, username string, opts *ListOptions) (*CommentsLinks, *Response, error)

	GetPosts() *UserPostFinder
	GetPostsOf(username string) *UserPostFinder

	GetComments() *UserCommentFinder
	GetCommentsOf(username string) *UserCommentFinder

	GetUpvoted() *UserPostFinder
	GetDownvoted() *UserPostFinder
	GetHidden() *UserPostFinder
	GetSaved(ctx context.Context, opts *ListOptions) (*CommentsLinks, *Response, error)
	GetGilded(ctx context.Context, opts *ListOptions) (*CommentsLinks, *Response, error)

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
	// is not the full ID, watch out
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
func (s *UserServiceOp) Overview(ctx context.Context, opts *ListOptions) (*CommentsLinks, *Response, error) {
	return s.OverviewOf(ctx, s.client.Username, opts)
}

// OverviewOf returns a list of the user's comments and links
func (s *UserServiceOp) OverviewOf(ctx context.Context, username string, opts *ListOptions) (*CommentsLinks, *Response, error) {
	path := fmt.Sprintf("user/%s/overview", username)
	return s.getCommentsAndLinks(ctx, path, opts)
}

// GetPosts returns a list of the client's posts.
func (s *UserServiceOp) GetPosts() *UserPostFinder {
	return s.GetPostsOf(s.client.Username)
}

// GetPostsOf returns a list of the user's posts.
func (s *UserServiceOp) GetPostsOf(username string) *UserPostFinder {
	return newUserPostFinder(s.client, username, "submitted")
}

// GetComments returns a list of the client's comments.
func (s *UserServiceOp) GetComments() *UserCommentFinder {
	return s.GetCommentsOf(s.client.Username)
}

// GetCommentsOf returns a list of a user's comments.
func (s *UserServiceOp) GetCommentsOf(username string) *UserCommentFinder {
	f := new(UserCommentFinder)
	f.client = s.client
	f.username = username
	return f
}

// GetUpvoted returns a list of the client's upvoted submissions
func (s *UserServiceOp) GetUpvoted() *UserPostFinder {
	return newUserPostFinder(s.client, s.client.Username, "upvoted")
}

// GetDownvoted returns a list of the client's downvoted submissions
func (s *UserServiceOp) GetDownvoted() *UserPostFinder {
	return newUserPostFinder(s.client, s.client.Username, "downvoted")
}

// GetHidden returns a list of the client's hidden submissions
func (s *UserServiceOp) GetHidden() *UserPostFinder {
	return newUserPostFinder(s.client, s.client.Username, "hidden")
}

// GetSaved returns a list of the client's saved comments and links
func (s *UserServiceOp) GetSaved(ctx context.Context, opts *ListOptions) (*CommentsLinks, *Response, error) {
	path := fmt.Sprintf("user/%s/saved", s.client.Username)
	return s.getCommentsAndLinks(ctx, path, opts)
}

// GetGilded returns a list of the client's gilded comments and links
func (s *UserServiceOp) GetGilded(ctx context.Context, opts *ListOptions) (*CommentsLinks, *Response, error) {
	path := fmt.Sprintf("user/%s/gilded", s.client.Username)
	return s.getCommentsAndLinks(ctx, path, opts)
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

func (s *UserServiceOp) getLinks(ctx context.Context, path string, opts *ListOptions) (*Links, *Response, error) {
	listing, resp, err := s.getListing(ctx, path, opts)
	if err != nil {
		return nil, resp, err
	}
	return listing.getLinks(), resp, nil
}

func (s *UserServiceOp) getComments(ctx context.Context, path string, opts *ListOptions) (*Comments, *Response, error) {
	listing, resp, err := s.getListing(ctx, path, opts)
	if err != nil {
		return nil, resp, err
	}
	return listing.getComments(), resp, nil
}

func (s *UserServiceOp) getCommentsAndLinks(ctx context.Context, path string, opts *ListOptions) (*CommentsLinks, *Response, error) {
	listing, resp, err := s.getListing(ctx, path, opts)
	if err != nil {
		return nil, resp, err
	}

	v := new(CommentsLinks)
	v.Comments = listing.getComments().Comments
	v.Links = listing.getLinks().Links
	v.After = listing.getAfter()
	v.Before = listing.getBefore()

	return v, resp, nil
}

func (s *UserServiceOp) getListing(ctx context.Context, path string, opts *ListOptions) (*rootListing, *Response, error) {
	path, err := addOptions(path, opts)
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

	return root, resp, err
}

// UserPostFinder finds the posts of a user.
type UserPostFinder struct {
	client   *Client
	username string
	// where can be submitted, upvoted, downvoted, hidden
	// https://www.reddit.com/dev/api/#GET_user_{username}_{where}
	where string
	opts  struct {
		After  string `url:"after,omitempty"`
		Before string `url:"before,omitempty"`
		Limit  int    `url:"limit,omitempty"`
		Sort   string `url:"sort,omitempty"`
	}
}

func newUserPostFinder(cli *Client, username string, where string) *UserPostFinder {
	f := new(UserPostFinder)
	f.client = cli
	f.username = username
	f.where = where
	return f
}

// After sets the after option.
func (f *UserPostFinder) After(after string) *UserPostFinder {
	f.opts.After = after
	return f
}

// Before sets the before option.
func (f *UserPostFinder) Before(before string) *UserPostFinder {
	f.opts.Before = before
	return f
}

// Limit sets the limit option.
func (f *UserPostFinder) Limit(limit int) *UserPostFinder {
	f.opts.Limit = limit
	return f
}

// Sort sets the sort option.
func (f *UserPostFinder) Sort(sort Sort) *UserPostFinder {
	f.opts.Sort = sort.String()
	return f
}

// Do conducts the search.
func (f *UserPostFinder) Do(ctx context.Context) (*Links, *Response, error) {
	path := fmt.Sprintf("user/%s/%s", f.username, f.where)
	path, err := addOptions(path, f.opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := f.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := f.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// UserCommentFinder finds the comments of a user.
type UserCommentFinder struct {
	client   *Client
	username string
	opts     struct {
		After  string `url:"after,omitempty"`
		Before string `url:"before,omitempty"`
		Limit  int    `url:"limit,omitempty"`
		Sort   string `url:"sort,omitempty"`
	}
}

// OfUser specified the user we want to get the comments of.
func (f *UserCommentFinder) OfUser(username string) *UserCommentFinder {
	f.username = username
	return f
}

// After sets the after option.
func (f *UserCommentFinder) After(after string) *UserCommentFinder {
	f.opts.After = after
	return f
}

// Before sets the before option.
func (f *UserCommentFinder) Before(before string) *UserCommentFinder {
	f.opts.Before = before
	return f
}

// Limit sets the limit option.
func (f *UserCommentFinder) Limit(limit int) *UserCommentFinder {
	f.opts.Limit = limit
	return f
}

// Sort sets the sort option.
func (f *UserCommentFinder) Sort(sort Sort) *UserCommentFinder {
	f.opts.Sort = sort.String()
	return f
}

// Do conducts the search.
func (f *UserCommentFinder) Do(ctx context.Context) (*Comments, *Response, error) {
	path := fmt.Sprintf("user/%s/comments", f.username)
	path, err := addOptions(path, f.opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := f.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := f.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getComments(), resp, nil
}
