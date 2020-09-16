package reddit

import (
	"context"
	"fmt"
	"net/http"
)

// LiveThreadService handles communication with the live thread
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_live
type LiveThreadService struct {
	client *Client
}

// LiveThread is a thread on Reddit that provides real-time updates.
type LiveThread struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Resources   string `json:"resources,omitempty"`

	State             string `json:"state,omitempty"`
	ViewerCount       int    `json:"viewer_count"`
	ViewerCountFuzzed bool   `json:"viewer_count_fuzzed"`

	WebSocketURL string `json:"websocket_url,omitempty"`

	Announcement bool `json:"is_announcement"`
	NSFW         bool `json:"nsfw"`
}

// Get information about a live thread.
func (s *LiveThreadService) Get(ctx context.Context, id string) (*LiveThread, *Response, error) {
	path := fmt.Sprintf("live/%s/about", id)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	t, _ := root.LiveThread()
	return t, resp, nil
}
