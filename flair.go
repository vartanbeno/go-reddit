package reddit

import (
	"context"
	"fmt"
	"net/http"
)

// FlairService handles communication with the flair
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_flair
type FlairService struct {
	client *Client
}

// Flair is a flair on Reddit
type Flair struct {
	ID      string `json:"id,omitempty"`
	Text    string `json:"text,omitempty"`
	Type    string `json:"type,omitempty"`
	CSS     string `json:"css_class,omitempty"`
	ModOnly bool   `json:"mod_only"`
}

// GetFromSubreddit returns the flairs from the subreddit
func (s *FlairService) GetFromSubreddit(ctx context.Context, name string) ([]Flair, *Response, error) {
	path := fmt.Sprintf("r/%s/api/user_flair_v2", name)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var flairs []Flair
	resp, err := s.client.Do(ctx, req, &flairs)
	if err != nil {
		return nil, resp, err
	}

	return flairs, resp, nil
}
