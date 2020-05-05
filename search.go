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
	Users(ctx context.Context, q string, opts *ListOptions) (*Users, *Response, error)

	LinksByRelevance(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	LinksByHottest(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	LinksByTop(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	LinksByComments(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error)
	LinksByRelevanceInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)
	LinksByHottestInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)
	LinksByTopInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)
	LinksByCommentsInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error)

	Subreddits(ctx context.Context, q string, opts *ListOptions) (*Subreddits, *Response, error)
	SubredditNames(ctx context.Context, q string) ([]string, *Response, error)
	SubredditInfo(ctx context.Context, q string) ([]SubredditShort, *Response, error)
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

type subredditNamesRoot struct {
	Names []string `json:"names,omitempty"`
}

type subredditShortsRoot struct {
	Subreddits []SubredditShort `json:"subreddits,omitempty"`
}

// SubredditShort represents minimal information about a subreddit
type SubredditShort struct {
	Name        string `json:"name,omitempty"`
	Subscribers int    `json:"subscriber_count"`
	ActiveUsers int    `json:"active_user_count"`
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

// Users searches for users
func (s *SearchServiceOp) Users(ctx context.Context, q string, opts *ListOptions) (*Users, *Response, error) {
	query := newSearchQuery(q, "user", "", opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getUsers(), resp, nil
}

// LinksByRelevance searches for links sorted by relevance to the search query in all of Reddit
func (s *SearchServiceOp) LinksByRelevance(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortRelevance], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// LinksByHottest searches for the hottest links in all of Reddit
func (s *SearchServiceOp) LinksByHottest(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortHot], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// LinksByTop searches for the top links in all of Reddit
func (s *SearchServiceOp) LinksByTop(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortTop], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// LinksByComments searches for links with the highest number of comments in all of Reddit
func (s *SearchServiceOp) LinksByComments(ctx context.Context, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortComments], opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// LinksByRelevanceInSubreddit searches for link sorted by relevance to the search query in the specified subreddit
func (s *SearchServiceOp) LinksByRelevanceInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortRelevance], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// LinksByHottestInSubreddit searches for the hottest links in the specified subreddit
func (s *SearchServiceOp) LinksByHottestInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortHot], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// LinksByTopInSubreddit searches for the top links in the specified subreddit
func (s *SearchServiceOp) LinksByTopInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortTop], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// LinksByCommentsInSubreddit searches for links with the highest number of comments in the specified subreddit
func (s *SearchServiceOp) LinksByCommentsInSubreddit(ctx context.Context, subreddit, q string, opts *ListOptions) (*Links, *Response, error) {
	query := newSearchQuery(q, "link", sorts[sortComments], opts)

	root, resp, err := s.search(ctx, subreddit, query)
	if err != nil {
		return nil, nil, err
	}

	return root.getLinks(), resp, nil
}

// Subreddits searches for subreddits
func (s *SearchServiceOp) Subreddits(ctx context.Context, q string, opts *ListOptions) (*Subreddits, *Response, error) {
	query := newSearchQuery(q, "sr", "", opts)

	root, resp, err := s.search(ctx, "", query)
	if err != nil {
		return nil, resp, err
	}

	return root.getSubreddits(), resp, nil
}

// SubredditNames searches for subreddits with names beginning with the query provided
func (s *SearchServiceOp) SubredditNames(ctx context.Context, q string) ([]string, *Response, error) {
	path := fmt.Sprintf("api/search_reddit_names?query=%s", q)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subredditNamesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Names, resp, nil
}

// SubredditInfo searches for subreddits with names beginning with the query provided
// They hold a bit more info that just the name
func (s *SearchServiceOp) SubredditInfo(ctx context.Context, q string) ([]SubredditShort, *Response, error) {
	path := fmt.Sprintf("api/search_subreddits?query=%s", q)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subredditShortsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Subreddits, resp, nil
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
