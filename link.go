package geddit

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

// LinkService handles communication with the link (post)
// related methods of the Reddit API
type LinkService interface {
	Hide(ctx context.Context, ids ...string) (*Response, error)
	Unhide(ctx context.Context, ids ...string) (*Response, error)
}

// LinkServiceOp implements the Vote interface
type LinkServiceOp struct {
	client *Client
}

var _ LinkService = &LinkServiceOp{}

// Hide hides links with the specified ids
// On successful calls, it just returns {}
func (s *LinkServiceOp) Hide(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/hide"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

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

// Unhide unhides links with the specified ids
// On successful calls, it just returns {}
func (s *LinkServiceOp) Unhide(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/unhide"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

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
