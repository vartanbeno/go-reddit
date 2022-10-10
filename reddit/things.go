package reddit

import (
	"encoding/json"
	"fmt"
)

const (
	kindComment           = "t1"
	kindUser              = "t2"
	kindPost              = "t3"
	kindMessage           = "t4"
	kindSubreddit         = "t5"
	kindTrophy            = "t6"
	kindListing           = "Listing"
	kindSubredditSettings = "subreddit_settings"
	kindKarmaList         = "KarmaList"
	kindTrophyList        = "TrophyList"
	kindUserList          = "UserList"
	kindMore              = "more"
	kindLiveThread        = "LiveUpdateEvent"
	kindLiveThreadUpdate  = "LiveUpdate"
	kindModAction         = "modaction"
	kindMulti             = "LabeledMulti"
	kindMultiDescription  = "LabeledMultiDescription"
	kindWikiPage          = "wikipage"
	kindWikiPageListing   = "wikipagelisting"
	kindWikiPageSettings  = "wikipagesettings"
	kindStyleSheet        = "stylesheet"
)

type anchor interface {
	After() string
}

// thing is an entity on Reddit.
// Its kind reprsents what it is and what is stored in the Data field.
// e.g. t1 = comment, t2 = user, t3 = post, etc.
type thing struct {
	Kind string      `json:"kind"`
	Data interface{} `json:"data"`
}

func (t *thing) After() string {
	if t == nil {
		return ""
	}
	a, ok := t.Data.(anchor)
	if !ok {
		return ""
	}
	return a.After()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *thing) UnmarshalJSON(b []byte) error {
	root := new(struct {
		Kind string          `json:"kind"`
		Data json.RawMessage `json:"data"`
	})

	err := json.Unmarshal(b, root)
	if err != nil {
		return err
	}

	t.Kind = root.Kind
	var v interface{}

	switch t.Kind {
	case kindListing:
		v = new(listing)
	case kindComment:
		v = new(Comment)
	case kindMore:
		v = new(More)
	case kindUser:
		v = new(User)
	case kindPost:
		v = new(Post)
	case kindSubreddit:
		v = new(Subreddit)
	case kindSubredditSettings:
		v = new(SubredditSettings)
	case kindLiveThread:
		v = new(LiveThread)
	case kindLiveThreadUpdate:
		v = new(LiveThreadUpdate)
	case kindModAction:
		v = new(ModAction)
	case kindMulti:
		v = new(Multi)
	case kindMultiDescription:
		v = new(rootMultiDescription)
	case kindTrophy:
		v = new(Trophy)
	case kindTrophyList:
		v = new(trophyList)
	case kindKarmaList:
		v = new([]*SubredditKarma)
	case kindWikiPage:
		v = new(WikiPage)
	case kindWikiPageListing:
		v = new([]string)
	case kindWikiPageSettings:
		v = new(WikiPageSettings)
	case kindStyleSheet:
		v = new(SubredditStyleSheet)
	default:
		return fmt.Errorf("unrecognized kind: %q", t.Kind)
	}

	err = json.Unmarshal(root.Data, v)
	if err != nil {
		return err
	}

	t.Data = v
	return nil
}

func (t *thing) Listing() (v *listing, ok bool) {
	v, ok = t.Data.(*listing)
	return
}

func (t *thing) Comment() (v *Comment, ok bool) {
	v, ok = t.Data.(*Comment)
	return
}

func (t *thing) More() (v *More, ok bool) {
	v, ok = t.Data.(*More)
	return
}

func (t *thing) User() (v *User, ok bool) {
	v, ok = t.Data.(*User)
	return
}

func (t *thing) Post() (v *Post, ok bool) {
	v, ok = t.Data.(*Post)
	return
}

func (t *thing) Subreddit() (v *Subreddit, ok bool) {
	v, ok = t.Data.(*Subreddit)
	return
}

func (t *thing) SubredditSettings() (v *SubredditSettings, ok bool) {
	v, ok = t.Data.(*SubredditSettings)
	return
}

func (t *thing) LiveThread() (v *LiveThread, ok bool) {
	v, ok = t.Data.(*LiveThread)
	return
}

func (t *thing) LiveThreadUpdate() (v *LiveThreadUpdate, ok bool) {
	v, ok = t.Data.(*LiveThreadUpdate)
	return
}

func (t *thing) ModAction() (v *ModAction, ok bool) {
	v, ok = t.Data.(*ModAction)
	return
}

func (t *thing) Multi() (v *Multi, ok bool) {
	v, ok = t.Data.(*Multi)
	return
}

func (t *thing) MultiDescription() (s string, ok bool) {
	v, ok := t.Data.(*rootMultiDescription)
	if ok {
		s = v.Body
	}
	return
}

func (t *thing) Trophy() (v *Trophy, ok bool) {
	v, ok = t.Data.(*Trophy)
	return
}

func (t *thing) TrophyList() ([]*Trophy, bool) {
	v, ok := t.Data.(*trophyList)
	if !ok {
		return nil, ok
	}
	return *v, ok
}

func (t *thing) Karma() ([]*SubredditKarma, bool) {
	v, ok := t.Data.(*[]*SubredditKarma)
	if !ok {
		return nil, ok
	}
	return *v, ok
}

func (t *thing) WikiPage() (v *WikiPage, ok bool) {
	v, ok = t.Data.(*WikiPage)
	return
}

func (t *thing) WikiPages() ([]string, bool) {
	v, ok := t.Data.(*[]string)
	if !ok {
		return nil, ok
	}
	return *v, ok
}

func (t *thing) WikiPageSettings() (v *WikiPageSettings, ok bool) {
	v, ok = t.Data.(*WikiPageSettings)
	return
}

func (t *thing) StyleSheet() (v *SubredditStyleSheet, ok bool) {
	v, ok = t.Data.(*SubredditStyleSheet)
	return
}

// listing is a list of things coming from the Reddit API.
// It also contains the after anchor useful to get the next results via subsequent requests.
type listing struct {
	things things
	after  string
}

func (l *listing) After() string {
	return l.after
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *listing) UnmarshalJSON(b []byte) error {
	root := new(struct {
		Things things `json:"children"`
		After  string `json:"after"`
	})

	err := json.Unmarshal(b, root)
	if err != nil {
		return err
	}

	l.things = root.Things
	l.after = root.After

	return nil
}

func (l *listing) Comments() []*Comment {
	if l == nil {
		return nil
	}
	return l.things.Comments
}

func (l *listing) Mores() []*More {
	if l == nil {
		return nil
	}
	return l.things.Mores
}

func (l *listing) Users() []*User {
	if l == nil {
		return nil
	}
	return l.things.Users
}

func (l *listing) Posts() []*Post {
	if l == nil {
		return nil
	}
	return l.things.Posts
}

func (l *listing) Subreddits() []*Subreddit {
	if l == nil {
		return nil
	}
	return l.things.Subreddits
}

func (l *listing) ModActions() []*ModAction {
	if l == nil {
		return nil
	}
	return l.things.ModActions
}

func (l *listing) Multis() []*Multi {
	if l == nil {
		return nil
	}
	return l.things.Multis
}

func (l *listing) LiveThreads() []*LiveThread {
	if l == nil {
		return nil
	}
	return l.things.LiveThreads
}

func (l *listing) LiveThreadUpdates() []*LiveThreadUpdate {
	if l == nil {
		return nil
	}
	return l.things.LiveThreadUpdates
}

type things struct {
	Comments          []*Comment
	Mores             []*More
	Users             []*User
	Posts             []*Post
	Subreddits        []*Subreddit
	ModActions        []*ModAction
	Multis            []*Multi
	LiveThreads       []*LiveThread
	LiveThreadUpdates []*LiveThreadUpdate
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *things) UnmarshalJSON(b []byte) error {
	var things []thing
	if err := json.Unmarshal(b, &things); err != nil {
		return err
	}

	t.add(things...)
	return nil
}

func (t *things) add(things ...thing) {
	for _, thing := range things {
		switch v := thing.Data.(type) {
		case *Comment:
			t.Comments = append(t.Comments, v)
		case *More:
			t.Mores = append(t.Mores, v)
		case *User:
			t.Users = append(t.Users, v)
		case *Post:
			t.Posts = append(t.Posts, v)
		case *Subreddit:
			t.Subreddits = append(t.Subreddits, v)
		case *ModAction:
			t.ModActions = append(t.ModActions, v)
		case *Multi:
			t.Multis = append(t.Multis, v)
		case *LiveThread:
			t.LiveThreads = append(t.LiveThreads, v)
		case *LiveThreadUpdate:
			t.LiveThreadUpdates = append(t.LiveThreadUpdates, v)
		}
	}
}

type trophyList []*Trophy

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *trophyList) UnmarshalJSON(b []byte) error {
	root := new(struct {
		Trophies []thing `json:"trophies"`
	})

	err := json.Unmarshal(b, root)
	if err != nil {
		return err
	}

	*l = make(trophyList, 0, len(root.Trophies))
	for _, thing := range root.Trophies {
		if trophy, ok := thing.Trophy(); ok {
			*l = append(*l, trophy)
		}
	}

	return nil
}

// Comment is a comment posted by a user.
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

	SubredditName         string `json:"subreddit,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string `json:"subreddit_id,omitempty"`

	// Indicates if you've upvote/downvoted (true/false).
	// If neither, it will be nil.
	Likes *bool `json:"likes"`

	Score            int `json:"score"`
	Controversiality int `json:"controversiality"`

	PostID string `json:"link_id,omitempty"`
	// This doesn't appear consistently.
	PostTitle string `json:"link_title,omitempty"`
	// This doesn't appear consistently.
	PostPermalink string `json:"link_permalink,omitempty"`
	// This doesn't appear consistently.
	PostAuthor string `json:"link_author,omitempty"`
	// This doesn't appear consistently.
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

// HasMore determines whether the comment has more replies to load in its reply tree.
func (c *Comment) HasMore() bool {
	return c.Replies.More != nil && len(c.Replies.More.Children) > 0
}

// addCommentToReplies traverses the comment tree to find the one
// that the 2nd comment is replying to. It then adds it to its replies.
func (c *Comment) addCommentToReplies(comment *Comment) {
	if c.FullID == comment.ParentID {
		c.Replies.Comments = append(c.Replies.Comments, comment)
		return
	}

	for _, reply := range c.Replies.Comments {
		reply.addCommentToReplies(comment)
	}
}

func (c *Comment) addMoreToReplies(more *More) {
	if c.FullID == more.ParentID {
		c.Replies.More = more
		return
	}

	for _, reply := range c.Replies.Comments {
		reply.addMoreToReplies(more)
	}
}

// Replies holds replies to a comment.
// It contains both comments and "more" comments, which are entrypoints to other
// comments that were left out.
type Replies struct {
	Comments []*Comment `json:"comments,omitempty"`
	More     *More      `json:"-"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Replies) UnmarshalJSON(data []byte) error {
	// if a comment has no replies, its "replies" field is set to ""
	if string(data) == `""` {
		r = nil
		return nil
	}

	root := new(thing)
	err := json.Unmarshal(data, root)
	if err != nil {
		return err
	}

	listing, _ := root.Listing()

	r.Comments = listing.Comments()
	if len(listing.Mores()) > 0 {
		r.More = listing.Mores()[0]
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (r *Replies) MarshalJSON() ([]byte, error) {
	if r == nil || len(r.Comments) == 0 {
		return []byte(`null`), nil
	}
	return json.Marshal(r.Comments)
}

// More holds information used to retrieve additional comments omitted from a base comment tree.
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

	Permalink string `json:"permalink,omitempty"`
	URL       string `json:"url,omitempty"`

	Title string `json:"title,omitempty"`
	Body  string `json:"selftext,omitempty"`

	// Indicates if you've upvoted/downvoted (true/false).
	// If neither, it will be nil.
	Likes *bool `json:"likes"`

	Score            int     `json:"score"`
	UpvoteRatio      float32 `json:"upvote_ratio"`
	NumberOfComments int     `json:"num_comments"`

	SubredditName         string `json:"subreddit,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`
	SubredditID           string `json:"subreddit_id,omitempty"`
	SubredditSubscribers  int    `json:"subreddit_subscribers"`

	Author   string `json:"author,omitempty"`
	AuthorID string `json:"author_fullname,omitempty"`

	Spoiler    bool `json:"spoiler"`
	Locked     bool `json:"locked"`
	NSFW       bool `json:"over_18"`
	IsSelfPost bool `json:"is_self"`
	Saved      bool `json:"saved"`
	Stickied   bool `json:"stickied"`
}

/*
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
	Subscribed      bool `json:"user_is_subscriber"`
	Favorite        bool `json:"user_has_favorited"`
}
*/
// Subreddit holds information about a subreddit
// This struct was generated from the response from the api.
// Some fields are interface{} because they had empty json values, eg:
// [], {}, null. Where the types were easily inferred from the name, the type
// was changed to a real go type, where it was not interface{} was left. In
// the future if someone finds a real response these can be changed.
type Subreddit struct {
	AcceptFollowers             bool        `json:"accept_followers"`
	AccountsActive              int         `json:"accounts_active"`
	AccountsActiveIsFuzzed      bool        `json:"accounts_active_is_fuzzed"`
	ActiveUserCount             int         `json:"active_user_count"`
	AdvertiserCategory          string      `json:"advertiser_category"`
	AllOriginalContent          bool        `json:"all_original_content"`
	AllowChatPostCreation       bool        `json:"allow_chat_post_creation"`
	AllowDiscovery              bool        `json:"allow_discovery"`
	AllowGalleries              bool        `json:"allow_galleries"`
	AllowImages                 bool        `json:"allow_images"`
	AllowPolls                  bool        `json:"allow_polls"`
	AllowPredictionContributors bool        `json:"allow_prediction_contributors"`
	AllowPredictions            bool        `json:"allow_predictions"`
	AllowPredictionsTournament  bool        `json:"allow_predictions_tournament"`
	AllowTalks                  bool        `json:"allow_talks"`
	AllowVideogifs              bool        `json:"allow_videogifs"`
	AllowVideos                 bool        `json:"allow_videos"`
	AllowedMediaInComments      interface{} `json:"allowed_media_in_comments"`
	BannerBackgroundColor       string      `json:"banner_background_color"`
	BannerBackgroundImage       string      `json:"banner_background_image"`
	BannerImg                   string      `json:"banner_img"`
	BannerSize                  []float64   `json:"banner_size"`
	CanAssignLinkFlair          bool        `json:"can_assign_link_flair"`
	CanAssignUserFlair          bool        `json:"can_assign_user_flair"`
	CollapseDeletedComments     bool        `json:"collapse_deleted_comments"`
	CommentContributionSettings struct {
	} `json:"comment_contribution_settings"`
	CommentScoreHideMins             int           `json:"comment_score_hide_mins"`
	CommunityIcon                    string        `json:"community_icon"`
	CommunityReviewed                bool          `json:"community_reviewed"`
	Created                          *Timestamp    `json:"created"`
	CreatedUtc                       *Timestamp    `json:"created_utc"`
	Description                      string        `json:"description"`
	DescriptionHTML                  string        `json:"description_html"`
	DisableContributorRequests       bool          `json:"disable_contributor_requests"`
	DisplayName                      string        `json:"display_name"`
	DisplayNamePrefixed              string        `json:"display_name_prefixed"`
	EmojisCustomSize                 interface{}   `json:"emojis_custom_size"`
	EmojisEnabled                    bool          `json:"emojis_enabled"`
	FreeFormReports                  bool          `json:"free_form_reports"`
	HasMenuWidget                    bool          `json:"has_menu_widget"`
	HeaderImg                        string        `json:"header_img"`
	HeaderSize                       []int         `json:"header_size"`
	HeaderTitle                      string        `json:"header_title"`
	HideAds                          bool          `json:"hide_ads"`
	IconImg                          string        `json:"icon_img"`
	IconSize                         []float64     `json:"icon_size"`
	ID                               string        `json:"id"`
	IsChatPostFeatureEnabled         bool          `json:"is_chat_post_feature_enabled"`
	IsCrosspostableSubreddit         bool          `json:"is_crosspostable_subreddit"`
	IsEnrolledInNewModmail           bool          `json:"is_enrolled_in_new_modmail"`
	KeyColor                         string        `json:"key_color"`
	Lang                             string        `json:"lang"`
	LinkFlairEnabled                 bool          `json:"link_flair_enabled"`
	LinkFlairPosition                string        `json:"link_flair_position"`
	MobileBannerImage                string        `json:"mobile_banner_image"`
	Name                             string        `json:"name"`
	NotificationLevel                interface{}   `json:"notification_level"`
	OriginalContentTagEnabled        bool          `json:"original_content_tag_enabled"`
	Over18                           bool          `json:"over18"`
	PredictionLeaderboardEntryType   string        `json:"prediction_leaderboard_entry_type"`
	PrimaryColor                     string        `json:"primary_color"`
	PublicDescription                string        `json:"public_description"`
	PublicDescriptionHTML            string        `json:"public_description_html"`
	PublicTraffic                    bool          `json:"public_traffic"`
	Quarantine                       bool          `json:"quarantine"`
	RestrictCommenting               bool          `json:"restrict_commenting"`
	RestrictPosting                  bool          `json:"restrict_posting"`
	ShouldArchivePosts               bool          `json:"should_archive_posts"`
	ShouldShowMediaInCommentsSetting bool          `json:"should_show_media_in_comments_setting"`
	ShowMedia                        bool          `json:"show_media"`
	ShowMediaPreview                 bool          `json:"show_media_preview"`
	SpoilersEnabled                  bool          `json:"spoilers_enabled"`
	SubmissionType                   string        `json:"submission_type"`
	SubmitLinkLabel                  string        `json:"submit_link_label"`
	SubmitText                       string        `json:"submit_text"`
	SubmitTextHTML                   string        `json:"submit_text_html"`
	SubmitTextLabel                  string        `json:"submit_text_label"`
	SubredditType                    string        `json:"subreddit_type"`
	Subscribers                      int           `json:"subscribers"`
	SuggestedCommentSort             string        `json:"suggested_comment_sort"`
	Title                            string        `json:"title"`
	URL                              string        `json:"url"`
	UserCanFlairInSr                 interface{}   `json:"user_can_flair_in_sr"`
	UserFlairBackgroundColor         interface{}   `json:"user_flair_background_color"`
	UserFlairCSSClass                interface{}   `json:"user_flair_css_class"`
	UserFlairEnabledInSr             bool          `json:"user_flair_enabled_in_sr"`
	UserFlairPosition                string        `json:"user_flair_position"`
	UserFlairRichtext                []interface{} `json:"user_flair_richtext"`
	UserFlairTemplateID              interface{}   `json:"user_flair_template_id"`
	UserFlairText                    interface{}   `json:"user_flair_text"`
	UserFlairTextColor               interface{}   `json:"user_flair_text_color"`
	UserFlairType                    string        `json:"user_flair_type"`
	UserHasFavorited                 bool          `json:"user_has_favorited"`
	UserIsBanned                     bool          `json:"user_is_banned"`
	UserIsContributor                bool          `json:"user_is_contributor"`
	UserIsModerator                  bool          `json:"user_is_moderator"`
	UserIsMuted                      bool          `json:"user_is_muted"`
	UserIsSubscriber                 bool          `json:"user_is_subscriber"`
	UserSrFlairEnabled               bool          `json:"user_sr_flair_enabled"`
	UserSrThemeEnabled               bool          `json:"user_sr_theme_enabled"`
	WhitelistStatus                  string        `json:"whitelist_status"`
	WikiEnabled                      bool          `json:"wiki_enabled"`
	Wls                              int           `json:"wls"`
}

// PostAndComments is a post and its comments.
type PostAndComments struct {
	Post     *Post      `json:"post"`
	Comments []*Comment `json:"comments"`
	More     *More      `json:"-"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// When getting a sticky post, you get an array of 2 Listings
// The 1st one contains the single post in its children array
// The 2nd one contains the comments to the post
func (pc *PostAndComments) UnmarshalJSON(data []byte) error {
	var root [2]thing

	err := json.Unmarshal(data, &root)
	if err != nil {
		return err
	}

	listing1, _ := root[0].Listing()
	listing2, _ := root[1].Listing()

	pc.Post = listing1.Posts()[0]
	pc.Comments = listing2.Comments()
	if len(listing2.Mores()) > 0 {
		pc.More = listing2.Mores()[0]
	}

	return nil
}

// HasMore determines whether the post has more replies to load in its reply tree.
func (pc *PostAndComments) HasMore() bool {
	return pc.More != nil && len(pc.More.Children) > 0
}

func (pc *PostAndComments) addCommentToTree(comment *Comment) {
	if pc.Post.FullID == comment.ParentID {
		pc.Comments = append(pc.Comments, comment)
		return
	}

	for _, reply := range pc.Comments {
		reply.addCommentToReplies(comment)
	}
}

func (pc *PostAndComments) addMoreToTree(more *More) {
	if pc.Post.FullID == more.ParentID {
		pc.More = more
	}

	for _, reply := range pc.Comments {
		reply.addMoreToReplies(more)
	}
}
