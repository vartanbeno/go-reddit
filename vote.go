package geddit

import (
	"context"
	"fmt"
	"net/url"
)

// VoteService handles communication with the upvote/downvote
// related methods of the Reddit API
type VoteService interface {
	Up(ctx context.Context, id string) (*Response, error)
	Down(ctx context.Context, id string) (*Response, error)
	Remove(ctx context.Context, id string) (*Response, error)
}

// VoteServiceOp implements the Vote interface
type VoteServiceOp struct {
	client *Client
}

var _ VoteService = &VoteServiceOp{}

type vote int

// Reddit interprets -1, 0, 1 as downvote, no vote, and upvote, respectively.
const (
	downvote vote = iota - 1
	novote
	upvote
)

func (s *VoteServiceOp) vote(ctx context.Context, id string, vote vote) (*Response, error) {
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

// Up upvotes a link or a comment
func (s *VoteServiceOp) Up(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, upvote)
}

// Down downvotes a link or a comment
func (s *VoteServiceOp) Down(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, downvote)
}

// Remove removes the user's vote on a link or a comment
func (s *VoteServiceOp) Remove(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, novote)
}
