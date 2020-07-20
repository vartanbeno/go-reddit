package reddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SubredditService handles communication with the subreddit
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_subreddits
type SubredditService service

type rootSubreddit struct {
	Kind string     `json:"kind,omitempty"`
	Data *Subreddit `json:"data,omitempty"`
}

type rootSubredditNames struct {
	Names []string `json:"names,omitempty"`
}

type rootSubredditShorts struct {
	Subreddits []SubredditShort `json:"subreddits,omitempty"`
}

// SubredditShort represents minimal information about a subreddit
type SubredditShort struct {
	Name        string `json:"name,omitempty"`
	Subscribers int    `json:"subscriber_count"`
	ActiveUsers int    `json:"active_user_count"`
}

type rootModeratorList struct {
	Kind string `json:"kind,omitempty"`
	Data struct {
		Moderators []Moderator `json:"children"`
	} `json:"data"`
}

// Moderator is a user who moderates a subreddit.
type Moderator struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Permissions []string `json:"mod_permissions"`
}

// GetPosts returns posts.
// By default, it'll look for the hottest posts from r/all.
// Note: when looking for hot posts in a subreddit, it will include the
// sticked posts (if any) PLUS posts from the limit parameter (25 by default).
func (s *SubredditService) GetPosts() *PostFinder {
	f := new(PostFinder)
	f.client = s.client
	return f.Sort(SortHot).FromAll()
}

// GetByName gets a subreddit by name
func (s *SubredditService) GetByName(ctx context.Context, subreddit string) (*Subreddit, *Response, error) {
	if subreddit == "" {
		return nil, nil, errors.New("empty subreddit name provided")
	}

	path := fmt.Sprintf("r/%s/about", subreddit)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootSubreddit)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, nil
}

// GetPopular returns popular subreddits.
func (s *SubredditService) GetPopular(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/popular", opts)
}

// GetNew returns new subreddits.
func (s *SubredditService) GetNew(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/new", opts)
}

// GetGold returns gold subreddits.
func (s *SubredditService) GetGold(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/gold", opts)
}

// GetDefault returns default subreddits.
func (s *SubredditService) GetDefault(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/default", opts)
}

// GetSubscribed returns the list of subreddits the client is subscribed to.
func (s *SubredditService) GetSubscribed(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/subscriber", opts)
}

// GetApproved returns the list of subreddits the client is an approved user in.
func (s *SubredditService) GetApproved(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// GetModerated returns the list of subreddits the client is a moderator of.
func (s *SubredditService) GetModerated(ctx context.Context, opts *ListOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/moderator", opts)
}

// GetSticky1 returns the first stickied post on a subreddit (if it exists).
func (s *SubredditService) GetSticky1(ctx context.Context, name string) (*Post, []*Comment, *Response, error) {
	return s.getSticky(ctx, name, 1)
}

// GetSticky2 returns the second stickied post on a subreddit (if it exists).
func (s *SubredditService) GetSticky2(ctx context.Context, name string) (*Post, []*Comment, *Response, error) {
	return s.getSticky(ctx, name, 2)
}

// Subscribe subscribes to subreddits based on their names.
func (s *SubredditService) Subscribe(ctx context.Context, subreddits ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "sub")
	form.Set("sr_name", strings.Join(subreddits, ","))
	return s.handleSubscription(ctx, form)
}

// SubscribeByID subscribes to subreddits based on their id.
func (s *SubredditService) SubscribeByID(ctx context.Context, ids ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "sub")
	form.Set("sr", strings.Join(ids, ","))
	return s.handleSubscription(ctx, form)
}

// Unsubscribe unsubscribes from subreddits based on their names.
func (s *SubredditService) Unsubscribe(ctx context.Context, subreddits ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "unsub")
	form.Set("sr_name", strings.Join(subreddits, ","))
	return s.handleSubscription(ctx, form)
}

// UnsubscribeByID unsubscribes from subreddits based on their id.
func (s *SubredditService) UnsubscribeByID(ctx context.Context, ids ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "unsub")
	form.Set("sr", strings.Join(ids, ","))
	return s.handleSubscription(ctx, form)
}

// SearchSubredditNames searches for subreddits with names beginning with the query provided.
func (s *SubredditService) SearchSubredditNames(ctx context.Context, query string) ([]string, *Response, error) {
	path := fmt.Sprintf("api/search_reddit_names?query=%s", query)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootSubredditNames)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Names, resp, nil
}

// SearchSubredditInfo searches for subreddits with names beginning with the query provided.
// They hold a bit more info that just the name, but still not much.
func (s *SubredditService) SearchSubredditInfo(ctx context.Context, query string) ([]SubredditShort, *Response, error) {
	path := fmt.Sprintf("api/search_subreddits?query=%s", query)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootSubredditShorts)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Subreddits, resp, nil
}

func (s *SubredditService) handleSubscription(ctx context.Context, form url.Values) (*Response, error) {
	path := "api/subscribe"

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SubredditService) getSubreddits(ctx context.Context, path string, opts *ListOptions) (*Subreddits, *Response, error) {
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

// getSticky returns one of the 2 stickied posts of the subreddit (if they exist).
// Num should be equal to 1 or 2, depending on which one you want.
func (s *SubredditService) getSticky(ctx context.Context, subreddit string, num int) (*Post, []*Comment, *Response, error) {
	type query struct {
		Num int `url:"num"`
	}

	path := fmt.Sprintf("r/%s/about/sticky", subreddit)
	path, err := addOptions(path, query{num})
	if err != nil {
		return nil, nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	root := new(postAndComments)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, resp, err
	}

	return root.Post, root.Comments, resp, nil
}

// PostFinder finds posts from the specified subreddits.
// If no subreddits are specified, it finds posts from the ones the client is subscribed to.
type PostFinder struct {
	client     *Client
	subreddits []string
	sort       string
	opts       struct {
		After    string `url:"after,omitempty"`
		Before   string `url:"before,omitempty"`
		Limit    int    `url:"limit,omitempty"`
		Timespan string `url:"t,omitempty"`
	}
}

// After sets the after option.
func (f *PostFinder) After(after string) *PostFinder {
	f.opts.After = after
	return f
}

// Before sets the before option.
func (f *PostFinder) Before(before string) *PostFinder {
	f.opts.Before = before
	return f
}

// Limit sets the limit option.
func (f *PostFinder) Limit(limit int) *PostFinder {
	f.opts.Limit = limit
	return f
}

// FromSubreddits restricts the search to the subreddits.
func (f *PostFinder) FromSubreddits(subreddits ...string) *PostFinder {
	f.subreddits = subreddits
	return f
}

// FromAll allows the finder to find posts from r/all.
func (f *PostFinder) FromAll() *PostFinder {
	f.subreddits = []string{"all"}
	return f
}

// Sort sets the sort option.
func (f *PostFinder) Sort(sort Sort) *PostFinder {
	f.sort = sort.String()
	return f
}

// Timespan sets the timespan option.
func (f *PostFinder) Timespan(timespan Timespan) *PostFinder {
	f.opts.Timespan = timespan.String()
	return f
}

// Do conducts the search.
func (f *PostFinder) Do(ctx context.Context) (*Posts, *Response, error) {
	path := f.sort
	if len(f.subreddits) > 0 {
		path = fmt.Sprintf("r/%s/%s", strings.Join(f.subreddits, "+"), f.sort)
	}

	path, err := addOptions(path, f.opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := f.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := f.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getPosts(), resp, nil
}

// Moderators returns the moderators of a subreddit.
func (s *SubredditService) Moderators(ctx context.Context, subreddit string) (interface{}, *Response, error) {
	path := fmt.Sprintf("r/%s/about/moderators", subreddit)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootModeratorList)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data.Moderators, resp, nil
}
