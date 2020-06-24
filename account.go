package geddit

import (
	"context"
	"fmt"
	"net/http"
)

// AccountService handles communication with the account
// related methods of the Reddit API.
type AccountService interface {
	Karma(ctx context.Context) ([]SubredditKarma, *Response, error)
}

// AccountServiceOp implements the AccountService interface.
type AccountServiceOp struct {
	client *Client
}

var _ AccountService = &AccountServiceOp{}

type rootSubredditKarma struct {
	Kind string           `json:"kind,omitempty"`
	Data []SubredditKarma `json:"data,omitempty"`
}

// SubredditKarma holds user karma data for the subreddit.
type SubredditKarma struct {
	Subreddit    string `json:"sr"`
	PostKarma    int    `json:"link_karma"`
	CommentKarma int    `json:"comment_karma"`
}

// Settings are the user's account settings.
// Some of the fields' descriptions are taken from:
// https://praw.readthedocs.io/en/latest/code_overview/other/preferences.html#praw.models.Preferences.update
// todo: these should probably be pointers with omitempty
type Settings struct {
	AcceptPrivateMessages bool `json:"accept_pms"`
	// Allow Reddit to use your activity on Reddit to show you more relevant advertisements.
	ActivityRelevantAds bool `json:"activity_relevant_ads"`
	// Allow reddit to log my outbound clicks for personalization.
	AllowClickTracking bool `json:"allow_clicktracking"`

	// I would like to beta test features for reddit. By enabling, you will join r/beta immediately.
	Beta                     bool   `json:"beta"`
	ClickGadget              bool   `json:"clickgadget"`
	CollapseReadMessages     bool   `json:"collapse_read_messages"`
	Compress                 bool   `json:"compress"`
	CredditAutorenew         bool   `json:"creddit_autorenew"`
	DefaultCommentSort       string `json:"default_comment_sort"`
	ShowDomainDetails        bool   `json:"domain_details"`
	SendEmailDigests         bool   `json:"email_digests"`
	SendMessagesAsEmails     bool   `json:"email_messages"`
	UnsubscribeFromAllEmails bool   `json:"email_unsubscribe_all"`
	DisableCustomThemes      bool   `json:"enable_default_themes"`
	Location                 string `json:"g"`
	HideAds                  bool   `json:"hide_ads"`

	// Don't allow search engines to index my user profile.
	HideFromSearchEngines bool `json:"hide_from_robots"`

	HideUpvotedPosts   bool `json:"hide_ups"`
	HideDownvotedPosts bool `json:"hide_downs"`

	HighlightControversialComments bool `json:"highlight_controversial"`
	HighlightNewComments           bool `json:"highlight_new_comments"`
	IgnoreSuggestedSorts           bool `json:"ignore_suggested_sort"`
	// Use new Reddit as my default experience.
	UseNewReddit        bool   `json:"in_redesign_beta"`
	LabelNSFW           bool   `json:"label_nsfw"`
	Language            string `json:"lang"`
	ShowOldSearchPage   bool   `json:"legacy_search"`
	EnableNotifications bool   `json:"live_orangereds"`
	MarkMessagesAsRead  bool   `json:"mark_messages_read"`

	// Determine whether to show thumbnails next to posts in subreddits.
	// - "on": show thumbnails next to posts
	// - "off": do not show thumbnails next to posts
	// - "subreddit": show thumbnails next to posts based on the subreddit's preferences
	ShowThumbnails string `json:"media"`

	// Determine whether to auto-expand media in subreddits.
	// - "on": auto-expand media previews
	// - "off": do not auto-expand media previews
	// - "subreddit": auto-expand media previews based on the subreddit's preferences
	AutoExpandMedia            string `json:"media_preview"`
	MinimumCommentScore        *int   `json:"min_comment_score"`
	MinimumPostScore           *int   `json:"min_link_score"`
	EnableMentionNotifications bool   `json:"monitor_mentions"`
	OpenLinksInNewWindow       bool   `json:"newwindow"`
	// todo: test this
	DarkMode         bool `json:"nightmode"`
	DisableProfanity bool `json:"no_profanity"`
	NumComments      int  `json:"num_comments,omitempty"`
	NumPosts         int  `json:"numsites,omitempty"`
	ShowSpotlightBox bool `json:"organic"`
	// todo: test this
	SubredditTheme        string `json:"other_theme"`
	ShowNSFW              bool   `json:"over_18"`
	EnablePrivateRSSFeeds bool   `json:"private_feeds"`
	ProfileOptOut         bool   `json:"profile_opt_out"`
	// Make my upvotes and downvotes public.
	PublicizeVotes bool `json:"public_votes"`

	// Allow my data to be used for research purposes.
	AllowResearch            bool `json:"research"`
	IncludeNSFWSearchResults bool `json:"search_include_over_18"`
	ReceiveCrosspostMessages bool `json:"send_crosspost_messages"`
	ReceiveWelcomeMessages   bool `json:"send_welcome_messages"`

	// Show a user's flair (next to their name on a post or comment).
	ShowUserFlair bool `json:"show_flair"`
	// Show a post's flair.
	ShowPostFlair bool `json:"show_link_flair"`

	ShowGoldExpiration               bool   `json:"show_gold_expiration"`
	ShowLocationBasedRecommendations bool   `json:"show_location_based_recommendations"`
	ShowPromote                      bool   `json:"show_promote"`
	ShowCustomSubredditThemes        bool   `json:"show_stylesheets"`
	ShowTrendingSubreddits           bool   `json:"show_trending"`
	ShowTwitter                      bool   `json:"show_twitter"`
	StoreVisits                      bool   `json:"store_visits"`
	ThemeSelector                    string `json:"theme_selector"`

	// Allow Reddit to use data provided by third-parties to show you more relevant advertisements on Reddit.i
	AllowThirdPartyDataAdPersonalization bool `json:"third_party_data_personalized_ads"`
	// Allow personalization of advertisements using data from third-party websites.
	AllowThirdPartySiteDataAdPersonalization bool `json:"third_party_site_data_personalized_ads"`
	// Allow personalization of content using data from third-party websites.
	AllowThirdPartySiteDataContentPersonalization bool `json:"third_party_site_data_personalized_content"`

	EnableThreadedMessages bool `json:"threaded_messages"`
	EnableThreadedModmail  bool `json:"threaded_modmail"`
	TopKarmaSubreddits     bool `json:"top_karma_subreddits"`
	UseGlobalDefaults      bool `json:"use_global_defaults"`
	// todo: test this, no_video_autoplay
	EnableVideoAutoplay bool `json:"video_autoplay"`
}

// Karma returns a breakdown of your karma per subreddit.
func (s *AccountServiceOp) Karma(ctx context.Context) ([]SubredditKarma, *Response, error) {
	path := "api/v1/me/karma"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootSubredditKarma)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Data, resp, nil
}

// Settings returns your account settings.
func (s *AccountServiceOp) Settings(ctx context.Context) {
	path := "api/v1/me/prefs"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return
	}

	fmt.Println(req)

	root := new(Settings)
	fmt.Println(root.ShowThumbnails)
}
