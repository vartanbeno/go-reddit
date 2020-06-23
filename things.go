package geddit

import (
	"encoding/json"
	"errors"
)

const (
	kindComment    = "t1"
	kindAccount    = "t2"
	kindLink       = "t3" // a link is a post
	kindMessage    = "t4"
	kindSubreddit  = "t5"
	kindAward      = "t6"
	kindListing    = "Listing"
	kingTrophyList = "TrophyList"
	kindUserList   = "UserList"
	kindMode       = "more"
)

// Sort is a sorting option.
type Sort int

var sorts = [...]string{
	"hot",
	"best",
	"new",
	"rising",
	"controversial",
	"top",
	"relevance",
	"comments",
}

// Different sorting options.
const (
	SortHot Sort = iota
	SortBest
	SortNew
	SortRising
	SortControversial
	SortTop
	SortRelevance
	SortComments
)

func (s Sort) String() string {
	if s < SortHot || s > SortComments {
		return ""
	}
	return sorts[s]
}

// Timespan is a timespan option.
// E.g. "hour" means in the last hour, "all" means all-time.
// It is used when conducting searches.
type Timespan int

var timespans = [...]string{
	"hour",
	"day",
	"week",
	"month",
	"year",
	"all",
}

// Different timespan options.
const (
	TimespanHour Timespan = iota
	TimespanDay
	TimespanWeek
	TimespanMonth
	TimespanYear
	TimespanAll
)

func (t Timespan) String() string {
	if t < TimespanHour || t > TimespanAll {
		return ""
	}
	return timespans[t]
}

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
	Comments   []Comment   `json:"comments,omitempty"`
	Users      []User      `json:"users,omitempty"`
	Posts      []Post      `json:"posts,omitempty"`
	Subreddits []Subreddit `json:"subreddits,omitempty"`
	// todo: add the other kinds of things
}

type commentRoot struct {
	Kind string   `json:"kind,omitempty"`
	Data *Comment `json:"data,omitempty"`
}

type userRoot struct {
	Kind string `json:"kind,omitempty"`
	Data *User  `json:"data,omitempty"`
}

type postRoot struct {
	Kind string `json:"kind,omitempty"`
	Data *Post  `json:"data,omitempty"`
}

type subredditRoot struct {
	Kind string     `json:"kind,omitempty"`
	Data *Subreddit `json:"data,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Things) UnmarshalJSON(b []byte) error {
	var children []map[string]interface{}
	if err := json.Unmarshal(b, &children); err != nil {
		return err
	}

	for _, child := range children {
		byteValue, _ := json.Marshal(child)
		switch child["kind"] {
		// todo: kindMore
		case kindComment:
			root := new(commentRoot)
			if err := json.Unmarshal(byteValue, root); err == nil && root.Data != nil {
				t.Comments = append(t.Comments, *root.Data)
			}
		case kindAccount:
			root := new(userRoot)
			if err := json.Unmarshal(byteValue, root); err == nil && root.Data != nil {
				t.Users = append(t.Users, *root.Data)
			}
		case kindLink:
			root := new(postRoot)
			if err := json.Unmarshal(byteValue, root); err == nil && root.Data != nil {
				t.Posts = append(t.Posts, *root.Data)
			}
		case kindMessage:
		case kindSubreddit:
			root := new(subredditRoot)
			if err := json.Unmarshal(byteValue, root); err == nil && root.Data != nil {
				t.Subreddits = append(t.Subreddits, *root.Data)
			}
		case kindAward:
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

	ParentID  string `json:"parent_id,omitempty"`
	Permalink string `json:"permalink,omitempty"`

	Body            string `json:"body,omitempty"`
	Author          string `json:"author,omitempty"`
	AuthorID        string `json:"author_fullname,omitempty"`
	AuthorFlairText string `json:"author_flair_text,omitempty"`
	AuthorFlairID   string `json:"author_flair_template_id,omitempty"`

	Subreddit             string `json:"subreddit,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string `json:"subreddit_id,omitempty"`

	// Indicates if the client has upvote/downvoted (true/false)
	// If neither, it will be nil
	Likes *bool `json:"likes"`

	Score            int `json:"score"`
	Controversiality int `json:"controversiality"`

	PostID string `json:"link_id,omitempty"`

	// These don't appear when submitting a comment
	PostTitle       string `json:"link_title,omitempty"`
	PostPermalink   string `json:"link_permalink,omitempty"`
	PostAuthor      string `json:"link_author,omitempty"`
	PostNumComments int    `json:"num_comments"`

	IsSubmitter bool `json:"is_submitter"`
	ScoreHidden bool `json:"score_hidden"`
	Saved       bool `json:"saved"`
	Stickied    bool `json:"stickied"`
	Locked      bool `json:"locked"`
	CanGild     bool `json:"can_gild"`
	NSFW        bool `json:"over_18"`

	Replies Replies `json:"replies"`
}

// Replies are replies to a comment
type Replies []Comment

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Replies) UnmarshalJSON(data []byte) error {
	// if a comment has no replies, its "replies" field is set to ""
	if string(data) == `""` {
		return nil
	}

	root := new(rootListing)
	err := json.Unmarshal(data, root)
	if err != nil {
		return err
	}

	*r = root.getComments().Comments
	return nil
}

// Post is a submitted post on Reddit.
type Post struct {
	ID      string     `json:"id,omitempty"`
	FullID  string     `json:"name,omitempty"`
	Created *Timestamp `json:"created_utc,omitempty"`
	Edited  *Timestamp `json:"edited,omitempty"`

	Permalink string `json:"permalink,omitempty"`
	URL       string `json:"url,omitempty"`

	Title string `json:"title,omitempty"`
	Body  string `json:"selftext,omitempty"`

	// Indicates if the client has upvote/downvoted (true/false)
	// If neither, it will be nil
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

// Comments is a list of comments
type Comments struct {
	Comments []Comment `json:"comments,omitempty"`
	After    string    `json:"after"`
	Before   string    `json:"before"`
}

// Users is a list of users
type Users struct {
	Users  []User `json:"users,omitempty"`
	After  string `json:"after"`
	Before string `json:"before"`
}

// Subreddits is a list of subreddits
type Subreddits struct {
	Subreddits []Subreddit `json:"subreddits,omitempty"`
	After      string      `json:"after"`
	Before     string      `json:"before"`
}

// Posts is a list of posts.
type Posts struct {
	Posts  []Post `json:"submissions,omitempty"`
	After  string `json:"after"`
	Before string `json:"before"`
}

// PostAndComments is a post and its comments
type PostAndComments struct {
	Post     Post      `json:"post,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// When getting a sticky post, you get an array of 2 Listings
// The 1st one contains the single post in its children array
// The 2nd one contains the comments to the post
func (pc *PostAndComments) UnmarshalJSON(data []byte) error {
	var l []rootListing

	err := json.Unmarshal(data, &l)
	if err != nil {
		return err
	}

	if len(l) < 2 {
		return errors.New("unexpected json response when getting post")
	}

	post := l[0].getPosts().Posts[0]
	comments := l[1].getComments().Comments

	pc.Post = post
	pc.Comments = comments

	return nil
}
