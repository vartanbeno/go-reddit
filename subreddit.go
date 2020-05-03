package geddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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

	GetHotLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error)
	GetBestLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error)
	GetNewLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error)
	GetRisingLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error)
	GetControversialLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error)
	GetTopLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error)

	// GetSticky1(ctx context.Context, name string) (interface{}, *Response, error)
	// GetSticky2(ctx context.Context, name string) (interface{}, *Response, error)

	Subscribe(ctx context.Context, names ...string) (*Response, error)
	SubscribeByID(ctx context.Context, ids ...string) (*Response, error)
	Unsubscribe(ctx context.Context, names ...string) (*Response, error)
	UnsubscribeByID(ctx context.Context, ids ...string) (*Response, error)

	StreamLinks(ctx context.Context, names ...string) (<-chan Link, chan<- bool, error)
}

// SubredditServiceOp implements the SubredditService interface
type SubredditServiceOp struct {
	client *Client
}

var _ SubredditService = &SubredditServiceOp{}

// SubredditList holds information about a list of subreddits
// The after and before fields help decide the anchor point for a subsequent
// call that returns a list
type SubredditList struct {
	Subreddits []Subreddit `json:"subreddits,omitempty"`
	After      string      `json:"after,omitempty"`
	Before     string      `json:"before,omitempty"`
}

// LinkList holds information about a list of links
// The after and before fields help decide the anchor point for a subsequent
// call that returns a list
// Note: not to be confused with linked lists
type LinkList struct {
	Links  []Link `json:"submissions,omitempty"`
	After  string `json:"after,omitempty"`
	Before string `json:"before,omitempty"`
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

// GetHotLinks returns the hot links
// If no subreddit names are provided, then it runs the search against all those the client is subscribed to
// IMPORTANT: for subreddits, this will include the stickied posts (if any)
// PLUS the number of posts from the limit parameter (which is 25 by default)
func (s *SubredditServiceOp) GetHotLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error) {
	return s.getLinks(ctx, sortHot, opts, names...)
}

// GetBestLinks returns the best links
// If no subreddit names are provided, then it runs the search against all those the client is subscribed to
// IMPORTANT: for subreddits, this will include the stickied posts (if any)
// PLUS the number of posts from the limit parameter (which is 25 by default)
func (s *SubredditServiceOp) GetBestLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error) {
	return s.getLinks(ctx, sortBest, opts, names...)
}

// GetNewLinks returns the new links
// If no subreddit names are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetNewLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error) {
	return s.getLinks(ctx, sortNew, opts, names...)
}

// GetRisingLinks returns the rising links
// If no subreddit names are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetRisingLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error) {
	return s.getLinks(ctx, sortRising, opts, names...)
}

// GetControversialLinks returns the controversial links
// If no subreddit names are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetControversialLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error) {
	return s.getLinks(ctx, sortControversial, opts, names...)
}

// GetTopLinks returns the top links
// If no subreddit names are provided, then it runs the search against all those the client is subscribed to
func (s *SubredditServiceOp) GetTopLinks(ctx context.Context, opts *ListOptions, names ...string) (*LinkList, *Response, error) {
	return s.getLinks(ctx, sortTop, opts, names...)
}

type sticky int

const (
	sticky1 sticky = iota + 1
	sticky2
)

// // GetSticky1 returns the first stickied post on a subreddit (if it exists)
// func (s *SubredditServiceOp) GetSticky1(ctx context.Context, name string) (interface{}, *Response, error) {
// 	return s.getSticky(ctx, name, sticky1)
// }

// // GetSticky2 returns the second stickied post on a subreddit (if it exists)
// func (s *SubredditServiceOp) GetSticky2(ctx context.Context, name string) (interface{}, *Response, error) {
// 	return s.getSticky(ctx, name, sticky2)
// }

// Subscribe subscribes to subreddits based on their name
// Returns {} on success
func (s *SubredditServiceOp) Subscribe(ctx context.Context, names ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "sub")
	form.Set("sr_name", strings.Join(names, ","))
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

// Unsubscribe unsubscribes from subreddits
// Returns {} on success
func (s *SubredditServiceOp) Unsubscribe(ctx context.Context, names ...string) (*Response, error) {
	form := url.Values{}
	form.Set("action", "unsub")
	form.Set("sr_name", strings.Join(names, ","))
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

func (s *SubredditServiceOp) getSubreddits(ctx context.Context, path string, opts *ListOptions) (*SubredditList, *Response, error) {
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

	l := new(SubredditList)

	if root.Data != nil {
		l.Subreddits = root.Data.Things.Subreddits
		l.After = root.Data.After
		l.Before = root.Data.Before
	}

	return l, resp, nil
}

func (s *SubredditServiceOp) getLinks(ctx context.Context, sort sort, opts *ListOptions, names ...string) (*LinkList, *Response, error) {
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

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	l := new(LinkList)

	if root.Data != nil {
		l.Links = root.Data.Things.Links
		l.After = root.Data.After
		l.Before = root.Data.Before
	}

	return l, resp, nil
}

// getSticky returns one of the 2 stickied posts of the subreddit
// Num should be equal to 1 or 2, depending on which one you want
// If it's <= 1, it's 1
// If it's >= 2, it's 2
// todo
// func (s *SubredditServiceOp) getSticky(ctx context.Context, name string, num sticky) (interface{}, *Response, error) {
// 	type query struct {
// 		Num sticky `url:"num"`
// 	}

// 	path := fmt.Sprintf("r/%s/about/sticky", name)
// 	path, err := addOptions(path, query{num})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	req, err := s.client.NewRequest(http.MethodGet, path, nil)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	var root []rootListing
// 	resp, err := s.client.Do(ctx, req, &root)
// 	if err != nil {
// 		return nil, resp, err
// 	}

// 	// test, _ := json.MarshalIndent(root, "", "  ")
// 	// fmt.Println(string(test))

// 	linkRoot := new(linkRoot)

// 	link := root[0].Data.Children[0]
// 	byteValue, err := json.Marshal(link)
// 	if err != nil {
// 		return nil, resp, err
// 	}

// 	err = json.Unmarshal(byteValue, linkRoot)
// 	if err != nil {
// 		return nil, resp, err
// 	}

// 	// these are all the comments in the post
// 	comments := root[1].Data.Children

// 	var commentsRoot []commentRoot
// 	byteValue, err = json.Marshal(comments)
// 	if err != nil {
// 		return nil, resp, err
// 	}

// 	err = json.Unmarshal(byteValue, &commentsRoot)
// 	if err != nil {
// 		return nil, resp, err
// 	}

// 	test, _ := json.MarshalIndent(commentsRoot, "", "  ")
// 	fmt.Println(string(test))

// 	for _, comment := range commentsRoot {
// 		if string(comment.Data.RepliesRaw) == `""` {
// 			comment.Data.Replies = nil
// 			continue
// 		}

// 		// var
// 	}

// 	return commentsRoot, resp, nil
// }

// func handleComments(comments []commentRoot) {
// 	for _, comment := range comments {
// 		if string(comment.Data.RepliesRaw) == `""` {
// 			comment.Data.Replies = nil
// 			continue
// 		}

// 	}
// }

// StreamLinks returns a channel that receives new submissions from the subreddits
// To stop the stream, simply send a bool value to the stop channel
func (s *SubredditServiceOp) StreamLinks(ctx context.Context, names ...string) (<-chan Link, chan<- bool, error) {
	if len(names) == 0 {
		return nil, nil, errors.New("must specify at least one subreddit")
	}

	submissionCh := make(chan Link)
	stop := make(chan bool, 1)

	go func() {
		// todo: if the post with the before gets deleted, you keep getting 0 posts
		var last *Timestamp
		for {
			select {
			case <-stop:
				close(submissionCh)
				return
			default:
				sl, _, err := s.GetNewLinks(ctx, nil, names...)
				if err != nil {
					continue
				}

				var newest *Timestamp
				for i, submission := range sl.Links {
					if i == 0 {
						newest = submission.Created
					}
					if last == nil {
						submissionCh <- submission
						continue
					}
					if last.Before(*submission.Created) {
						submissionCh <- submission
					}
				}
				last = newest
			}
			<-time.After(time.Second * 3)
			fmt.Println()
		}
	}()

	// go func() {
	// 	var before string
	// 	for {
	// 		select {
	// 		case <-stop:
	// 			close(submissionCh)
	// 			return
	// 		default:
	// 			sl, _, err := s.GetSubmissions(ctx, SortNew, &ListOptions{Before: before}, names...)
	// 			if err != nil {
	// 				continue
	// 			}
	// 			fmt.Printf("Received %d posts\n", len(sl.Submissions))

	// 			if len(sl.Submissions) == 0 {
	// 				continue
	// 			}

	// 			for _, submission := range sl.Submissions {
	// 				submissionCh <- submission
	// 			}
	// 			before = sl.Submissions[0].FullID
	// 		}
	// 		<-time.After(time.Second * 5)
	// 		fmt.Println()
	// 	}
	// }()

	return submissionCh, stop, nil
}
