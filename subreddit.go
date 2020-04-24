package geddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// SubredditService handles communication with the subreddit
// related methods of the Reddit API
type SubredditService interface {
	GetByName(ctx context.Context, name string) (*Subreddit, *Response, error)
}

// SubredditServiceOp implements the SubredditService interface
type SubredditServiceOp struct {
	client *Client
}

var _ SubredditService = &SubredditServiceOp{}

type subredditRoot struct {
	Kind *string    `json:"kind,omitempty"`
	Data *Subreddit `json:"data,omitempty"`
}

// Subreddit holds information about a subreddit
type Subreddit struct {
	ID      *string  `json:"id,omitempty"`
	FullID  *string  `json:"name,omitempty"`
	Created *float64 `json:"created_utc,omitempty"`

	URL                 *string `json:"url,omitempty"`
	DisplayName         *string `json:"display_name,omitempty"`
	DisplayNamePrefixed *string `json:"display_name_prefixed,omitempty"`
	Title               *string `json:"title,omitempty"`
	PublicDescription   *string `json:"public_description,omitempty"`

	Subscribers     *int `json:"subscribers,omitempty"`
	ActiveUserCount *int `json:"active_user_count,omitempty"`
}

// GetByName gets a subreddit by name
func (s *SubredditServiceOp) GetByName(ctx context.Context, name string) (*Subreddit, *Response, error) {
	if name == "" {
		return nil, nil, errors.New("empty subreddit name provided")
	}

	path := fmt.Sprintf("r/%s/about.json", name)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subredditRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, nil
}
