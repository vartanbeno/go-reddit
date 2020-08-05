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
type SubredditService struct {
	client *Client
}

type rootSubreddit struct {
	Kind string     `json:"kind,omitempty"`
	Data *Subreddit `json:"data,omitempty"`
}

type rootSubredditNames struct {
	Names []string `json:"names,omitempty"`
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

func (s *SubredditService) getPosts(ctx context.Context, sort string, subreddit string, opts *ListPostOptions) (*Posts, *Response, error) {
	path := sort
	if subreddit != "" {
		path = fmt.Sprintf("r/%s/%s", subreddit, sort)
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

	return root.getPosts(), resp, nil
}

// HotPosts returns the hottest posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
// Note: when looking for hot posts in a subreddit, it will include the stickied
// posts (if any) PLUS posts from the limit parameter (25 by default).
func (s *SubredditService) HotPosts(ctx context.Context, subreddit string, opts *ListPostOptions) (*Posts, *Response, error) {
	return s.getPosts(ctx, "hot", subreddit, opts)
}

// NewPosts returns the newest posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
func (s *SubredditService) NewPosts(ctx context.Context, subreddit string, opts *ListPostOptions) (*Posts, *Response, error) {
	return s.getPosts(ctx, "new", subreddit, opts)
}

// RisingPosts returns the rising posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
func (s *SubredditService) RisingPosts(ctx context.Context, subreddit string, opts *ListPostOptions) (*Posts, *Response, error) {
	return s.getPosts(ctx, "rising", subreddit, opts)
}

// ControversialPosts returns the most controversial posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
func (s *SubredditService) ControversialPosts(ctx context.Context, subreddit string, opts *ListPostOptions) (*Posts, *Response, error) {
	return s.getPosts(ctx, "controversial", subreddit, opts)
}

// TopPosts returns the top posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
func (s *SubredditService) TopPosts(ctx context.Context, subreddit string, opts *ListPostOptions) (*Posts, *Response, error) {
	return s.getPosts(ctx, "top", subreddit, opts)
}

// Get gets a subreddit by name.
func (s *SubredditService) Get(ctx context.Context, name string) (*Subreddit, *Response, error) {
	if name == "" {
		return nil, nil, errors.New("name: cannot be empty")
	}

	path := fmt.Sprintf("r/%s/about", name)
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

// Popular returns popular subreddits.
func (s *SubredditService) Popular(ctx context.Context, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/popular", opts)
}

// New returns new subreddits.
func (s *SubredditService) New(ctx context.Context, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/new", opts)
}

// Gold returns gold subreddits (i.e. only accessible to users with gold).
// It seems like it returns an empty list if you don't have gold.
func (s *SubredditService) Gold(ctx context.Context, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/gold", opts)
}

// Default returns default subreddits.
func (s *SubredditService) Default(ctx context.Context, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/default", opts)
}

// Subscribed returns the list of subreddits you are subscribed to.
func (s *SubredditService) Subscribed(ctx context.Context, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/subscriber", opts)
}

// Approved returns the list of subreddits you are an approved user in.
func (s *SubredditService) Approved(ctx context.Context, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// Moderated returns the list of subreddits you are a moderator of.
func (s *SubredditService) Moderated(ctx context.Context, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/moderator", opts)
}

// GetSticky1 returns the first stickied post on a subreddit (if it exists).
func (s *SubredditService) GetSticky1(ctx context.Context, name string) (*PostAndComments, *Response, error) {
	return s.getSticky(ctx, name, 1)
}

// GetSticky2 returns the second stickied post on a subreddit (if it exists).
func (s *SubredditService) GetSticky2(ctx context.Context, name string) (*PostAndComments, *Response, error) {
	return s.getSticky(ctx, name, 2)
}

func (s *SubredditService) handleSubscription(ctx context.Context, form url.Values) (*Response, error) {
	path := "api/subscribe"

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
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

// Search searches for subreddits.
func (s *SubredditService) Search(ctx context.Context, query string, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
	path := "subreddits/search"
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	type params struct {
		Query string `url:"q"`
	}
	path, err = addOptions(path, params{query})
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

	return root.getSubreddits(), resp, nil
}

// SearchNames searches for subreddits with names beginning with the query provided.
func (s *SubredditService) SearchNames(ctx context.Context, query string) ([]string, *Response, error) {
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

// SearchPosts searches for posts in the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If no subreddit is provided, the search is run against r/all.
func (s *SubredditService) SearchPosts(ctx context.Context, query string, subreddit string, opts *ListPostSearchOptions) (*Posts, *Response, error) {
	if subreddit == "" {
		subreddit = "all"
	}

	path := fmt.Sprintf("r/%s/search", subreddit)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	type params struct {
		Query              string `url:"q"`
		RestrictSubreddits bool   `url:"restrict_sr,omitempty"`
	}

	notAll := !strings.EqualFold(subreddit, "all")
	path, err = addOptions(path, params{query, notAll})
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

	return root.getPosts(), resp, nil
}

func (s *SubredditService) getSubreddits(ctx context.Context, path string, opts *ListSubredditOptions) (*Subreddits, *Response, error) {
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

	return root.getSubreddits(), resp, nil
}

// getSticky returns one of the 2 stickied posts of the subreddit (if they exist).
// Num should be equal to 1 or 2, depending on which one you want.
func (s *SubredditService) getSticky(ctx context.Context, subreddit string, num int) (*PostAndComments, *Response, error) {
	type query struct {
		Num int `url:"num"`
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

	root := new(PostAndComments)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
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

// todo: sr_detail's NSFW indicator is over_18 instead of over18
func (s *SubredditService) random(ctx context.Context, nsfw bool) (*Subreddit, *Response, error) {
	path := "r/random"
	if nsfw {
		path = "r/randnsfw"
	}

	type query struct {
		ExpandSubreddit bool `url:"sr_detail"`
		Limit           int  `url:"limit,omitempty"`
	}

	path, err := addOptions(path, query{true, 1})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	type rootResponse struct {
		Data struct {
			Children []struct {
				Data struct {
					Subreddit *Subreddit `json:"sr_detail"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	root := new(rootResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	var sr *Subreddit
	if len(root.Data.Children) > 0 {
		sr = root.Data.Children[0].Data.Subreddit
	}

	return sr, resp, nil
}

// Random returns a random SFW subreddit.
func (s *SubredditService) Random(ctx context.Context) (*Subreddit, *Response, error) {
	return s.random(ctx, false)
}

// RandomNSFW returns a random NSFW subreddit.
func (s *SubredditService) RandomNSFW(ctx context.Context) (*Subreddit, *Response, error) {
	return s.random(ctx, true)
}

// SubmissionText gets the submission text for the subreddit.
// This text is set by the subreddit moderators and intended to be displayed on the submission form.
func (s *SubredditService) SubmissionText(ctx context.Context, name string) (string, *Response, error) {
	if name == "" {
		return "", nil, errors.New("name: cannot be empty")
	}

	path := fmt.Sprintf("r/%s/api/submit_text", name)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", nil, err
	}

	type response struct {
		Text string `json:"submit_text"`
	}
	root := new(response)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}

	return root.Text, resp, err
}
