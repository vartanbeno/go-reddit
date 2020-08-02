package reddit

import (
	"fmt"
	"net/url"
)

/*

For searches to include NSFW results, the user must
enable the following setting in their preferences:
"include not safe for work (NSFW) search results in searches"
Note: The "limit" parameter in searches is prone to inconsistent
behaviour, e.g. sometimes limit=1 returns nothing when it should.

*/

func newSearchOptions(opts ...SearchOptionSetter) url.Values {
	searchOptions := make(url.Values)
	for _, opt := range opts {
		opt(searchOptions)
	}
	return searchOptions
}

// SearchOptionSetter sets values for the options.
type SearchOptionSetter func(opts url.Values)

// SetAfter sets the after option.
func SetAfter(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("after", v)
	}
}

// SetBefore sets the before option.
func SetBefore(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("before", v)
	}
}

// SetLimit sets the limit option.
// Warning: It seems like setting the limit to 1 sometimes returns 0 results.
func SetLimit(v int) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("limit", fmt.Sprint(v))
	}
}

// SortByHot sets the sort option to return the hottest results first.
func SortByHot(opts url.Values) {
	opts.Set("sort", "hot")
}

// SortByBest sets the sort option to return the best results first.
func SortByBest(opts url.Values) {
	opts.Set("sort", "best")
}

// SortByNew sets the sort option to return the newest results first.
func SortByNew(opts url.Values) {
	opts.Set("sort", "new")
}

// SortByRising sets the sort option to return the rising results first.
func SortByRising(opts url.Values) {
	opts.Set("sort", "rising")
}

// SortByControversial sets the sort option to return the most controversial results first.
func SortByControversial(opts url.Values) {
	opts.Set("sort", "controversial")
}

// SortByTop sets the sort option to return the top results first.
func SortByTop(opts url.Values) {
	opts.Set("sort", "top")
}

// SortByRelevance sets the sort option to return the most relevant results first.
// This can be used when searching for subreddits and users.
func SortByRelevance(opts url.Values) {
	opts.Set("sort", "relevance")
}

// SortByActivity sets the sort option to return results with the most activity first.
// This can be used when searching for subreddits and users.
func SortByActivity(opts url.Values) {
	opts.Set("sort", "activity")
}

// SortByNumberOfComments sets the sort option to return the results with the highest
// number of comments first.
func SortByNumberOfComments(opts url.Values) {
	opts.Set("sort", "comments")
}

// FromThePastHour sets the timespan option to return results from the past hour.
func FromThePastHour(opts url.Values) {
	opts.Set("t", "hour")
}

// FromThePastDay sets the timespan option to return results from the past day.
func FromThePastDay(opts url.Values) {
	opts.Set("t", "day")
}

// FromThePastWeek sets the timespan option to return results from the past week.
func FromThePastWeek(opts url.Values) {
	opts.Set("t", "week")
}

// FromThePastMonth sets the timespan option to return results from the past month.
func FromThePastMonth(opts url.Values) {
	opts.Set("t", "month")
}

// FromThePastYear sets the timespan option to return results from the past year.
func FromThePastYear(opts url.Values) {
	opts.Set("t", "year")
}

// FromAllTime sets the timespan option to return results from all time.
func FromAllTime(opts url.Values) {
	opts.Set("t", "all")
}

// setType sets the type option.
// For mod actions, it's for the type of action (e.g. "banuser", "spamcomment").
func setType(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("type", v)
	}
}

// setQuery sets the q option.
func setQuery(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("q", v)
	}
}

// setRestrict sets the restrict_sr option.
func setRestrict(opts url.Values) {
	opts.Set("restrict_sr", "true")
}
