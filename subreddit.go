package geddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// SubredditService handles communication with the subreddit
// related methods of the Reddit API
type SubredditService interface {
	GetByName(ctx context.Context, name string) (*Subreddit, *Response, error)

	GetPopular(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)
	GetNew(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)
	GetGold(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)
	GetDefault(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)

	GetMineWhereSubscriber(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)
	GetMineWhereContributor(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)
	GetMineWhereModerator(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)
	GetMineWhereStreams(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error)

	GetHotPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error)
	GetNewPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error)
	GetRisingPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error)
	GetControversialPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error)
	GetTopPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error)

	GetSticky1(ctx context.Context, name string) (interface{}, *Response, error)
	GetSticky2(ctx context.Context, name string) (interface{}, *Response, error)
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

type subredditRootListing struct {
	Kind *string `json:"kind,omitempty"`
	Data *struct {
		Dist   int             `json:"dist"`
		Roots  []subredditRoot `json:"children,omitempty"`
		After  string          `json:"after,omitempty"`
		Before string          `json:"before,omitempty"`
	} `json:"data,omitempty"`
}

// Subreddit holds information about a subreddit
type Subreddit struct {
	ID         string  `json:"id,omitempty"`
	FullID     string  `json:"name,omitempty"`
	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`

	URL                  string `json:"url,omitempty"`
	DisplayName          string `json:"display_name,omitempty"`
	DisplayNamePrefixed  string `json:"display_name_prefixed,omitempty"`
	Title                string `json:"title,omitempty"`
	PublicDescription    string `json:"public_description,omitempty"`
	Type                 string `json:"subreddit_type,omitempty"`
	SuggestedCommentSort string `json:"suggested_comment_sort,omitempty"`

	Subscribers     int  `json:"subscribers"`
	ActiveUserCount *int `json:"active_user_count,omitempty"`
	NSFW            bool `json:"over18"`
	UserIsMod       bool `json:"user_is_moderator"`
}

// SubredditList holds information about a list of subreddits
// The after and before fields help decide the anchor point for a subsequent
// call that returns a list
type SubredditList struct {
	Subreddits []Subreddit `json:"subreddits,omitempty"`
	After      string      `json:"after,omitempty"`
	Before     string      `json:"before,omitempty"`
}

type submissionRoot struct {
	Kind *string     `json:"kind,omitempty"`
	Data *Submission `json:"data,omitempty"`
}

type submissionRootListing struct {
	Kind *string `json:"kind,omitempty"`
	Data *struct {
		Dist   int              `json:"dist"`
		Roots  []submissionRoot `json:"children,omitempty"`
		After  string           `json:"after,omitempty"`
		Before string           `json:"before,omitempty"`
	} `json:"data,omitempty"`
}

// Submission is a submitted post on Reddit
type Submission struct {
	ID         string  `json:"id,omitempty"`
	FullID     string  `json:"name,omitempty"`
	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`

	Permalink string `json:"permalink,omitempty"`
	URL       string `json:"url,omitempty"`

	Title            string `json:"title,omitempty"`
	Body             string `json:"selftext,omitempty"`
	Score            int    `json:"score"`
	NumberOfComments int    `json:"num_comments"`

	SubredditID           string `json:"t5_2qo4s,omitempty"`
	SubredditName         string `json:"subreddit,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`

	AuthorID   string `json:"author_fullname,omitempty"`
	AuthorName string `json:"author,omitempty"`

	Spoiler    bool `json:"spoiler"`
	Locked     bool `json:"locked"`
	NSFW       bool `json:"over_18"`
	IsSelfPost bool `json:"is_self"`
	Saved      bool `json:"saved"`
	Stickied   bool `json:"stickied"`
}

// SubmissionList holds information about a list of subreddits
// The after and before fields help decide the anchor point for a subsequent
// call that returns a list
type SubmissionList struct {
	Submissions []Submission `json:"submissions,omitempty"`
	After       string       `json:"after,omitempty"`
	Before      string       `json:"before,omitempty"`
}

// GetByName gets a subreddit by name
func (s *SubredditServiceOp) GetByName(ctx context.Context, name string) (*Subreddit, *Response, error) {
	if name == "" {
		return nil, nil, errors.New("empty subreddit name provided")
	}

	path := fmt.Sprintf("r/%s/about", name)
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
func (s *SubredditServiceOp) GetPopular(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/popular", opts)
}

// GetNew returns new subreddits
func (s *SubredditServiceOp) GetNew(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/new", opts)
}

// GetGold returns gold subreddits
func (s *SubredditServiceOp) GetGold(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/gold", opts)
}

// GetDefault returns default subreddits
func (s *SubredditServiceOp) GetDefault(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/default", opts)
}

// GetMineWhereSubscriber returns the list of subreddits the client is subscribed to
func (s *SubredditServiceOp) GetMineWhereSubscriber(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/subscriber", opts)
}

// GetMineWhereContributor returns the list of subreddits the client is a contributor to
func (s *SubredditServiceOp) GetMineWhereContributor(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// GetMineWhereModerator returns the list of subreddits the client is a moderator in
func (s *SubredditServiceOp) GetMineWhereModerator(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

// GetMineWhereStreams returns the list of subreddits the client is subscribed to and has hosted videos in
func (s *SubredditServiceOp) GetMineWhereStreams(ctx context.Context, opts *ListOptions) (*SubredditList, *Response, error) {
	return s.getSubreddits(ctx, "subreddits/mine/contributor", opts)
}

type sort int

const (
	sortHot sort = iota
	sortBest
	sortNew
	sortRising
	sortControversial
	sortTop
)

var sorts = [...]string{
	"hot",
	"best",
	"new",
	"rising",
	"controversial",
	"top",
}

// GetHotPosts returns the hot posts
// If no subreddit names are provided, then it runs the search against /r/all
// IMPORTANT: for subreddits, this will include the stickied posts (if any)
// PLUS the number of posts from the limit parameter (which is 25 by default)
func (s *SubredditServiceOp) GetHotPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error) {
	return s.getPosts(ctx, sortHot, opts, names...)
}

// GetBestPosts returns the best posts
// If no subreddit names are provided, then it runs the search against /r/all
// IMPORTANT: for subreddits, this will include the stickied posts (if any)
// PLUS the number of posts from the limit parameter (which is 25 by default)
func (s *SubredditServiceOp) GetBestPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error) {
	return s.getPosts(ctx, sortBest, opts, names...)
}

// GetNewPosts returns the new posts
// If no subreddit names are provided, then it runs the search against /r/all
func (s *SubredditServiceOp) GetNewPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error) {
	return s.getPosts(ctx, sortNew, opts, names...)
}

// GetRisingPosts returns the rising posts
// If no subreddit names are provided, then it runs the search against /r/all
func (s *SubredditServiceOp) GetRisingPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error) {
	return s.getPosts(ctx, sortRising, opts, names...)
}

// GetControversialPosts returns the controversial posts
// If no subreddit names are provided, then it runs the search against /r/all
func (s *SubredditServiceOp) GetControversialPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error) {
	return s.getPosts(ctx, sortControversial, opts, names...)
}

// GetTopPosts returns the top posts
// If no subreddit names are provided, then it runs the search against /r/all
func (s *SubredditServiceOp) GetTopPosts(ctx context.Context, opts *ListOptions, names ...string) (*SubmissionList, *Response, error) {
	return s.getPosts(ctx, sortTop, opts, names...)
}

type sticky int

const (
	sticky1 sticky = iota + 1
	sticky2
)

// GetSticky1 returns the first stickied post on a subreddit (if it exists)
func (s *SubredditServiceOp) GetSticky1(ctx context.Context, name string) (interface{}, *Response, error) {
	return s.getSticky(ctx, name, sticky1)
}

// GetSticky2 returns the second stickied post on a subreddit (if it exists)
func (s *SubredditServiceOp) GetSticky2(ctx context.Context, name string) (interface{}, *Response, error) {
	return s.getSticky(ctx, name, sticky2)
}

func (s *SubredditServiceOp) getSubreddits(ctx context.Context, path string, opts *ListOptions) (*SubredditList, *Response, error) {
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subredditRootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.Data == nil {
		return nil, resp, nil
	}

	sl := new(SubredditList)
	var subreddits []Subreddit

	for _, child := range root.Data.Roots {
		subreddits = append(subreddits, *child.Data)
	}
	sl.Subreddits = subreddits
	sl.After = root.Data.After
	sl.Before = root.Data.Before

	return sl, resp, nil
}

func (s *SubredditServiceOp) getPosts(ctx context.Context, sort sort, opts *ListOptions, names ...string) (*SubmissionList, *Response, error) {
	path := sorts[sort]
	if len(names) > 0 {
		path = fmt.Sprintf("r/%s/%s", strings.Join(names, "+"), sorts[sort])
	}

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(submissionRootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.Data == nil {
		return nil, resp, nil
	}

	sl := new(SubmissionList)
	var submissions []Submission

	for _, child := range root.Data.Roots {
		submissions = append(submissions, *child.Data)
	}
	sl.Submissions = submissions
	sl.After = root.Data.After
	sl.Before = root.Data.Before

	return sl, resp, nil
}

// getSticky returns one of the 2 stickied posts of the subreddit
// Num should be equal to 1 or 2, depending on which one you want
// If it's <= 1, it's 1
// If it's >= 2, it's 2
// todo
func (s *SubredditServiceOp) getSticky(ctx context.Context, name string, num sticky) (interface{}, *Response, error) {
	// type query struct {
	// 	Num sticky `url:"num"`
	// }

	// path := fmt.Sprintf("r/%s/about/sticky", name)
	// path, err := addOptions(path, query{num})
	// if err != nil {
	// 	return nil, nil, err
	// }

	// req, err := s.client.NewRequest(http.MethodGet, path, nil)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// root := new(submissionRootListing)
	// resp, err := s.client.Do(ctx, req, root)
	// if err != nil {
	// 	return nil, resp, err
	// }

	// return nil, resp, nil
	return nil, nil, nil
}
