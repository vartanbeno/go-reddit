package reddit

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ListingsService handles communication with the listing
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_listings
type ListingsService struct {
	client *Client
}

// Get posts, comments, and subreddits from their full IDs.
func (s *ListingsService) Get(ctx context.Context, ids ...string) ([]*Post, []*Comment, []*Subreddit, *Response, error) {
	path := fmt.Sprintf("api/info?id=%s", strings.Join(ids, ","))

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

	posts := root.getPosts().Posts
	return posts, resp, nil
}
