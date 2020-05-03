package geddit

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// CommentService handles communication with the comment
// related methods of the Reddit API
type CommentService interface {
	Submit(ctx context.Context, id string, text string) (*Comment, *Response, error)
	Edit(ctx context.Context, id string, text string) (*Comment, *Response, error)
	Delete(ctx context.Context, id string) (*Response, error)

	Save(ctx context.Context, id string) (*Response, error)
	Unsave(ctx context.Context, id string) (*Response, error)
}

// CommentServiceOp implements the CommentService interface
type CommentServiceOp struct {
	client *Client
}

var _ CommentService = &CommentServiceOp{}

// CommentList holds information about a list of comments
// The after and before fields help decide the anchor point for a subsequent
// call that returns a list
type CommentList struct {
	Comments []Comment `json:"comments,omitempty"`
	After    string    `json:"after,omitempty"`
	Before   string    `json:"before,omitempty"`
}

func (s *CommentServiceOp) isCommentID(id string) bool {
	return strings.HasPrefix(id, kindComment+"_")
}

// Submit submits a comment as a reply to a link or to another comment
func (s *CommentServiceOp) Submit(ctx context.Context, id string, text string) (*Comment, *Response, error) {
	path := "api/comment"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("return_rtjson", "true")
	form.Set("parent", id)
	form.Set("text", text)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(Comment)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Edit edits the comment with the id provided
// todo: don't forget to do this for links (i.e. posts)
func (s *CommentServiceOp) Edit(ctx context.Context, id string, text string) (*Comment, *Response, error) {
	if !s.isCommentID(id) {
		return nil, nil, fmt.Errorf("must provide comment id (starting with t1_); id provided: %q", id)
	}

	path := "api/editusertext"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("return_rtjson", "true")
	form.Set("thing_id", id)
	form.Set("text", text)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(Comment)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Delete deletes a comment via the id
// todo: don't forget to do this for links (i.e. posts)
// Seems like this always returns {} as a response, no matter if an id is even provided
func (s *CommentServiceOp) Delete(ctx context.Context, id string) (*Response, error) {
	if !s.isCommentID(id) {
		return nil, fmt.Errorf("must provide comment id (starting with t1_); id provided: %q", id)
	}

	path := "api/del"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Save saves a comment
// Seems like this just returns {} on success
func (s *CommentServiceOp) Save(ctx context.Context, id string) (*Response, error) {
	if !s.isCommentID(id) {
		return nil, fmt.Errorf("must provide comment id (starting with t1_); id provided: %q", id)
	}

	path := "api/save"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Unsave unsaves a comment
// Seems like this just returns {} on success
func (s *CommentServiceOp) Unsave(ctx context.Context, id string) (*Response, error) {
	if !s.isCommentID(id) {
		return nil, fmt.Errorf("must provide comment id (starting with t1_); id provided: %q", id)
	}

	path := "api/unsave"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
