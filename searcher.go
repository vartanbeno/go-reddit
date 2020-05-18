package geddit

import (
	"context"
	"net/http"
)

// todo: query parameter "show" = "all"

// Searcher defines some parameters common to all requests
// used to conduct searches against the Reddit API.
type Searcher interface {
	setAfter(string)
	setBefore(string)
	setLimit(int)
	setSort(Sort)
	setTimespan(Timespan)
}

// Contains all options used for searching.
// Not all are used for every search endpoint.
// For example, for getting a user's posts, "q" is not used.
// After/Before are used as the anchor points for subsequent searches.
// Limit is the maximum number of items to be returned (default: 25, max: 100).
// Sort: hot, new, top, controversial, etc.
// Timespan: hour, day, week, month, year, all.
type searchOpts struct {
	Query              string `url:"q,omitempty"`
	Type               string `url:"type,omitempty"`
	After              string `url:"after,omitempty"`
	Before             string `url:"before,omitempty"`
	Limit              int    `url:"limit,omitempty"`
	RestrictSubreddits bool   `url:"restrict_sr,omitempty"`
	Sort               string `url:"sort,omitempty"`
	Timespan           string `url:"t,omitempty"`
}

type clientSearcher struct {
	client *Client
	opts   searchOpts
}

var _ Searcher = &clientSearcher{}

func (s *clientSearcher) setAfter(v string) {
	s.opts.After = v
}

func (s *clientSearcher) setBefore(v string) {
	s.opts.Before = v
}

func (s *clientSearcher) setLimit(v int) {
	s.opts.Limit = v
}

func (s *clientSearcher) setSort(v Sort) {
	s.opts.Sort = v.String()
}

func (s *clientSearcher) setTimespan(v Timespan) {
	s.opts.Timespan = v.String()
}

func (s *clientSearcher) Do(ctx context.Context, path string) (*rootListing, *Response, error) {
	path, err := addOptions(path, s.opts)
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

// SearchOpt sets search options.
type SearchOpt func(s Searcher)

// SetAfter sets the after option.
func SetAfter(v string) SearchOpt {
	return func(s Searcher) {
		s.setAfter(v)
	}
}

// SetBefore sets the before option.
func SetBefore(v string) SearchOpt {
	return func(s Searcher) {
		s.setBefore(v)
	}
}

// SetLimit sets the limit option.
func SetLimit(v int) SearchOpt {
	return func(s Searcher) {
		s.setLimit(v)
	}
}

// SetSort sets the sort option.
func SetSort(v Sort) SearchOpt {
	return func(s Searcher) {
		s.setSort(v)
	}
}

// SetTimespan sets the timespan option.
func SetTimespan(v Timespan) SearchOpt {
	return func(s Searcher) {
		s.setTimespan(v)
	}
}

// FromSubreddits is an option that restricts the
// search to happen in the specified subreddits.
// If none are specified, it's like searching r/all.
// This option is only applicable to the PostSearcher.
func FromSubreddits(subreddits ...string) SearchOpt {
	return func(s Searcher) {
		if ps, ok := s.(*PostSearcher); ok {
			ps.subreddits = subreddits
			ps.opts.RestrictSubreddits = len(subreddits) > 0
		}
	}
}
