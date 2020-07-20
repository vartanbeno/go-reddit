package reddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ListingsService handles communication with the listing
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_listings
type ListingsService service

// Get returns posts, comments, and subreddits from their IDs.
func (s *ListingsService) Get(ctx context.Context, ids ...string) ([]*Post, []*Comment, []*Subreddit, *Response, error) {
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
func (s *ListingsService) GetPosts(ctx context.Context, ids ...string) ([]*Post, *Response, error) {
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
