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

	Overview(opts ...SearchOpt) *UserCommentPostSearcher
	OverviewOf(username string, opts ...SearchOpt) *UserCommentPostSearcher

	Posts(opts ...SearchOpt) *UserPostSearcher
	PostsOf(username string, opts ...SearchOpt) *UserPostSearcher

	Comments(opts ...SearchOpt) *UserCommentSearcher
	CommentsOf(username string, opts ...SearchOpt) *UserCommentSearcher

	GetUpvoted(opts ...SearchOpt) *UserPostSearcher
	GetDownvoted(opts ...SearchOpt) *UserPostSearcher
	GetHidden(opts ...SearchOpt) *UserPostSearcher
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
func (s *UserServiceOp) Overview(opts ...SearchOpt) *UserCommentPostSearcher {
	return s.OverviewOf(s.client.Username, opts...)
}

// OverviewOf returns a list of the user's comments and links
func (s *UserServiceOp) OverviewOf(username string, opts ...SearchOpt) *UserCommentPostSearcher {
	sr := new(UserCommentPostSearcher)
	sr.client = s.client
	sr.username = username
	sr.where = "overview"
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// Posts returns a list of the client's posts.
func (s *UserServiceOp) Posts(opts ...SearchOpt) *UserPostSearcher {
	return s.PostsOf(s.client.Username, opts...)
}

// PostsOf returns a list of the user's posts.
func (s *UserServiceOp) PostsOf(username string, opts ...SearchOpt) *UserPostSearcher {
	sr := new(UserPostSearcher)
	sr.client = s.client
	sr.username = username
	sr.where = "submitted"
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// Comments returns a list of the client's comments.
func (s *UserServiceOp) Comments(opts ...SearchOpt) *UserCommentSearcher {
	return s.CommentsOf(s.client.Username, opts...)
}

// CommentsOf returns a list of a user's comments.
func (s *UserServiceOp) CommentsOf(username string, opts ...SearchOpt) *UserCommentSearcher {
	sr := new(UserCommentSearcher)
	sr.client = s.client
	sr.username = username
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// GetUpvoted returns a list of the client's upvoted submissions.
func (s *UserServiceOp) GetUpvoted(opts ...SearchOpt) *UserPostSearcher {
	sr := new(UserPostSearcher)
	sr.client = s.client
	sr.username = s.client.Username
	sr.where = "upvoted"
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// GetDownvoted returns a list of the client's downvoted submissions.
func (s *UserServiceOp) GetDownvoted(opts ...SearchOpt) *UserPostSearcher {
	sr := new(UserPostSearcher)
	sr.client = s.client
	sr.username = s.client.Username
	sr.where = "downvoted"
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// GetHidden returns a list of the client's hidden submissions.
func (s *UserServiceOp) GetHidden(opts ...SearchOpt) *UserPostSearcher {
	sr := new(UserPostSearcher)
	sr.client = s.client
	sr.username = s.client.Username
	sr.where = "hidden"
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// GetSaved returns a list of the client's saved comments and links.
func (s *UserServiceOp) GetSaved(ctx context.Context, opts *ListOptions) (*CommentsLinks, *Response, error) {
	path := fmt.Sprintf("user/%s/saved", s.client.Username)
	return s.getCommentsAndLinks(ctx, path, opts)
}

// GetGilded returns a list of the client's gilded comments and links.
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

// UserPostSearcher finds the posts of a user.
type UserPostSearcher struct {
	clientSearcher
	username string
	// where can be submitted, upvoted, downvoted, hidden
	// https://www.reddit.com/dev/api/#GET_user_{username}_{where}
	where   string
	after   string
	Results []Link
}

func (s *UserPostSearcher) search(ctx context.Context) (*Links, *Response, error) {
	path := fmt.Sprintf("user/%s/%s", s.username, s.where)
	root, resp, err := s.clientSearcher.Do(ctx, path)
	if err != nil {
		return nil, resp, err
	}
	return root.getLinks(), resp, nil
}

// Search runs the searcher.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserPostSearcher) Search(ctx context.Context) (bool, *Response, error) {
	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = root.Links
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// More runs the searcher again and adds to the results.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserPostSearcher) More(ctx context.Context) (bool, *Response, error) {
	if s.after == "" {
		return s.Search(ctx)
	}

	s.setAfter(s.after)

	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = append(s.Results, root.Links...)
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// All runs the searcher until it yields no more results.
// The limit is set to 100, just to make the least amount
// of requests possible. It is reset to its original value after.
func (s *UserPostSearcher) All(ctx context.Context) error {
	limit := s.opts.Limit

	s.setLimit(100)
	defer s.setLimit(limit)

	var ok = true
	var err error

	for ok {
		ok, _, err = s.More(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// UserCommentSearcher finds the comments of a user.
type UserCommentSearcher struct {
	clientSearcher
	username string
	after    string
	Results  []Comment
}

func (s *UserCommentSearcher) search(ctx context.Context) (*Comments, *Response, error) {
	path := fmt.Sprintf("user/%s/comments", s.username)
	root, resp, err := s.clientSearcher.Do(ctx, path)
	if err != nil {
		return nil, resp, err
	}
	return root.getComments(), resp, nil
}

// Search runs the searcher.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserCommentSearcher) Search(ctx context.Context) (bool, *Response, error) {
	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = root.Comments
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// More runs the searcher again and adds to the results.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserCommentSearcher) More(ctx context.Context) (bool, *Response, error) {
	if s.after == "" {
		return s.Search(ctx)
	}

	s.setAfter(s.after)

	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = append(s.Results, root.Comments...)
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// All runs the searcher until it yields no more results.
// The limit is set to 100, just to make the least amount
// of requests possible. It is reset to its original value after.
func (s *UserCommentSearcher) All(ctx context.Context) error {
	limit := s.opts.Limit

	s.setLimit(100)
	defer s.setLimit(limit)

	var ok = true
	var err error

	for ok {
		ok, _, err = s.More(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// UserCommentPostSearcher finds the comments and posts of a user.
type UserCommentPostSearcher struct {
	clientSearcher
	username string
	where    string
	after    string
	Results  struct {
		Comments []Comment `json:"comments"`
		Posts    []Link    `json:"posts"`
	}
}

func (s *UserCommentPostSearcher) search(ctx context.Context) (*Comments, *Links, *Response, error) {
	path := fmt.Sprintf("user/%s/%s", s.username, s.where)
	root, resp, err := s.clientSearcher.Do(ctx, path)
	if err != nil {
		return nil, nil, resp, err
	}
	return root.getComments(), root.getLinks(), resp, nil
}

// Search runs the searcher.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserCommentPostSearcher) Search(ctx context.Context) (bool, *Response, error) {
	rootComments, rootPosts, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results.Comments = rootComments.Comments
	s.Results.Posts = rootPosts.Links
	s.after = rootComments.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// More runs the searcher again and adds to the results.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserCommentPostSearcher) More(ctx context.Context) (bool, *Response, error) {
	if s.after == "" {
		return s.Search(ctx)
	}

	s.setAfter(s.after)

	rootComments, rootPosts, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results.Comments = append(s.Results.Comments, rootComments.Comments...)
	s.Results.Posts = append(s.Results.Posts, rootPosts.Links...)
	s.after = rootComments.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// All runs the searcher until it yields no more results.
// The limit is set to 100, just to make the least amount
// of requests possible. It is reset to its original value after.
func (s *UserCommentPostSearcher) All(ctx context.Context) error {
	limit := s.opts.Limit

	s.setLimit(100)
	defer s.setLimit(limit)

	var ok = true
	var err error

	for ok {
		ok, _, err = s.More(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
