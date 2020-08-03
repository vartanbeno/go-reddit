package reddit

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// CommentService handles communication with the comment
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_links_and_comments
type CommentService struct {
	*postAndCommentService
	client *Client
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

// LoadMoreReplies retrieves more replies that were left out when initially fetching the comment.
func (s *CommentService) LoadMoreReplies(ctx context.Context, comment *Comment) (*Response, error) {
	if comment == nil {
		return nil, errors.New("comment: cannot be nil")
	}

	if !comment.hasMore() {
		return nil, nil
	}

	postID := comment.PostID
	commentIDs := comment.Replies.MoreComments.Children

	type query struct {
		PostID  string   `url:"link_id"`
		IDs     []string `url:"children,comma"`
		APIType string   `url:"api_type"`
	}

	path := "api/morechildren"
	path, err := addOptions(path, query{postID, commentIDs, "json"})
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	type rootResponse struct {
		JSON struct {
			Data struct {
				Things Things `json:"things"`
			} `json:"data"`
		} `json:"json"`
	}

	root := new(rootResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	comments := root.JSON.Data.Things.Comments
	for _, c := range comments {
		addCommentToReplies(comment, c)
	}

	comment.Replies.MoreComments = nil
	return resp, nil
}

// addCommentToReplies traverses the comment tree to find the one
// that the 2nd comment is replying to. It then adds it to its replies.
func addCommentToReplies(parent *Comment, comment *Comment) {
	if parent.FullID == comment.ParentID {
		parent.Replies.Comments = append(parent.Replies.Comments, comment)
		return
	}

	for _, reply := range parent.Replies.Comments {
		addCommentToReplies(reply, comment)
	}
}
