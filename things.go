package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	kindComment    = "t1"
	kindAccount    = "t2"
	kindLink       = "t3" // a link is a post
	kindMessage    = "t4"
	kindSubreddit  = "t5"
	kindAward      = "t6"
	kindListing    = "Listing"
	kindKarmaList  = "KarmaList"
	kindTrophyList = "TrophyList"
	kindUserList   = "UserList"
	kindMore       = "more"
	kindModAction  = "modaction"
)

// Permalink is the link to a post or comment.
type Permalink string

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *Permalink) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	v = "https://www.reddit.com" + v
	*p = Permalink(v)
	return nil
}

// todo: rename this to thing
type root struct {
	Kind string      `json:"kind,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type rootListing struct {
	Kind string   `json:"kind,omitempty"`
	Data *Listing `json:"data"`
}

// Listing holds things coming from the Reddit API
// It also contains the after/before anchors useful for subsequent requests
type Listing struct {
	Things Things `json:"children"`
	After  string `json:"after"`
	Before string `json:"before"`
}

// Things are objects/entities coming from the Reddit API.
type Things struct {
	Comments     []*Comment
	MoreComments *More

	Users      []*User
	Posts      []*Post
	Subreddits []*Subreddit
	ModActions []*ModAction
	// todo: add the other kinds of things
}

func (t *Things) init() {
	if t.Comments == nil {
		t.Comments = make([]*Comment, 0)
	}
	if t.Users == nil {
		t.Users = make([]*User, 0)
	}
	if t.Posts == nil {
		t.Posts = make([]*Post, 0)
	}
	if t.Subreddits == nil {
		t.Subreddits = make([]*Subreddit, 0)
	}
	if t.ModActions == nil {
		t.ModActions = make([]*ModAction, 0)
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Things) UnmarshalJSON(b []byte) error {
	t.init()

	type thing struct {
		Kind string          `json:"kind"`
		Data json.RawMessage `json:"data"`
	}

	var things []thing
	if err := json.Unmarshal(b, &things); err != nil {
		return err
	}

	for _, thing := range things {
		switch thing.Kind {
		case kindComment:
			v := new(Comment)
			if err := json.Unmarshal(thing.Data, v); err == nil {
				t.Comments = append(t.Comments, v)
			}
		case kindMore:
			v := new(More)
			if err := json.Unmarshal(thing.Data, v); err == nil {
				t.MoreComments = v
			}
		case kindAccount:
			v := new(User)
			if err := json.Unmarshal(thing.Data, v); err == nil {
				t.Users = append(t.Users, v)
			}
		case kindLink:
			v := new(Post)
			if err := json.Unmarshal(thing.Data, v); err == nil {
				t.Posts = append(t.Posts, v)
			}
		case kindMessage:
		case kindSubreddit:
			v := new(Subreddit)
			if err := json.Unmarshal(thing.Data, v); err == nil {
				t.Subreddits = append(t.Subreddits, v)
			}
		case kindAward:
		case kindModAction:
			v := new(ModAction)
			if err := json.Unmarshal(thing.Data, v); err == nil {
				t.ModActions = append(t.ModActions, v)
			}
		}
	}

	return nil
}

// Comment is a comment posted by a user
type Comment struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`
	Edited  *Timestamp `json:"edited,omitempty"`

	ParentID  string    `json:"parent_id,omitempty"`
	Permalink Permalink `json:"permalink,omitempty"`

	Body            string `json:"body,omitempty"`
	Author          string `json:"author,omitempty"`
	AuthorID        string `json:"author_fullname,omitempty"`
	AuthorFlairText string `json:"author_flair_text,omitempty"`
	AuthorFlairID   string `json:"author_flair_template_id,omitempty"`

	Subreddit             string `json:"subreddit,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string `json:"subreddit_id,omitempty"`

	// Indicates if you've upvote/downvoted (true/false).
	// If neither, it will be nil.
	Likes *bool `json:"likes"`

	Score            int `json:"score"`
	Controversiality int `json:"controversiality"`

	// todo: check the validity of these comments
	PostID string `json:"link_id,omitempty"`
	// This doesn't appear when submitting a comment.
	PostTitle string `json:"link_title,omitempty"`
	// This doesn't appear when submitting a comment.
	PostPermalink string `json:"link_permalink,omitempty"`
	// This doesn't appear when submitting a comment.
	PostAuthor string `json:"link_author,omitempty"`
	// This doesn't appear when submitting a comment
	// or when getting a post with its comments.
	PostNumComments *int `json:"num_comments,omitempty"`

	IsSubmitter bool `json:"is_submitter"`
	ScoreHidden bool `json:"score_hidden"`
	Saved       bool `json:"saved"`
	Stickied    bool `json:"stickied"`
	Locked      bool `json:"locked"`
	CanGild     bool `json:"can_gild"`
	NSFW        bool `json:"over_18"`

	Replies Replies `json:"replies"`
}

func (c *Comment) hasMore() bool {
	return c.Replies.MoreComments != nil && len(c.Replies.MoreComments.Children) > 0
}

// Replies holds replies to a comment.
// It contains both comments and "more" comments, which are entrypoints to other
// comments that were left out.
type Replies struct {
	Comments     []*Comment `json:"comments,omitempty"`
	MoreComments *More      `json:"more,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Replies) UnmarshalJSON(data []byte) error {
	// if a comment has no replies, its "replies" field is set to ""
	if string(data) == `""` {
		r = nil
		return nil
	}

	root := new(rootListing)
	err := json.Unmarshal(data, root)
	if err != nil {
		return err
	}

	if root.Data != nil {
		r.Comments = root.Data.Things.Comments
		r.MoreComments = root.Data.Things.MoreComments
	}

	return nil
}

// todo: should we implement json.Marshaler?

// More holds information
type More struct {
	ID       string `json:"id"`
	FullID   string `json:"name"`
	ParentID string `json:"parent_id"`
	// Total number of replies to the parent + replies to those replies (recursively).
	Count int `json:"count"`
	// Number of comment nodes from the parent down to the furthest comment node.
	Depth    int      `json:"depth"`
	Children []string `json:"children"`
}

// Post is a submitted post on Reddit.
type Post struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`
	Edited  *Timestamp `json:"edited,omitempty"`

	Permalink Permalink `json:"permalink,omitempty"`
	URL       string    `json:"url,omitempty"`

	Title string `json:"title,omitempty"`
	Body  string `json:"selftext,omitempty"`

	// Indicates if you've upvote/downvoted (true/false).
	// If neither, it will be nil.
	Likes *bool `json:"likes"`

	Score            int     `json:"score"`
	UpvoteRatio      float32 `json:"upvote_ratio"`
	NumberOfComments int     `json:"num_comments"`

	SubredditID           string `json:"subreddit_id,omitempty"`
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

func (p Post) String() string {
	chunks := []string{
		fmt.Sprintf("[%d]", p.Score),
		p.SubredditNamePrefixed,
		"-",
		p.Title,
		"-",
		string(p.Permalink),
	}
	return strings.Join(chunks, " ")
}

// Subreddit holds information about a subreddit
type Subreddit struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`

	URL                  string `json:"url,omitempty"`
	Name                 string `json:"display_name,omitempty"`
	NamePrefixed         string `json:"display_name_prefixed,omitempty"`
	Title                string `json:"title,omitempty"`
	Description          string `json:"public_description,omitempty"`
	Type                 string `json:"subreddit_type,omitempty"`
	SuggestedCommentSort string `json:"suggested_comment_sort,omitempty"`

	Subscribers     int  `json:"subscribers"`
	ActiveUserCount *int `json:"active_user_count,omitempty"`
	NSFW            bool `json:"over18"`
	UserIsMod       bool `json:"user_is_moderator"`
	Favorite        bool `json:"user_has_favorited"`
}

func (rl *rootListing) getAfter() string {
	if rl == nil || rl.Data == nil {
		return ""
	}
	return rl.Data.After
}

func (rl *rootListing) getBefore() string {
	if rl == nil || rl.Data == nil {
		return ""
	}
	return rl.Data.Before
}

func (rl *rootListing) getComments() *Comments {
	v := new(Comments)
	if rl != nil && rl.Data != nil {
		v.Comments = rl.Data.Things.Comments
		v.After = rl.Data.After
		v.Before = rl.Data.Before
	}
	return v
}

func (rl *rootListing) getMoreComments() *More {
	if rl == nil || rl.Data == nil {
		return nil
	}
	return rl.Data.Things.MoreComments
}

func (rl *rootListing) getUsers() *Users {
	v := new(Users)
	if rl != nil && rl.Data != nil {
		v.Users = rl.Data.Things.Users
		v.After = rl.Data.After
		v.Before = rl.Data.Before
	}
	return v
}

func (rl *rootListing) getPosts() *Posts {
	v := new(Posts)
	if rl != nil && rl.Data != nil {
		v.Posts = rl.Data.Things.Posts
		v.After = rl.Data.After
		v.Before = rl.Data.Before
	}
	return v
}

func (rl *rootListing) getSubreddits() *Subreddits {
	v := new(Subreddits)
	if rl != nil && rl.Data != nil {
		v.Subreddits = rl.Data.Things.Subreddits
		v.After = rl.Data.After
		v.Before = rl.Data.Before
	}
	return v
}

func (rl *rootListing) getModeratorActions() *ModActions {
	v := new(ModActions)
	if rl != nil && rl.Data != nil {
		v.ModActions = rl.Data.Things.ModActions
		v.After = rl.Data.After
		v.Before = rl.Data.Before
	}
	return v
}

// Comments is a list of comments
type Comments struct {
	Comments []*Comment `json:"comments"`
	After    string     `json:"after"`
	Before   string     `json:"before"`
}

// Users is a list of users
type Users struct {
	Users  []*User `json:"users"`
	After  string  `json:"after"`
	Before string  `json:"before"`
}

// Subreddits is a list of subreddits
type Subreddits struct {
	Subreddits []*Subreddit `json:"subreddits"`
	After      string       `json:"after"`
	Before     string       `json:"before"`
}

// Posts is a list of posts.
type Posts struct {
	Posts  []*Post `json:"posts"`
	After  string  `json:"after"`
	Before string  `json:"before"`
}

// ModActions is a list of moderator action.
type ModActions struct {
	ModActions []*ModAction `json:"moderator_actions"`
	After      string       `json:"after"`
	Before     string       `json:"before"`
}

// PostAndComments is a post and its comments.
type PostAndComments struct {
	Post         *Post      `json:"post"`
	Comments     []*Comment `json:"comments"`
	moreComments *More
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// When getting a sticky post, you get an array of 2 Listings
// The 1st one contains the single post in its children array
// The 2nd one contains the comments to the post
func (pc *PostAndComments) UnmarshalJSON(data []byte) error {
	var l [2]rootListing

	err := json.Unmarshal(data, &l)
	if err != nil {
		return err
	}

	post := l[0].getPosts().Posts[0]
	comments := l[1].getComments().Comments
	moreComments := l[1].getMoreComments()

	pc.Post = post
	pc.Comments = comments
	pc.moreComments = moreComments

	return nil
}

func (pc *PostAndComments) hasMore() bool {
	return pc.moreComments != nil && len(pc.moreComments.Children) > 0
}
