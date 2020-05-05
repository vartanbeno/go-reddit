package geddit

import (
	"context"
	"fmt"
	"net/http"
)

// SearchService handles communication with the search
// related methods of the Reddit API
// IMPORTANT: for searches to include NSFW results, the
// user must check the following in their preferences:
// "include not safe for work (NSFW) search results in searches"
type SearchService interface {
	SearchUsers(ctx context.Context, q string, opts *ListOptions) (*Users, *Response, error)

	SearchLinksByRelevance(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	SearchLinksByHottest(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	SearchLinksByTop(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	SearchLinksByComments(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	SearchLinksByRelevanceInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)
	SearchLinksByHottestInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)
	SearchLinksByTopInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)
	SearchLinksByCommentsInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)

	SearchSubreddits(ctx context.Context, q string, opts *ListOptions) (*Subreddits, *Response, error)
}

// SearchServiceOp implements the VoteService interface
type SearchServiceOp struct {
	client *Client
}

var _ SearchService = &SearchServiceOp{}

type searchQuery struct {
	ListOptions
	Query string `url:"q,omitempty"`
	Type  string `url:"type,omitempty"`
	Sort  string `url:"sort,omitempty"`
}

func newSearchQuery(query, _type, sort string, opts *ListOptions) *searchQuery {
	if opts == nil {
		opts = &ListOptions{}
	}
	return &searchQuery{
		ListOptions: *opts,
		Query:       query,
		Type:        _type,
		Sort:        sort,
	}
}

// SearchUsers searches for users
func (s *SearchServiceOp) SearchUsers(ctx context.Context, q string, opts *ListOptions) (*Users, *Response, error) {
	query := newSearchQuery(q, "user", "", opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getUsers(), resp, nil
}

// SearchLinksByRelevance searches for link sorted by relevance to the search query
func (s *SearchServiceOp) SearchLinksByRelevance(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortRelevance], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchLinksByHottest searches for the hottest links
func (s *SearchServiceOp) SearchLinksByHottest(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortHot], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchLinksByTop searches for the top links
func (s *SearchServiceOp) SearchLinksByTop(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortTop], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchLinksByComments searches for links with the highest number of comments
func (s *SearchServiceOp) SearchLinksByComments(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortComments], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchLinksByRelevanceInSubreddit searches for link sorted by relevance to the search query in the specified subreddit
func (s *SearchServiceOp) SearchLinksByRelevanceInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortRelevance], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchLinksByHottestInSubreddit searches for the hottest links in the specified subreddit
func (s *SearchServiceOp) SearchLinksByHottestInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortHot], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchLinksByTopInSubreddit searches for the top links in the specified subreddit
func (s *SearchServiceOp) SearchLinksByTopInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortTop], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchLinksByCommentsInSubreddit searches for links with the highest number of comments in the specified subreddit
func (s *SearchServiceOp) SearchLinksByCommentsInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortComments], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// SearchSubreddits searches for subreddits
func (s *SearchServiceOp) SearchSubreddits(ctx context.Context, q string, opts *ListOptions) (*Subreddits, *Response, error) {
	query := newSearchQuery(q, "sr", "", opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, resp, err
	}

	return root.getSubreddits(), resp, nil
}

func (s *SearchServiceOp) search(ctx context.Context, subreddit string, opts *searchQuery) (*rootListing, *Response, error) {
	path := "search"
	if subreddit != "" {
		path = fmt.Sprintf("r/%s/search?restrict_sr=true", subreddit)
	}

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
