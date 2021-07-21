package reddit

import (
	"context"
	"fmt"
	"path"
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
	p := "api/info"
	params := struct {
		IDs []string `url:"id,omitempty,comma"`
	}{ids}

	l, resp, err := s.client.getListing(ctx, p, params)
	if err != nil {
		return nil, nil, nil, resp, err
	}

	return l.Posts(), l.Comments(), l.Subreddits(), resp, nil
}

// GetPosts returns posts from their full IDs.
func (s *ListingsService) GetPosts(ctx context.Context, ids ...string) ([]*Post, *Response, error) {
	p := fmt.Sprintf("by_id/%s", strings.Join(ids, ","))
	l, resp, err := s.client.getListing(ctx, p, nil)
	if err != nil {
		return nil, resp, err
	}
	return l.Posts(), resp, nil
}

// Comments returns comments from a subreddit or post.
func (s *ListingsService) Comments(ctx context.Context, subID, postID string, ops *ListOptions) ([]*Comment, *Response, error) {
	p := path.Join("r", subID, postID, "comments")
	l, resp, err := s.client.getListing(ctx, p, ops)
	if err != nil {
		return nil, resp, err
	}
	return l.Comments(), resp, nil
}
