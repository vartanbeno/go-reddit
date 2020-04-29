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

type commentRoot struct {
	Kind *string  `json:"kind,omitempty"`
	Data *Comment `json:"data,omitempty"`
}

type commentRootListing struct {
	Kind *string `json:"kind,omitempty"`
	Data *struct {
		Dist   int           `json:"dist"`
		Roots  []commentRoot `json:"children,omitempty"`
		After  string        `json:"after,omitempty"`
		Before string        `json:"before,omitempty"`
	} `json:"data,omitempty"`
}

// Comment is a comment posted by a user
type Comment struct {
	ID        string `json:"id,omitempty"`
	FullID    string `json:"name,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`
	Permalink string `json:"permalink,omitempty"`

	Body            string `json:"body,omitempty"`
	BodyHTML        string `json:"body_html,omitempty"`
	Author          string `json:"author,omitempty"`
	AuthorID        string `json:"author_fullname,omitempty"`
	AuthorFlairText string `json:"author_flair_text,omitempty"`

	Subreddit             string `json:"subreddit,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string `json:"subreddit_id,omitempty"`

	Score            int `json:"score"`
	Controversiality int `json:"controversiality"`

	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`

	LinkID string `json:"link_id,omitempty"`

	// These don't appear when submitting a comment
	LinkTitle       string `json:"link_title,omitempty"`
	LinkPermalink   string `json:"link_permalink,omitempty"`
	LinkAuthor      string `json:"link_author,omitempty"`
	LinkNumComments int    `json:"num_comments"`

	IsSubmitter bool `json:"is_submitter"`
	ScoreHidden bool `json:"score_hidden"`
	Saved       bool `json:"saved"`
	Stickied    bool `json:"stickied"`
	Locked      bool `json:"locked"`
	CanGild     bool `json:"can_gild"`
	NSFW        bool `json:"over_18"`
}

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
