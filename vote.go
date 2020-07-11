package reddit

import (
	"context"
	"fmt"
	"net/url"
)

// VoteService handles communication with the upvote/downvote
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#POST_api_vote
type VoteService service

type vote int

// Reddit interprets -1, 0, 1 as downvote, no vote, and upvote, respectively.
const (
	downvote vote = iota - 1
	novote
	upvote
)

func (s *VoteService) vote(ctx context.Context, id string, vote vote) (*Response, error) {
	path := "api/vote"

	form := url.Values{}
	form.Set("id", id)
	form.Set("dir", fmt.Sprint(vote))
	form.Set("rank", "10")

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Up upvotes a post or a comment.
func (s *VoteService) Up(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, upvote)
}

// Down downvotes a post or a comment.
func (s *VoteService) Down(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, downvote)
}

// Remove removes the user's vote on a post or a comment.
func (s *VoteService) Remove(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, novote)
}
