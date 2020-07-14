package reddit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PostAndCommentService handles communication with the post and comment
// related methods of the Reddit API.
// This service holds functionality common to both posts and comments.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_links_and_comments
type PostAndCommentService service

type vote int

// Reddit interprets -1, 0, 1 as downvote, no vote, and upvote, respectively.
const (
	downvote vote = iota - 1
	novote
	upvote
)

// Delete deletes a post or comment via its full ID.
func (s *PostAndCommentService) Delete(ctx context.Context, id string) (*Response, error) {
	path := "api/del"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Save saves a post or comment.
func (s *PostAndCommentService) Save(ctx context.Context, id string) (*Response, error) {
	path := "api/save"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unsave unsaves a post or comment.
func (s *PostAndCommentService) Unsave(ctx context.Context, id string) (*Response, error) {
	path := "api/unsave"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// EnableReplies enables inbox replies for one of your posts or comments.
func (s *PostAndCommentService) EnableReplies(ctx context.Context, id string) (*Response, error) {
	path := "api/sendreplies"

	form := url.Values{}
	form.Set("id", id)
	form.Set("state", "true")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DisableReplies dsables inbox replies for one of your posts or comments.
func (s *PostAndCommentService) DisableReplies(ctx context.Context, id string) (*Response, error) {
	path := "api/sendreplies"

	form := url.Values{}
	form.Set("id", id)
	form.Set("state", "false")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Lock locks a post or comment, preventing it from receiving new comments.
func (s *PostAndCommentService) Lock(ctx context.Context, id string) (*Response, error) {
	path := "api/lock"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unlock unlocks a post or comment, allowing it to receive new comments.
func (s *PostAndCommentService) Unlock(ctx context.Context, id string) (*Response, error) {
	path := "api/unlock"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *PostAndCommentService) vote(ctx context.Context, id string, vote vote) (*Response, error) {
	path := "api/vote"

	form := url.Values{}
	form.Set("id", id)
	form.Set("dir", fmt.Sprint(vote))
	form.Set("rank", "10")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Upvote upvotes a post or a comment.
func (s *PostAndCommentService) Upvote(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, upvote)
}

// Downvote downvotes a post or a comment.
func (s *PostAndCommentService) Downvote(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, downvote)
}

// RemoveVote removes your vote on a post or a comment.
func (s *PostAndCommentService) RemoveVote(ctx context.Context, id string) (*Response, error) {
	return s.vote(ctx, id, novote)
}
