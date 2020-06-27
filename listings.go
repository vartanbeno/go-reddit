package geddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ListingsService handles communication with the post
// related methods of the Reddit API.
type ListingsService interface {
	Get(ctx context.Context, ids ...string) ([]Post, []Comment, []Subreddit, *Response, error)
	GetPosts(ctx context.Context, ids ...string) ([]Post, *Response, error)
	GetPost(ctx context.Context, id string) (*PostAndComments, *Response, error)
}

// ListingsServiceOp implements the Vote interface.
type ListingsServiceOp struct {
	client *Client
}

var _ ListingsService = &ListingsServiceOp{}

// Get returns posts, comments, and subreddits from their IDs.
func (s *ListingsServiceOp) Get(ctx context.Context, ids ...string) ([]Post, []Comment, []Subreddit, *Response, error) {
	type query struct {
		IDs []string `url:"id,comma"`
	}

	path := "api/info"
	path, err := addOptions(path, query{ids})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, nil, resp, err
	}

	posts := root.getPosts().Posts
	comments := root.getComments().Comments
	subreddits := root.getSubreddits().Subreddits

	return posts, comments, subreddits, resp, nil
}

// GetPosts returns posts from their full IDs.
func (s *ListingsServiceOp) GetPosts(ctx context.Context, ids ...string) ([]Post, *Response, error) {
	if len(ids) == 0 {
		return nil, nil, errors.New("must provide at least 1 id")
	}

	path := fmt.Sprintf("by_id/%s", strings.Join(ids, ","))
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getPosts().Posts, resp, nil
}

// GetPost returns a post with its comments.
// The id here is the ID36 of the post, not its full id.
// Example: instead of t3_abc123, use abc123.
func (s *ListingsServiceOp) GetPost(ctx context.Context, id string) (*PostAndComments, *Response, error) {
	path := fmt.Sprintf("comments/%s", id)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(PostAndComments)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
