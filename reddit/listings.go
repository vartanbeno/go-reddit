package reddit

import (
	"context"
	"fmt"
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
	path := "api/info"
	params := struct {
		IDs []string `url:"id,omitempty,comma"`
	}{ids}

	l, resp, err := s.client.getListing(ctx, path, params)
	if err != nil {
		return nil, nil, nil, resp, err
	}

	return l.Posts(), l.Comments(), l.Subreddits(), resp, nil
}

// GetPosts returns posts from their full IDs.
func (s *ListingsService) GetPosts(ctx context.Context, ids ...string) ([]*Post, *Response, error) {
	converted_ids := []string{}
	for _, id := range ids {
		converted_ids = append(converted_ids, "t3_"+id)
	}
	path := fmt.Sprintf("by_id/%s", strings.Join(converted_ids, ","))
	l, resp, err := s.client.getListing(ctx, path, nil)
	if err != nil {
		return nil, resp, err
	}
	return l.Posts(), resp, nil
}

func (s *ListingsService) GetComments(ctx context.Context, ids ...string) ([]*Comment, *Response, error) {
	converted_ids := []string{}
	for _, id := range ids {
		converted_ids = append(converted_ids, "t1_"+id)
	}
	path := fmt.Sprintf("api/info?id=%s", strings.Join(converted_ids, ","))
	l, resp, err := s.client.getListing(ctx, path, nil)
	if err != nil {
		return nil, resp, err
	}
	return l.Comments(), resp, nil
}
