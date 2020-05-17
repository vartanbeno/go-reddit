package geddit

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// SearchService handles communication with the search
// related methods of the Reddit API
// IMPORTANT: for searches to include NSFW results, the
// user must check the following in their preferences:
// "include not safe for work (NSFW) search results in searches"
type SearchService interface {
	Posts(query string) *PostSearchBuilder
	Subreddits(query string) *SubredditSearchBuilder
	Users(query string) *UserSearchBuilder
}

// SearchServiceOp implements the VoteService interface
type SearchServiceOp struct {
	client *Client
}

var _ SearchService = &SearchServiceOp{}

// Posts searches for posts.
// By default, it searches for the most relevant posts of all time.
// To change the sorting, use PostSearchBuilder.Sort().
// Possible sort options: relevance, hot, top, new, comments.
// To change the timespan, use PostSearchBuilder.Timespan().
// Possible timespan options: hour, day, week, month, year, all.
func (s *SearchServiceOp) Posts(query string) *PostSearchBuilder {
	b := new(PostSearchBuilder)
	b.client = s.client
	b.opts.Query = query
	b.opts.Type = "link"
	return b.Sort(SortRelevance).Timespan(TimespanAll)
}

// Subreddits searches for subreddits.
func (s *SearchServiceOp) Subreddits(query string) *SubredditSearchBuilder {
	b := new(SubredditSearchBuilder)
	b.client = s.client
	b.opts.Query = query
	b.opts.Type = "sr"
	return b
}

// Users searches for users.
func (s *SearchServiceOp) Users(query string) *UserSearchBuilder {
	b := new(UserSearchBuilder)
	b.client = s.client
	b.opts.Query = query
	b.opts.Type = "user"
	return b
}

type searchOpts struct {
	Query              string `url:"q"`
	Type               string `url:"type,omitempty"`
	After              string `url:"after,omitempty"`
	Before             string `url:"before,omitempty"`
	Limit              int    `url:"limit,omitempty"`
	RestrictSubreddits bool   `url:"restrict_sr,omitempty"`
	Sort               string `url:"sort,omitempty"`
	Timespan           string `url:"t,omitempty"`
}

// PostSearchBuilder helps conducts searches that return posts.
type PostSearchBuilder struct {
	client     *Client
	subreddits []string
	opts       searchOpts
}

// After sets the after option.
func (b *PostSearchBuilder) After(after string) *PostSearchBuilder {
	b.opts.After = after
	return b
}

// Before sets the before option.
func (b *PostSearchBuilder) Before(before string) *PostSearchBuilder {
	b.opts.Before = before
	return b
}

// Limit sets the limit option.
func (b *PostSearchBuilder) Limit(limit int) *PostSearchBuilder {
	b.opts.Limit = limit
	return b
}

// FromSubreddits restricts the search to happen in the specified subreddits only.
func (b *PostSearchBuilder) FromSubreddits(subreddits ...string) *PostSearchBuilder {
	b.subreddits = subreddits
	b.opts.RestrictSubreddits = len(subreddits) > 0
	return b
}

// FromAll runs the search against r/all.
func (b *PostSearchBuilder) FromAll() *PostSearchBuilder {
	return b.FromSubreddits()
}

// Sort sets the sort option.
func (b *PostSearchBuilder) Sort(sort Sort) *PostSearchBuilder {
	b.opts.Sort = sort.String()
	return b
}

// Timespan sets the timespan option.
func (b *PostSearchBuilder) Timespan(timespan Timespan) *PostSearchBuilder {
	b.opts.Timespan = timespan.String()
	return b
}

// Do conducts the search.
func (b *PostSearchBuilder) Do(ctx context.Context) (*Links, *Response, error) {
	path := "search"
	if len(b.subreddits) > 0 {
		path = fmt.Sprintf("r/%s/search", strings.Join(b.subreddits, "+"))
	}

	path, err := addOptions(path, b.opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := b.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := b.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getLinks(), resp, nil
}

// SubredditSearchBuilder helps conducts searches that return subreddits.
type SubredditSearchBuilder struct {
	client *Client
	opts   searchOpts
}

// After sets the after option.
func (b *SubredditSearchBuilder) After(after string) *SubredditSearchBuilder {
	b.opts.After = after
	return b
}

// Before sets the before option.
func (b *SubredditSearchBuilder) Before(before string) *SubredditSearchBuilder {
	b.opts.Before = before
	return b
}

// Limit sets the limit option.
func (b *SubredditSearchBuilder) Limit(limit int) *SubredditSearchBuilder {
	b.opts.Limit = limit
	return b
}

// Do conducts the search.
func (b *SubredditSearchBuilder) Do(ctx context.Context) (*Subreddits, *Response, error) {
	path := "search"
	path, err := addOptions(path, b.opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := b.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := b.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getSubreddits(), resp, nil
}

// UserSearchBuilder helps conducts searches that return posts.
type UserSearchBuilder struct {
	client *Client
	opts   searchOpts
}

// After sets the after option.
func (b *UserSearchBuilder) After(after string) *UserSearchBuilder {
	b.opts.After = after
	return b
}

// Before sets the before option.
func (b *UserSearchBuilder) Before(before string) *UserSearchBuilder {
	b.opts.Before = before
	return b
}

// Limit sets the limit option.
func (b *UserSearchBuilder) Limit(limit int) *UserSearchBuilder {
	b.opts.Limit = limit
	return b
}

// Do conducts the search.
func (b *UserSearchBuilder) Do(ctx context.Context) (*Users, *Response, error) {
	path := "search"
	path, err := addOptions(path, b.opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := b.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := b.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getUsers(), resp, nil
}
