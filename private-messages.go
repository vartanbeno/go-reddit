package reddit

import (
	"context"
	"net/http"
	"net/url"
)

// PrivateMessageService handles communication with the private message
// related methods of the Reddit API
type PrivateMessageService interface {
	BlockUser(ctx context.Context, messageID string) (*Response, error)
}

// PrivateMessageServiceOp implements the PrivateMessageService interface
type PrivateMessageServiceOp struct {
	client *Client
}

var _ PrivateMessageService = &PrivateMessageServiceOp{}

// BlockUser blocks a user based on the ID of the private message
func (s *PrivateMessageServiceOp) BlockUser(ctx context.Context, messageID string) (*Response, error) {
	path := "api/block"

	form := url.Values{}
	form.Set("id", messageID)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, nil
	}

	return s.client.Do(ctx, req, nil)
}
