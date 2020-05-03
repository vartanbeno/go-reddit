package geddit

import (
	"context"
	"net/http"
)

// ListingsService handles communication with the link (post)
// related methods of the Reddit API
type ListingsService interface {
	Get(ctx context.Context, ids ...string) (*CommentsLinksSubreddits, *Response, error)
}

// ListingsServiceOp implements the Vote interface
type ListingsServiceOp struct {
	client *Client
}

var _ ListingsService = &ListingsServiceOp{}

// Get gets a list of things based on their IDs
// Only links, comments, and subreddits are allowed
func (s *ListingsServiceOp) Get(ctx context.Context, ids ...string) (*CommentsLinksSubreddits, *Response, error) {
	type query struct {
		IDs []string `url:"id,comma"`
	}

	path := "api/info"
	path, err := addOptions(path, query{ids})
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

	v := new(CommentsLinksSubreddits)
	v.Comments = root.getComments().Comments
	v.Links = root.getLinks().Links
	v.Subreddits = root.getSubreddits().Subreddits

	return v, resp, nil
}

// todo: do by_id next
