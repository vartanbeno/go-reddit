package geddit

import (
	"context"
	"fmt"
	"strings"
)

// SearchService handles communication with the search
// related methods of the Reddit API
// IMPORTANT: for searches to include NSFW results, the
// user must check the following in their preferences:
// "include not safe for work (NSFW) search results in searches"
// Note: The "limit" parameter in searches is prone to inconsistent
// behaviour.
type SearchService interface {
	Posts(query string, opts ...SearchOpt) *PostSearcher
	Subreddits(query string, opts ...SearchOpt) *SubredditSearcher
	Users(query string, opts ...SearchOpt) *UserSearcher
}

// SearchServiceOp implements the VoteService interface
type SearchServiceOp struct {
	client *Client
}

var _ SearchService = &SearchServiceOp{}

// Posts searches for posts.
// By default, it searches for the most relevant posts of all time.
func (s *SearchServiceOp) Posts(query string, opts ...SearchOpt) *PostSearcher {
	sr := new(PostSearcher)
	sr.client = s.client
	sr.opts.Query = query
	sr.opts.Type = "link"
	sr.opts.Sort = SortRelevance.String()
	sr.opts.Timespan = TimespanAll.String()
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// Subreddits searches for subreddits.
func (s *SearchServiceOp) Subreddits(query string, opts ...SearchOpt) *SubredditSearcher {
	sr := new(SubredditSearcher)
	sr.client = s.client
	sr.opts.Query = query
	sr.opts.Type = "sr"
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// Users searches for users.
func (s *SearchServiceOp) Users(query string, opts ...SearchOpt) *UserSearcher {
	sr := new(UserSearcher)
	sr.client = s.client
	sr.opts.Query = query
	sr.opts.Type = "user"
	for _, opt := range opts {
		opt(sr)
	}
	return sr
}

// PostSearcher helps conducts searches that return posts.
type PostSearcher struct {
	clientSearcher
	subreddits []string
	after      string
	Results    []Link
}

func (s *PostSearcher) search(ctx context.Context) (*Links, *Response, error) {
	path := "search"
	if len(s.subreddits) > 0 {
		path = fmt.Sprintf("r/%s/search", strings.Join(s.subreddits, "+"))
	}

	root, resp, err := s.clientSearcher.Do(ctx, path)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// Search runs the searcher.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *PostSearcher) Search(ctx context.Context) (bool, *Response, error) {
	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = root.Links
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// More runs the searcher again and adds to the results.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *PostSearcher) More(ctx context.Context) (bool, *Response, error) {
	if s.after == "" {
		return s.Search(ctx)
	}

	s.setAfter(s.after)

	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = append(s.Results, root.Links...)
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// All runs the searcher until it yields no more results.
// The limit is set to 100, just to make the least amount
// of requests possible. It is reset to its original value after.
func (s *PostSearcher) All(ctx context.Context) error {
	limit := s.opts.Limit

	s.setLimit(100)
	defer s.setLimit(limit)

	var ok = true
	var err error

	for ok {
		ok, _, err = s.More(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// SubredditSearcher helps conducts searches that return subreddits.
type SubredditSearcher struct {
	clientSearcher
	after   string
	Results []Subreddit
}

func (s *SubredditSearcher) search(ctx context.Context) (*Subreddits, *Response, error) {
	path := "search"
	root, resp, err := s.clientSearcher.Do(ctx, path)
	if err != nil {
		return nil, resp, err
	}
	return root.getSubreddits(), resp, nil
}

// Search runs the searcher.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *SubredditSearcher) Search(ctx context.Context) (bool, *Response, error) {
	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = root.Subreddits
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// More runs the searcher again and adds to the results.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *SubredditSearcher) More(ctx context.Context) (bool, *Response, error) {
	if s.after == "" {
		return s.Search(ctx)
	}

	s.setAfter(s.after)

	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = append(s.Results, root.Subreddits...)
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// All runs the searcher until it yields no more results.
// The limit is set to 100, just to make the least amount
// of requests possible. It is reset to its original value after.
func (s *SubredditSearcher) All(ctx context.Context) error {
	limit := s.opts.Limit

	s.setLimit(100)
	defer s.setLimit(limit)

	var ok = true
	var err error

	for ok {
		ok, _, err = s.More(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// UserSearcher helps conducts searches that return users.
type UserSearcher struct {
	clientSearcher
	after   string
	Results []User
}

func (s *UserSearcher) search(ctx context.Context) (*Users, *Response, error) {
	path := "search"
	root, resp, err := s.clientSearcher.Do(ctx, path)
	if err != nil {
		return nil, resp, err
	}
	return root.getUsers(), resp, nil
}

// Search runs the searcher.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserSearcher) Search(ctx context.Context) (bool, *Response, error) {
	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = root.Users
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// More runs the searcher again and adds to the results.
// The first return value tells the user if there are
// more results that were cut off (due to the limit).
func (s *UserSearcher) More(ctx context.Context) (bool, *Response, error) {
	if s.after == "" {
		return s.Search(ctx)
	}

	s.setAfter(s.after)

	root, resp, err := s.search(ctx)
	if err != nil {
		return false, resp, err
	}

	s.Results = append(s.Results, root.Users...)
	s.after = root.After

	// if the "after" value is non-empty, it
	// means there are more results to come.
	moreResultsExist := s.after != ""

	return moreResultsExist, resp, nil
}

// All runs the searcher until it yields no more results.
// The limit is set to 100, just to make the least amount
// of requests possible. It is reset to its original value after.
func (s *UserSearcher) All(ctx context.Context) error {
	limit := s.opts.Limit

	s.setLimit(100)
	defer s.setLimit(limit)

	var ok = true
	var err error

	for ok {
		ok, _, err = s.More(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
