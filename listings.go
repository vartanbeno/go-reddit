package geddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ListingsService handles communication with the link (post)
// related methods of the Reddit API
type ListingsService interface {
	Get(ctx context.Context, ids ...string) ([]Comment, []Link, []Subreddit, *Response, error)
	GetLinks(ctx context.Context, ids ...string) ([]Link, *Response, error)
	GetLink(ctx context.Context, id string) (*LinkAndComments, *Response, error)
}

// ListingsServiceOp implements the Vote interface
type ListingsServiceOp struct {
	client *Client
}

var _ ListingsService = &ListingsServiceOp{}

// Get returns comments, links, and subreddits from their IDs
func (s *ListingsServiceOp) Get(ctx context.Context, ids ...string) ([]Comment, []Link, []Subreddit, *Response, error) {
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

	comments := root.getComments().Comments
	links := root.getLinks().Links
	subreddits := root.getSubreddits().Subreddits

	return comments, links, subreddits, resp, nil
}

// GetLinks returns links from their full IDs
func (s *ListingsServiceOp) GetLinks(ctx context.Context, ids ...string) ([]Link, *Response, error) {
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

	return root.getLinks().Links, resp, nil
}

// GetLink returns a link with its comments
func (s *ListingsServiceOp) GetLink(ctx context.Context, id string) (*LinkAndComments, *Response, error) {
	path := fmt.Sprintf("comments/%s", id)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(LinkAndComments)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
