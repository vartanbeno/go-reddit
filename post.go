package reddit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// PostService handles communication with the post
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_links_and_comments
type PostService service

type submittedLinkRoot struct {
	JSON struct {
		Data *Submitted `json:"data,omitempty"`
	} `json:"json"`
}

// Submitted is a newly submitted post on Reddit.
type Submitted struct {
	ID     string `json:"id,omitempty"`
	FullID string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
}

// SubmitTextOptions are options used for text posts.
type SubmitTextOptions struct {
	Subreddit string `url:"sr,omitempty"`
	Title     string `url:"title,omitempty"`
	Text      string `url:"text,omitempty"`

	FlairID   string `url:"flair_id,omitempty"`
	FlairText string `url:"flair_text,omitempty"`

	SendReplies *bool `url:"sendreplies,omitempty"`
	NSFW        bool  `url:"nsfw,omitempty"`
	Spoiler     bool  `url:"spoiler,omitempty"`
}

// SubmitLinkOptions are options used for link posts.
type SubmitLinkOptions struct {
	Subreddit string `url:"sr,omitempty"`
	Title     string `url:"title,omitempty"`
	URL       string `url:"url,omitempty"`

	FlairID   string `url:"flair_id,omitempty"`
	FlairText string `url:"flair_text,omitempty"`

	SendReplies *bool `url:"sendreplies,omitempty"`
	Resubmit    bool  `url:"resubmit,omitempty"`
	NSFW        bool  `url:"nsfw,omitempty"`
	Spoiler     bool  `url:"spoiler,omitempty"`
}

// Get returns a post with its comments.
// id is the ID36 of the post, not its full id.
// Example: instead of t3_abc123, use abc123.
func (s *PostService) Get(ctx context.Context, id string) (*Post, []*Comment, *Response, error) {
	path := fmt.Sprintf("comments/%s", id)
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

func (s *PostService) submit(ctx context.Context, v interface{}) (*Submitted, *Response, error) {
	path := "api/submit"

	form, err := query.Values(v)
	if err != nil {
		return nil, nil, err
	}
	form.Set("api_type", "json")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(submittedLinkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.JSON.Data, resp, nil
}

// SubmitText submits a text post.
func (s *PostService) SubmitText(ctx context.Context, opts SubmitTextOptions) (*Submitted, *Response, error) {
	type submit struct {
		SubmitTextOptions
		Kind string `url:"kind,omitempty"`
	}
	return s.submit(ctx, &submit{opts, "self"})
}

// SubmitLink submits a link post.
func (s *PostService) SubmitLink(ctx context.Context, opts SubmitLinkOptions) (*Submitted, *Response, error) {
	type submit struct {
		SubmitLinkOptions
		Kind string `url:"kind,omitempty"`
	}
	return s.submit(ctx, &submit{opts, "link"})
}

// Edit edits a post.
func (s *PostService) Edit(ctx context.Context, id string, text string) (*Post, *Response, error) {
	path := "api/editusertext"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("return_rtjson", "true")
	form.Set("thing_id", id)
	form.Set("text", text)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(Post)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Hide hides posts.
func (s *PostService) Hide(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/hide"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unhide unhides posts.
func (s *PostService) Unhide(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/unhide"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// MarkNSFW marks a post as NSFW.
func (s *PostService) MarkNSFW(ctx context.Context, id string) (*Response, error) {
	path := "api/marknsfw"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UnmarkNSFW unmarks a post as NSFW.
func (s *PostService) UnmarkNSFW(ctx context.Context, id string) (*Response, error) {
	path := "api/unmarknsfw"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Spoiler marks a post as a spoiler.
func (s *PostService) Spoiler(ctx context.Context, id string) (*Response, error) {
	path := "api/spoiler"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unspoiler unmarks a post as a spoiler.
func (s *PostService) Unspoiler(ctx context.Context, id string) (*Response, error) {
	path := "api/unspoiler"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Sticky stickies a post in its subreddit.
// When bottom is true, the post will be set as the bottom sticky (the 2nd one).
// If no top sticky exists, the post will become the top sticky regardless.
// When attempting to sticky a post that's already stickied, it will return a 409 Conflict error.
func (s *PostService) Sticky(ctx context.Context, id string, bottom bool) (*Response, error) {
	path := "api/set_subreddit_sticky"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("id", id)
	form.Set("state", "true")
	if !bottom {
		form.Set("num", "1")
	}

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unsticky unstickies a post in its subreddit.
func (s *PostService) Unsticky(ctx context.Context, id string) (*Response, error) {
	path := "api/set_subreddit_sticky"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("id", id)
	form.Set("state", "false")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// PinToProfile pins one of your posts to your profile.
// TODO: very inconsistent behaviour, not sure I'm ready to include this parameter yet.
// The pos parameter should be a number between 1-4 (inclusive), indicating the position at which
// the post should appear on your profile.
// Note: The position will be bumped upward if there's space. E.g. if you only have 1 pinned post,
// and you try to pin another post to position 3, it will be pinned at 2.
// When attempting to pin a post that's already pinned, it will return a 409 Conflict error.
func (s *PostService) PinToProfile(ctx context.Context, id string) (*Response, error) {
	path := "api/set_subreddit_sticky"

	// if pos < 1 {
	// 	pos = 1
	// }
	// if pos > 4 {
	// 	pos = 4
	// }

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("id", id)
	form.Set("state", "true")
	form.Set("to_profile", "true")
	// form.Set("num", fmt.Sprint(pos))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UnpinFromProfile unpins one of your posts from your profile.
func (s *PostService) UnpinFromProfile(ctx context.Context, id string) (*Response, error) {
	path := "api/set_subreddit_sticky"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("id", id)
	form.Set("state", "false")
	form.Set("to_profile", "true")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// setSuggestedSort sets the suggested comment sort for the post.
// sort must be one of: confidence (i.e. best), top, new, controversial, old, random, qa, live
func (s *PostService) setSuggestedSort(ctx context.Context, id string, sort string) (*Response, error) {
	path := "api/set_suggested_sort"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("id", id)
	form.Set("sort", sort)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// SetSuggestedSortBest sets the suggested comment sort for the post to best.
func (s *PostService) SetSuggestedSortBest(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "confidence")
}

// SetSuggestedSortTop sets the suggested comment sort for the post to top.
func (s *PostService) SetSuggestedSortTop(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "top")
}

// SetSuggestedSortNew sets the suggested comment sort for the post to new.
func (s *PostService) SetSuggestedSortNew(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "new")
}

// SetSuggestedSortControversial sets the suggested comment sort for the post to controversial.
func (s *PostService) SetSuggestedSortControversial(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "controversial")
}

// SetSuggestedSortOld sorts the comments on the posts randomly.
func (s *PostService) SetSuggestedSortOld(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "old")
}

// SetSuggestedSortRandom sets the suggested comment sort for the post to random.
func (s *PostService) SetSuggestedSortRandom(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "random")
}

// SetSuggestedSortAMA sets the suggested comment sort for the post to a Q&A styled fashion.
func (s *PostService) SetSuggestedSortAMA(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "qa")
}

// SetSuggestedSortLive sets the suggested comment sort for the post to stream new comments as they're posted.
// As of now, this is still in beta, so it's not a fully developed feature yet. It just sets the sort as "new" for now.
func (s *PostService) SetSuggestedSortLive(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "live")
}

// ClearSuggestedSort clears the suggested comment sort for the post.
func (s *PostService) ClearSuggestedSort(ctx context.Context, id string) (*Response, error) {
	return s.setSuggestedSort(ctx, id, "")
}

// EnableContestMode enables contest mode for the post.
// Comments will be sorted randomly and regular users cannot see comment scores.
func (s *PostService) EnableContestMode(ctx context.Context, id string) (*Response, error) {
	path := "api/set_contest_mode"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("id", id)
	form.Set("state", "true")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DisableContestMode disables contest mode for the post.
func (s *PostService) DisableContestMode(ctx context.Context, id string) (*Response, error) {
	path := "api/set_contest_mode"

	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("id", id)
	form.Set("state", "false")

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// More retrieves more comments that were left out when initially fetching the post.
// id is the post's full ID.
// commentIDs are the ID36s of comments.
func (s *PostService) More(ctx context.Context, comment *Comment) (*Response, error) {
	if comment == nil {
		return nil, errors.New("comment: cannot be nil")
	}

	if comment.Replies.MoreComments == nil {
		return nil, nil
	}

	postID := comment.PostID
	commentIDs := comment.Replies.MoreComments.Children

	if len(commentIDs) == 0 {
		return nil, nil
	}

	type query struct {
		PostID  string   `url:"link_id"`
		IDs     []string `url:"children,comma"`
		APIType string   `url:"api_type"`
	}

	path := "api/morechildren"
	path, err := addOptions(path, query{postID, commentIDs, "json"})
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	type rootResponse struct {
		JSON struct {
			Data struct {
				Things Things `json:"things"`
			} `json:"data"`
		} `json:"json"`
	}

	root := new(rootResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	comments := root.JSON.Data.Things.Comments
	for _, c := range comments {
		addCommentToReplies(comment, c)
	}

	comment.Replies.MoreComments = nil
	return resp, nil
}

// addCommentToReplies traverses the comment tree to find the one
// that the 2nd comment is replying to. It then adds it to its replies.
func addCommentToReplies(parent *Comment, comment *Comment) {
	if parent.FullID == comment.ParentID {
		parent.Replies.Comments = append(parent.Replies.Comments, comment)
		return
	}

	for _, reply := range parent.Replies.Comments {
		addCommentToReplies(reply, comment)
	}
}

// RandomFromSubreddits returns a random post and its comments from the subreddits.
// If no subreddits are provided, it will run the query against your subscriptions.
func (s *PostService) RandomFromSubreddits(ctx context.Context, subreddits ...string) (*Post, []*Comment, *Response, error) {
	path := "random"
	if len(subreddits) > 0 {
		path = fmt.Sprintf("r/%s/random", strings.Join(subreddits, "+"))
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

// Random returns a random post and its comments from all of Reddit.
func (s *PostService) Random(ctx context.Context) (*Post, []*Comment, *Response, error) {
	return s.RandomFromSubreddits(ctx, "all")
}

// RandomFromSubscriptions returns a random post and its comments from your subscriptions.
func (s *PostService) RandomFromSubscriptions(ctx context.Context) (*Post, []*Comment, *Response, error) {
	return s.RandomFromSubreddits(ctx)
}
