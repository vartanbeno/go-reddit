package geddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SubredditService handles communication with the subreddit
// related methods of the Reddit API
type SubredditService interface {
	GetByName(ctx context.Context, subreddit string) (*Subreddit, *Response, error)

	GetPopular(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)
	GetNew(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)
	GetGold(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)
	GetDefault(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)

	GetMineWhereSubscriber(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)
	GetMineWhereContributor(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)
	GetMineWhereModerator(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)
	GetMineWhereStreams(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error)

	GetHotLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error)
	GetBestLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error)
	GetNewLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error)
	GetRisingLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error)
	GetControversialLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error)
	GetTopLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error)

	GetSticky1(ctx context.Context, subreddit string) (*LinkAndComments, *Response, error)
	GetSticky2(ctx context.Context, subreddit string) (*LinkAndComments, *Response, error)

	Subscribe(ctx context.Context, subreddits ...string) (*Response, error)
	SubscribeByID(ctx context.Context, ids ...string) (*Response, error)
	Unsubscribe(ctx context.Context, subreddits ...string) (*Response, error)
	UnsubscribeByID(ctx context.Context, ids ...string) (*Response, error)

	SearchSubredditNames(ctx context.Context, query string) ([]string, *Response, error)
	SearchSubredditInfo(ctx context.Context, query string) ([]SubredditShort, *Response, error)
}

// SubredditServiceOp implements the SubredditService interface
type SubredditServiceOp struct {
	client *Client
}

var _ SubredditService = &SubredditServiceOp{}

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

// GetByName gets a subreddit by name
func (s *SubredditServiceOp) GetByName(ctx context.Context, subreddit string) (*Subreddit, *Response, error) {
	if subreddit == "" {
		return nil, nil, errors.New("empty subreddit name provided")
	}

	path := fmt.Sprintf("r/%s/about", subreddit)
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

// GetPopular returns popular subreddits
func (s *SubredditServiceOp) GetPopular(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/popular", opts)
}

// GetNew returns new subreddits
func (s *SubredditServiceOp) GetNew(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/new", opts)
}

// GetGold returns gold subreddits
func (s *SubredditServiceOp) GetGold(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/gold", opts)
}

// GetDefault returns default subreddits
func (s *SubredditServiceOp) GetDefault(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/default", opts)
}

// GetMineWhereSubscriber returns the list of subreddits the client is subscribed to
func (s *SubredditServiceOp) GetMineWhereSubscriber(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/subscriber", opts)
}

// GetMineWhereContributor returns the list of subreddits the client is a contributor to
func (s *SubredditServiceOp) GetMineWhereContributor(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// GetMineWhereModerator returns the list of subreddits the client is a moderator in
func (s *SubredditServiceOp) GetMineWhereModerator(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// GetMineWhereStreams returns the list of subreddits the client is subscribed to and has hosted videos in
func (s *SubredditServiceOp) GetMineWhereStreams(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// GetHotLinks returns the hot links
// If no subreddit are provided, then it runs the search against all those the client is subscribed to
// IMPORTANT: for subreddits, this will include the stickied posts (if any)
// PLUS the number of posts from the limit parameter (which is 25 by default)
func (s *SubredditServiceOp) GetHotLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error) {
	return s.getLinks(ctx, SortHot, opts, subreddits...)
}

// GetBestLinks returns the best links
// If no subreddit are provided, then it runs the search against all those the client is subscribed to
// IMPORTANT: for subreddits, this will include the stickied posts (if any)
// PLUS the number of posts from the limit parameter (which is 25 by default)
func (s *SubredditServiceOp) GetBestLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error) {
	return s.getLinks(ctx, SortBest, opts, subreddits...)
}

// GetNewLinks returns the new links
// If no subreddit are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetNewLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error) {
	return s.getLinks(ctx, SortNew, opts, subreddits...)
}

// GetRisingLinks returns the rising links
// If no subreddit are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetRisingLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error) {
	return s.getLinks(ctx, SortRising, opts, subreddits...)
}

// GetControversialLinks returns the controversial links
// If no subreddit are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetControversialLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error) {
	return s.getLinks(ctx, SortControversial, opts, subreddits...)
}

// GetTopLinks returns the top links
// If no subreddit are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetTopLinks(ctx context.Context, opts *ListOptions, subreddits ...string) (*Links, *Response, error) {
	return s.getLinks(ctx, SortTop, opts, subreddits...)
}

// GetSticky1 returns the first stickied post on a subreddit (if it exists)
func (s *SubredditServiceOp) GetSticky1(ctx context.Context, name string) (*LinkAndComments, *Response, error) {
	return s.getSticky(ctx, name, sticky1)
}

// GetSticky2 returns the second stickied post on a subreddit (if it exists)
func (s *SubredditServiceOp) GetSticky2(ctx context.Context, name string) (*LinkAndComments, *Response, error) {
	return s.getSticky(ctx, name, sticky2)
}

// Subscribe subscribes to subreddits based on their names
// Returns {} on success
func (s *SubredditServiceOp) Subscribe(ctx context.Context, subreddits ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "sub")
	form.Set("sr_name", strings.Join(subreddits, ","))
	return s.handleSubscription(ctx, form)
}

// SubscribeByID subscribes to subreddits based on their id
// Returns {} on success
func (s *SubredditServiceOp) SubscribeByID(ctx context.Context, ids ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "sub")
	form.Set("sr", strings.Join(ids, ","))
	return s.handleSubscription(ctx, form)
}

// Unsubscribe unsubscribes from subreddits based on their names
// Returns {} on success
func (s *SubredditServiceOp) Unsubscribe(ctx context.Context, subreddits ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "unsub")
	form.Set("sr_name", strings.Join(subreddits, ","))
	return s.handleSubscription(ctx, form)
}

// UnsubscribeByID unsubscribes from subreddits based on their id
// Returns {} on success
func (s *SubredditServiceOp) UnsubscribeByID(ctx context.Context, ids ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "unsub")
	form.Set("sr", strings.Join(ids, ","))
	return s.handleSubscription(ctx, form)
}

// SearchSubredditNames searches for subreddits with names beginning with the query provided
func (s *SubredditServiceOp) SearchSubredditNames(ctx context.Context, query string) ([]string, *Response, error) {
	path := fmt.Sprintf("api/search_reddit_names?query=%s", query)

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

// SearchSubredditInfo searches for subreddits with names beginning with the query provided.
// They hold a bit more info that just the name, but still not much.
func (s *SubredditServiceOp) SearchSubredditInfo(ctx context.Context, query string) ([]SubredditShort, *Response, error) {
	path := fmt.Sprintf("api/search_subreddits?query=%s", query)

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

func (s *SubredditServiceOp) handleSubscription(ctx context.Context, form url.Values) (*Response, error) {
	path := "api/subscribe"
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

func (s *SubredditServiceOp) getSubreddits(ctx context.Context, path string, opts *ListOptions) (*Subreddits, *Response, error) {
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

	l := new(Subreddits)

	if root.Data != nil {
		l.Subreddits = root.Data.Things.Subreddits
		l.After = root.Data.After
		l.Before = root.Data.Before
	}

	return l, resp, nil
}

func (s *SubredditServiceOp) getLinks(ctx context.Context, sort Sort, opts *ListOptions, subreddits ...string) (*Links, *Response, error) {
	path := sorts[sort]
	if len(subreddits) > 0 {
		path = fmt.Sprintf("r/%s/%s", strings.Join(subreddits, "+"), sort)
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

	l := new(Links)

	if root.Data != nil {
		l.Links = root.Data.Things.Links
		l.After = root.Data.After
		l.Before = root.Data.Before
	}

	return l, resp, nil
}

// getSticky returns one of the 2 stickied posts of the subreddit (if they exist)
// Num should be equal to 1 or 2, depending on which one you want
// If it's <= 1, it's 1
// If it's >= 2, it's 2
func (s *SubredditServiceOp) getSticky(ctx context.Context, subreddit string, num sticky) (*LinkAndComments, *Response, error) {
	type query struct {
		Num sticky `url:"num"`
	}

	path := fmt.Sprintf("r/%s/about/sticky", subreddit)
	path, err := addOptions(path, query{num})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(LinkAndComments)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
