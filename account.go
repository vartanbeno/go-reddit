package geddit

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// AccountService handles communication with the account
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_account
type AccountService service

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
	// Control whose private messages you see: "everyone" or "whitelisted".
	AcceptPrivateMessages *string `json:"accept_pms,omitempty"`
	// Allow Reddit to use your activity on Reddit to show you more relevant advertisements.
	ActivityRelevantAds *bool `json:"activity_relevant_ads,omitempty"`
	// Allow reddit to log my outbound clicks for personalization.
	AllowClickTracking *bool `json:"allow_clicktracking,omitempty"`

	// Beta test features for reddit. By enabling, you will join r/beta immediately.
	Beta *bool `json:"beta,omitempty"`
	// Show me links I've recently viewed.
	ShowRecentlyViewedPosts *bool `json:"clickgadget,omitempty"`

	CollapseReadMessages *bool `json:"collapse_read_messages,omitempty"`

	// Compress the post display (make them look more compact).
	Compress *bool `json:"compress,omitempty"`

	CredditAutorenew *bool `json:"creddit_autorenew,omitempty"`

	// One of "confidence", "top", "new", "controversial", "old", "random", "qa", "live".
	DefaultCommentSort *string `json:"default_comment_sort,omitempty"`

	// Show additional details in the domain text when available,
	// such as the source subreddit or the content author’s url/name.
	ShowDomainDetails *bool `json:"domain_details,omitempty"`

	SendEmailDigests         *bool `json:"email_digests,omitempty"`
	SendMessagesAsEmails     *bool `json:"email_messages,omitempty"`
	UnsubscribeFromAllEmails *bool `json:"email_unsubscribe_all,omitempty"`

	// Disable subreddits from displaying their custom themes.
	DisableCustomThemes *bool `json:"enable_default_themes,omitempty"`

	// One of "GLOBAL", "AR", "AU", "BG", "CA", "CL", "CO", "CZ", "FI", "GB", "GR", "HR", "HU",
	// "IE", "IN", "IS", "JP", "MX", "MY", "NZ", "PH", "PL", "PR", "PT", "RO", "RS", "SE", "SG",
	// "TH", "TR", "TW", "US", "US_AK", "US_AL", "US_AR", "US_AZ", "US_CA", "US_CO", "US_CT",
	// "US_DC", "US_DE", "US_FL", "US_GA", "US_HI", "US_IA", "US_ID", "US_IL", "US_IN", "US_KS",
	// "US_KY", "US_LA", "US_MA", "US_MD", "US_ME", "US_MI", "US_MN", "US_MO", "US_MS", "US_MT",
	// "US_NC", "US_ND", "US_NE", "US_NH", "US_NJ", "US_NM", "US_NV", "US_NY", "US_OH", "US_OK",
	// "US_OR", "US_PA", "US_RI", "US_SC", "US_SD", "US_TN", "US_TX", "US_UT", "US_VA", "US_VT",
	// "US_WA", "US_WI", "US_WV", "US_WY".
	Location *string `json:"geopopular,omitempty"`

	HideAds *bool `json:"hide_ads,omitempty"`

	// Don't allow search engines to index my user profile.
	HideFromSearchEngines *bool `json:"hide_from_robots,omitempty"`

	// Don’t show me posts after I’ve upvoted them, except my own.
	HideUpvotedPosts *bool `json:"hide_ups,omitempty"`
	// Don’t show me posts after I’ve downvoted them, except my own.
	HideDownvotedPosts *bool `json:"hide_downs,omitempty"`

	// Show a dagger (†) on comments voted controversial (one that's been
	// upvoted and downvoted significantly).
	HighlightControversialComments *bool `json:"highlight_controversial,omitempty"`
	HighlightNewComments           *bool `json:"highlight_new_comments,omitempty"`

	// Ignore suggested sorts for specific threads/subreddits, like Q&As.
	IgnoreSuggestedSorts *bool `json:"ignore_suggested_sort,omitempty"`

	// Use new Reddit as my default experience.
	// Use this to SET the setting.
	UseNewReddit *bool `json:"in_redesign_beta,omitempty"`

	// Use new Reddit as my default experience.
	// Use this to GET the setting.
	UsesNewReddit *bool `json:"design_beta,omitempty"`

	// Label posts that are not safe for work (NSFW).
	LabelNSFW *bool `json:"label_nsfw,omitempty"`

	// A valid IETF language tag (underscore separated).
	Language *string `json:"lang,omitempty"`

	ShowOldSearchPage *bool `json:"legacy_search,omitempty"`

	// Send message notifications in my browser.
	EnableNotifications *bool `json:"live_orangereds,omitempty"`

	MarkMessagesAsRead *bool `json:"mark_messages_read,omitempty"`

	// Determine whether to show thumbnails next to posts in subreddits.
	// - "on": show thumbnails next to posts
	// - "off": do not show thumbnails next to posts
	// - "subreddit": show thumbnails next to posts based on the subreddit's preferences
	ShowThumbnails *string `json:"media,omitempty"`

	// Determine whether to auto-expand media in subreddits.
	// - "on": auto-expand media previews
	// - "off": do not auto-expand media previews
	// - "subreddit": auto-expand media previews based on the subreddit's preferences
	AutoExpandMedia *string `json:"media_preview,omitempty"`

	// Don't show me comments with a score less than this number.
	// Must be between -100 and 100 (inclusive).
	MinimumCommentScore *int `json:"min_comment_score,omitempty"`

	// Don't show me posts with a score less than this number.
	// Must be between -100 and 100 (inclusive).
	MinimumPostScore *int `json:"min_link_score,omitempty"`

	// Notify me when people say my username.
	EnableMentionNotifications *bool `json:"monitor_mentions,omitempty"`

	// Opens link in a new window/tab.
	OpenLinksInNewWindow *bool `json:"newwindow,omitempty"`

	DarkMode         *bool `json:"nightmode,omitempty"`
	DisableProfanity *bool `json:"no_profanity,omitempty"`

	// Display this many comments by default.
	// Must be between 1 and 500 (inclusive).
	NumberOfComments *int `json:"num_comments,omitempty,omitempty"`

	// Display this many posts by default.
	// Must be between 1 and 100 (inclusive).
	NumberOfPosts *int `json:"numsites,omitempty,omitempty"`

	// Show the spotlight box on the home feed.
	// Not sure what this is though...
	ShowSpotlightBox *bool `json:"organic,omitempty"`

	SubredditTheme *string `json:"other_theme,omitempty"`

	// Show content that is labeled not safe for work (NSFW).
	ShowNSFW *bool `json:"over_18,omitempty"`

	EnablePrivateRSSFeeds *bool `json:"private_feeds,omitempty"`

	// View user profiles on desktop using legacy mode.
	ProfileOptOut *bool `json:"profile_opt_out,omitempty"`
	// Make my upvotes and downvotes public.
	PublicizeVotes *bool `json:"public_votes,omitempty"`

	// Allow my data to be used for research purposes.
	AllowResearch *bool `json:"research,omitempty"`

	IncludeNSFWSearchResults *bool `json:"search_include_over_18,omitempty"`

	// Receive a message when my post gets cross-posted.
	ReceiveCrosspostMessages *bool `json:"send_crosspost_messages,omitempty"`
	// Receive welcome messages from moderators when I join a community.
	ReceiveWelcomeMessages *bool `json:"send_welcome_messages,omitempty"`

	// Show a user's flair (next to their name on a post or comment).
	ShowUserFlair *bool `json:"show_flair,omitempty"`
	// Show a post's flair.
	ShowPostFlair *bool `json:"show_link_flair,omitempty"`

	// Show how much gold you have remaining on your profile.
	ShowGoldExpiration               *bool `json:"show_gold_expiration,omitempty"`
	ShowLocationBasedRecommendations *bool `json:"show_location_based_recommendations,omitempty"`
	ShowPromote                      *bool `json:"show_promote,omitempty"`
	ShowCustomSubredditThemes        *bool `json:"show_stylesheets,omitempty"`

	// Show trending subreddits on the home feed.
	ShowTrendingSubreddits *bool   `json:"show_trending,omitempty"`
	ShowTwitter            *bool   `json:"show_twitter,omitempty"`
	StoreVisits            *bool   `json:"store_visits,omitempty"`
	ThemeSelector          *string `json:"theme_selector,omitempty"`

	// Allow Reddit to use data provided by third-parties to show you more relevant advertisements on Reddit.i
	AllowThirdPartyDataAdPersonalization *bool `json:"third_party_data_personalized_ads,omitempty"`
	// Allow personalization of advertisements using data from third-party websites.
	AllowThirdPartySiteDataAdPersonalization *bool `json:"third_party_site_data_personalized_ads,omitempty"`
	// Allow personalization of content using data from third-party websites.
	AllowThirdPartySiteDataContentPersonalization *bool `json:"third_party_site_data_personalized_content,omitempty"`

	EnableThreadedMessages *bool `json:"threaded_messages,omitempty"`
	EnableThreadedModmail  *bool `json:"threaded_modmail,omitempty"`

	// Show the communities you are active in on your profile (mobile only).
	TopKarmaSubreddits *bool `json:"top_karma_subreddits,omitempty"`

	UseGlobalDefaults   *bool `json:"use_global_defaults,omitempty"`
	EnableVideoAutoplay *bool `json:"video_autoplay,omitempty"`
}

// Info returns some general information about your account.
func (s *AccountService) Info(ctx context.Context) (*User, *Response, error) {
	path := "api/v1/me"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(User)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Karma returns a breakdown of your karma per subreddit.
func (s *AccountService) Karma(ctx context.Context) ([]SubredditKarma, *Response, error) {
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
func (s *AccountService) Settings(ctx context.Context) (*Settings, *Response, error) {
	path := "api/v1/me/prefs"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Settings)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// UpdateSettings updates your account settings and returns the modified version.
func (s *AccountService) UpdateSettings(ctx context.Context, settings *Settings) (*Settings, *Response, error) {
	path := "api/v1/me/prefs"

	req, err := s.client.NewRequest(http.MethodPatch, path, settings)
	if err != nil {
		return nil, nil, err
	}

	root := new(Settings)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Trophies returns a list of your trophies.
func (s *AccountService) Trophies(ctx context.Context) ([]Trophy, *Response, error) {
	path := "api/v1/me/trophies"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootTrophyListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	var trophies []Trophy
	for _, trophy := range root.Data.Trophies {
		trophies = append(trophies, trophy.Data)
	}

	return trophies, resp, nil
}

type rootFriendList struct {
	Friends []Friendship
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *rootFriendList) UnmarshalJSON(b []byte) error {
	var resBody []interface{}
	err := json.Unmarshal(b, &resBody)
	if err != nil {
		return err
	}

	if len(resBody) == 0 {
		return errors.New("unexpected data length received")
	}

	data, ok := resBody[0].(map[string]interface{})
	if !ok {
		return errors.New("unexpected data type received")
	}

	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		return errors.New("data does not contain expected field")
	}

	children, ok := dataMap["children"].([]interface{})
	if !ok {
		return errors.New("data does not contain expected field")
	}

	byteValue, err := json.Marshal(children)
	if err != nil {
		return nil
	}

	var friends []Friendship
	err = json.Unmarshal(byteValue, &friends)
	if err != nil {
		return err
	}

	l.Friends = friends
	return nil
}

// Friends returns a list of your friends.
func (s *AccountService) Friends(ctx context.Context) ([]Friendship, *Response, error) {
	path := "prefs/friends"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootFriendList)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Friends, resp, nil
}
