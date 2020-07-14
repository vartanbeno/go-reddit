package reddit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CommentService handles communication with the comment
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_links_and_comments
type CommentService service

func (s *CommentService) validateCommentID(id string) error {
	if strings.HasPrefix(id, kindComment+"_") {
		return nil
	}
	return fmt.Errorf("comment id %s does not start with %s_", id, kindComment)
}

// Submit submits a comment as a reply to a post, comment, or message.
// parentID is the full ID of the thing being replied to.
func (s *CommentService) Submit(ctx context.Context, parentID string, text string) (*Comment, *Response, error) {
	path := "api/comment"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("return_rtjson", "true")
	form.Set("parent", parentID)
	form.Set("text", text)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
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

// Edit edits a comment.
func (s *CommentService) Edit(ctx context.Context, id string, text string) (*Comment, *Response, error) {
	path := "api/editusertext"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("return_rtjson", "true")
	form.Set("thing_id", id)
	form.Set("text", text)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
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
