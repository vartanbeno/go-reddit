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

type rootSubredditNames struct {
	Names []string `json:"names,omitempty"`
}

// Relationship holds information about a relationship (friend/blocked).
type Relationship struct {
	ID      string     `json:"rel_id,omitempty"`
	User    string     `json:"name,omitempty"`
	UserID  string     `json:"id,omitempty"`
	Created *Timestamp `json:"date,omitempty"`
}

// Moderator is a user who moderates a subreddit.
type Moderator struct {
	*Relationship
	Permissions []string `json:"mod_permissions"`
}

// Ban represents a banned relationship.
type Ban struct {
	*Relationship
	// nil means the ban is permanent
	DaysLeft *int   `json:"days_left"`
	Note     string `json:"note,omitempty"`
}

// todo: interface{}, seriously?
func (s *SubredditService) getPosts(ctx context.Context, sort string, subreddit string, opts interface{}) ([]*Post, *Response, error) {
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

	root := new(listing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Posts, resp, nil
}

// HotPosts returns the hottest posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
// To search through all, just specify "all".
// To search through all and filter out subreddits, provide "all-name1-name2".
// Note: when looking for hot posts in a subreddit, it will include the stickied
// posts (if any) PLUS posts from the limit parameter (25 by default).
func (s *SubredditService) HotPosts(ctx context.Context, subreddit string, opts *ListOptions) ([]*Post, *Response, error) {
	return s.getPosts(ctx, "hot", subreddit, opts)
}

// NewPosts returns the newest posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
// To search through all, just specify "all".
// To search through all and filter out subreddits, provide "all-name1-name2".
func (s *SubredditService) NewPosts(ctx context.Context, subreddit string, opts *ListOptions) ([]*Post, *Response, error) {
	return s.getPosts(ctx, "new", subreddit, opts)
}

// RisingPosts returns the rising posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
// To search through all, just specify "all".
// To search through all and filter out subreddits, provide "all-name1-name2".
func (s *SubredditService) RisingPosts(ctx context.Context, subreddit string, opts *ListOptions) ([]*Post, *Response, error) {
	return s.getPosts(ctx, "rising", subreddit, opts)
}

// ControversialPosts returns the most controversial posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
// To search through all, just specify "all".
// To search through all and filter out subreddits, provide "all-name1-name2".
func (s *SubredditService) ControversialPosts(ctx context.Context, subreddit string, opts *ListPostOptions) ([]*Post, *Response, error) {
	return s.getPosts(ctx, "controversial", subreddit, opts)
}

// TopPosts returns the top posts from the specified subreddit.
// To search through multiple, separate the names with a plus (+), e.g. "golang+test".
// If none are defined, it returns the ones from your subscribed subreddits.
// To search through all, just specify "all".
// To search through all and filter out subreddits, provide "all-name1-name2".
func (s *SubredditService) TopPosts(ctx context.Context, subreddit string, opts *ListPostOptions) ([]*Post, *Response, error) {
	return s.getPosts(ctx, "top", subreddit, opts)
}

// Get a subreddit by name.
func (s *SubredditService) Get(ctx context.Context, name string) (*Subreddit, *Response, error) {
	if name == "" {
		return nil, nil, errors.New("name: cannot be empty")
	}

	path := fmt.Sprintf("r/%s/about", name)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	sr, _ := root.Subreddit()
	return sr, resp, nil
}

// Popular returns popular subreddits.
func (s *SubredditService) Popular(ctx context.Context, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/popular", opts)
}

// New returns new subreddits.
func (s *SubredditService) New(ctx context.Context, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/new", opts)
}

// Gold returns gold subreddits (i.e. only accessible to users with gold).
// It seems like it returns an empty list if you don't have gold.
func (s *SubredditService) Gold(ctx context.Context, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/gold", opts)
}

// Default returns default subreddits.
func (s *SubredditService) Default(ctx context.Context, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/default", opts)
}

// Subscribed returns the list of subreddits you are subscribed to.
func (s *SubredditService) Subscribed(ctx context.Context, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/subscriber", opts)
}

// Approved returns the list of subreddits you are an approved user in.
func (s *SubredditService) Approved(ctx context.Context, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// Moderated returns the list of subreddits you are a moderator of.
func (s *SubredditService) Moderated(ctx context.Context, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/moderator", opts)
}

// GetSticky1 returns the first stickied post on a subreddit (if it exists).
func (s *SubredditService) GetSticky1(ctx context.Context, subreddit string) (*PostAndComments, *Response, error) {
	return s.getSticky(ctx, subreddit, 1)
}

// GetSticky2 returns the second stickied post on a subreddit (if it exists).
func (s *SubredditService) GetSticky2(ctx context.Context, subreddit string) (*PostAndComments, *Response, error) {
	return s.getSticky(ctx, subreddit, 2)
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

// Favorite the subreddit.
func (s *SubredditService) Favorite(ctx context.Context, subreddit string) (*Response, error) {
	path := "api/favorite"

	form := url.Values{}
	form.Set("sr_name", subreddit)
	form.Set("make_favorite", "true")
	form.Set("api_type", "json")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unfavorite the subreddit.
func (s *SubredditService) Unfavorite(ctx context.Context, subreddit string) (*Response, error) {
	path := "api/favorite"

	form := url.Values{}
	form.Set("sr_name", subreddit)
	form.Set("make_favorite", "false")
	form.Set("api_type", "json")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Search for subreddits.
func (s *SubredditService) Search(ctx context.Context, query string, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
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

	root := new(listing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Subreddits, resp, nil
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
func (s *SubredditService) SearchPosts(ctx context.Context, query string, subreddit string, opts *ListPostSearchOptions) ([]*Post, *Response, error) {
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

	root := new(listing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Posts, resp, nil
}

func (s *SubredditService) getSubreddits(ctx context.Context, path string, opts *ListSubredditOptions) ([]*Subreddit, *Response, error) {
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(listing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Subreddits, resp, nil
}

// getSticky returns one of the 2 stickied posts of the subreddit (if they exist).
// Num should be equal to 1 or 2, depending on which one you want.
func (s *SubredditService) getSticky(ctx context.Context, subreddit string, num int) (*PostAndComments, *Response, error) {
	type params struct {
		Num int `url:"num"`
	}

	path := fmt.Sprintf("r/%s/about/sticky", subreddit)
	path, err := addOptions(path, params{num})
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

// todo: sr_detail's NSFW indicator is over_18 instead of over18
func (s *SubredditService) random(ctx context.Context, nsfw bool) (*Subreddit, *Response, error) {
	path := "r/random"
	if nsfw {
		path = "r/randnsfw"
	}

	type params struct {
		ExpandSubreddit bool `url:"sr_detail"`
		Limit           int  `url:"limit,omitempty"`
	}

	path, err := addOptions(path, params{true, 1})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Data struct {
			Children []struct {
				Data struct {
					Subreddit *Subreddit `json:"sr_detail"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	})
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

	root := new(struct {
		Text string `json:"submit_text"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}

	return root.Text, resp, err
}

// Banned gets banned users from the subreddit.
func (s *SubredditService) Banned(ctx context.Context, subreddit string, opts *ListOptions) ([]*Ban, *Response, error) {
	path := fmt.Sprintf("r/%s/about/banned", subreddit)

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Data struct {
			Bans   []*Ban `json:"children"`
			After  string `json:"after"`
			Before string `json:"before"`
		} `json:"data"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	resp.After = root.Data.After
	resp.Before = root.Data.Before

	return root.Data.Bans, resp, nil
}

// Muted gets muted users from the subreddit.
func (s *SubredditService) Muted(ctx context.Context, subreddit string, opts *ListOptions) ([]*Relationship, *Response, error) {
	path := fmt.Sprintf("r/%s/about/muted", subreddit)

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Data struct {
			Relationships []*Relationship `json:"children"`
			After         string          `json:"after"`
			Before        string          `json:"before"`
		} `json:"data"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	resp.After = root.Data.After
	resp.Before = root.Data.Before

	return root.Data.Relationships, resp, nil
}

// WikiBanned gets banned users from the subreddit.
func (s *SubredditService) WikiBanned(ctx context.Context, subreddit string, opts *ListOptions) ([]*Ban, *Response, error) {
	path := fmt.Sprintf("r/%s/about/wikibanned", subreddit)

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Data struct {
			Bans   []*Ban `json:"children"`
			After  string `json:"after"`
			Before string `json:"before"`
		} `json:"data"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	resp.After = root.Data.After
	resp.Before = root.Data.Before

	return root.Data.Bans, resp, nil
}

// Contributors gets contributors (also known as approved users) from the subreddit.
func (s *SubredditService) Contributors(ctx context.Context, subreddit string, opts *ListOptions) ([]*Relationship, *Response, error) {
	path := fmt.Sprintf("r/%s/about/contributors", subreddit)

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Data struct {
			Relationships []*Relationship `json:"children"`
			After         string          `json:"after"`
			Before        string          `json:"before"`
		} `json:"data"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	resp.After = root.Data.After
	resp.Before = root.Data.Before

	return root.Data.Relationships, resp, nil
}

// WikiContributors gets contributors of the wiki from the subreddit.
func (s *SubredditService) WikiContributors(ctx context.Context, subreddit string, opts *ListOptions) ([]*Relationship, *Response, error) {
	path := fmt.Sprintf("r/%s/about/wikicontributors", subreddit)

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Data struct {
			Relationships []*Relationship `json:"children"`
			After         string          `json:"after"`
			Before        string          `json:"before"`
		} `json:"data"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	resp.After = root.Data.After
	resp.Before = root.Data.Before

	return root.Data.Relationships, resp, nil
}

// Moderators gets the moderators of the subreddit.
func (s *SubredditService) Moderators(ctx context.Context, subreddit string) ([]*Moderator, *Response, error) {
	path := fmt.Sprintf("r/%s/about/moderators", subreddit)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Data struct {
			Moderators []*Moderator `json:"children"`
		} `json:"data"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data.Moderators, resp, nil
}
