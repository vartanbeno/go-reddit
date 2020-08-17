package reddit

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// MessageService handles communication with the message
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_messages
type MessageService struct {
	client *Client
}

// ReadAll marks all messages/comments as read. It queues up the task on Reddit's end.
// A successful response returns 202 to acknowledge acceptance of the request.
// This endpoint is heavily rate limited.
func (s *MessageService) ReadAll(ctx context.Context) (*Response, error) {
	path := "api/read_all_messages"

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Read marks a message/comment as read via its full ID.
func (s *MessageService) Read(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/read_message"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unread marks a message/comment as unread via its full ID.
func (s *MessageService) Unread(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/unread_message"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Block blocks the author of a thing via the thing's full ID.
// The thing can be a post, comment or message.
func (s *MessageService) Block(ctx context.Context, id string) (*Response, error) {
	path := "api/block"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Collapse collapses messages.
func (s *MessageService) Collapse(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/collapse_message"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Uncollapse uncollapses messages.
func (s *MessageService) Uncollapse(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/uncollapse_message"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Delete deletes a message.
func (s *MessageService) Delete(ctx context.Context, id string) (*Response, error) {
	path := "api/del_msg"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
