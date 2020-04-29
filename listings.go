package geddit

import (
	"context"
	"encoding/json"
	"net/http"
)

// ListingsService handles communication with the link (post)
// related methods of the Reddit API
type ListingsService interface {
	Get(ctx context.Context, ids ...string) (*Listing, *Response, error)
}

// ListingsServiceOp implements the Vote interface
type ListingsServiceOp struct {
	client *Client
}

var _ ListingsService = &ListingsServiceOp{}

type listingRoot struct {
	Kind string `json:"kind,omitempty"`
	Data *struct {
		Dist     int                      `json:"dist"`
		Children []map[string]interface{} `json:"children,omitempty"`
		After    string                   `json:"after,omitempty"`
		Before   string                   `json:"before,omitempty"`
	} `json:"data,omitempty"`
}

// Listing holds various types of things that all come from the Reddit API
type Listing struct {
	Links      []*Submission `json:"links,omitempty"`
	Comments   []*Comment    `json:"comments,omitempty"`
	Subreddits []*Subreddit  `json:"subreddits,omitempty"`
}

// Get gets a list of things based on their IDs
// Only links, comments, and subreddits are allowed
func (s *ListingsServiceOp) Get(ctx context.Context, ids ...string) (*Listing, *Response, error) {
	type query struct {
		IDs []string `url:"id,comma"`
	}

	path := "api/info"
	path, err := addOptions(path, query{ids})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(listingRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.Data == nil {
		return nil, resp, nil
	}

	l := new(Listing)

	for _, result := range root.Data.Children {
		kind, ok1 := result["kind"].(string)
		data, ok2 := result["data"]

		if ok1 && ok2 {
			byteValue, err := json.Marshal(data)
			if err != nil {
				return nil, resp, err
			}

			var v interface{}
			switch kind {
			case kindComment:
				v = new(Comment)
				l.Comments = append(l.Comments, v.(*Comment))
			case kindLink:
				v = new(Submission)
				l.Links = append(l.Links, v.(*Submission))
			case kindSubreddit:
				v = new(Subreddit)
				l.Subreddits = append(l.Subreddits, v.(*Subreddit))
			default:
				continue
			}

			err = json.Unmarshal(byteValue, v)
			if err != nil {
				return nil, resp, err
			}
		}
	}

	return l, resp, nil
}

// todo: do by_id next
