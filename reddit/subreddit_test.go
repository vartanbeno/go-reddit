package reddit

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedPosts = []*Post{
	{
		ID:      "agi5zf",
		FullID:  "t3_agi5zf",
		Created: &Timestamp{time.Date(2019, 1, 16, 5, 57, 51, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/agi5zf/test/",
		URL:       "https://www.reddit.com/r/test/comments/agi5zf/test/",

		Title: "test",
		Body:  "test",

		Score:            253,
		UpvoteRatio:      0.99,
		NumberOfComments: 1634,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8154,

		Author:   "kmiller0112",
		AuthorID: "t2_30a5ktgt",

		IsSelfPost: true,
		Stickied:   true,
	},
	{
		ID:      "hyhquk",
		FullID:  "t3_hyhquk",
		Created: &Timestamp{time.Date(2020, 7, 27, 0, 5, 10, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/test/comments/hyhquk/veggies/",
		URL:       "https://i.imgur.com/LrN2mPw.jpg",

		Title: "Veggies",

		Score:            4,
		UpvoteRatio:      1,
		NumberOfComments: 0,

		SubredditName:         "test",
		SubredditNamePrefixed: "r/test",
		SubredditID:           "t5_2qh23",
		SubredditSubscribers:  8154,

		Author:   "MuckleMcDuckle",
		AuthorID: "t2_6fqntbwq",
	},
}

var expectedSubreddit = &Subreddit{
	ID:      "2rc7j",
	Name:    "t5_2rc7j",
	Created: &Timestamp{time.Date(2009, 11, 11, 0, 54, 28, 0, time.UTC)},

	URL:                 "/r/golang/",
	DisplayName:         "golang",
	DisplayNamePrefixed: "r/golang",
	Title:               "The Go Programming Language",
	Description:         "Ask questions and post articles about the Go programming language and related tools, events etc.",
	SubredditType:       "public",

	Subscribers:      190423,
	ActiveUserCount:  368,
	Over18:           false,
	UserIsModerator:  false,
	UserIsSubscriber: false,
}

var expectedSubreddits = []*Subreddit{
	{
		ID:      "2qs0k",
		Name:    "t5_2qs0k",
		Created: &Timestamp{time.Date(2009, time.January, 25, 10, 25, 57, 0, time.UTC)},

		URL:                 "/r/Home/",
		DisplayName:         "Home",
		DisplayNamePrefixed: "r/Home",
		Title:               "Home",
		SubredditType:       "public",

		Subscribers:      15336,
		Over18:           false,
		UserIsModerator:  false,
		UserIsSubscriber: true,
		UserHasFavorited: false,

		// New
		AcceptFollowers:             false,
		AccountsActive:              0,
		AccountsActiveIsFuzzed:      false,
		ActiveUserCount:             0,
		AdvertiserCategory:          "",
		AllOriginalContent:          false,
		AllowChatPostCreation:       false,
		AllowDiscovery:              true,
		AllowGalleries:              true,
		AllowImages:                 true,
		AllowPolls:                  true,
		AllowPredictionContributors: false,
		AllowPredictions:            false,
		AllowPredictionsTournament:  false,
		AllowTalks:                  false,
		AllowVideogifs:              true,
		AllowVideos:                 true,
		AllowedMediaInComments:      false,
		BannerBackgroundColor:       "",
		BannerBackgroundImage:       "",
		BannerImg:                   "",
		BannerSize:                  []float64(nil),
		CanAssignLinkFlair:          false,
		CanAssignUserFlair:          false,
		CollapseDeletedComments:     false,
		CommentContributionSettings: struct{}{},
		CommentScoreHideMins:        0,
		CommunityIcon:               "",
		CommunityReviewed:           false,
		CreatedUtc:                  &Timestamp{time.Date(2009, time.January, 25, 2, 25, 57, 0, time.UTC)},
		Description:                 "Everything home related: interior design, home improvement, architecture.\n\n**Related subreddits**\n--------------------------\n* [/r/InteriorDesign](http://www.reddit.com/r/interiordesign)\n* [/r/architecture](http://www.reddit.com/r/architecture)\n* [/r/houseporn](http://www.reddit.com/r/houseporn)\n* [/r/roomporn](http://www.reddit.com/r/roomporn)\n* [/r/designmyroom](http://www.reddit.com/r/designmyroom)", DescriptionHTML: "&lt;!-- SC_OFF --&gt;&lt;div class=\"md\"&gt;&lt;p&gt;Everything home related: interior design, home improvement, architecture.&lt;/p&gt;\n\n&lt;h2&gt;&lt;strong&gt;Related subreddits&lt;/strong&gt;&lt;/h2&gt;\n\n&lt;ul&gt;\n&lt;li&gt;&lt;a href=\"http://www.reddit.com/r/interiordesign\"&gt;/r/InteriorDesign&lt;/a&gt;&lt;/li&gt;\n&lt;li&gt;&lt;a href=\"http://www.reddit.com/r/architecture\"&gt;/r/architecture&lt;/a&gt;&lt;/li&gt;\n&lt;li&gt;&lt;a href=\"http://www.reddit.com/r/houseporn\"&gt;/r/houseporn&lt;/a&gt;&lt;/li&gt;\n&lt;li&gt;&lt;a href=\"http://www.reddit.com/r/roomporn\"&gt;/r/roomporn&lt;/a&gt;&lt;/li&gt;\n&lt;li&gt;&lt;a href=\"http://www.reddit.com/r/designmyroom\"&gt;/r/designmyroom&lt;/a&gt;&lt;/li&gt;\n&lt;/ul&gt;\n&lt;/div&gt;&lt;!-- SC_ON --&gt;",
		DisableContributorRequests:       false,
		EmojisCustomSize:                 interface{}(nil),
		EmojisEnabled:                    false,
		FreeFormReports:                  true,
		HasMenuWidget:                    false,
		HeaderImg:                        "",
		HeaderSize:                       []int(nil),
		HeaderTitle:                      "",
		HideAds:                          false,
		IconImg:                          "",
		IconSize:                         []float64(nil),
		IsChatPostFeatureEnabled:         true,
		IsCrosspostableSubreddit:         true,
		IsEnrolledInNewModmail:           false,
		KeyColor:                         "",
		Lang:                             "en",
		LinkFlairEnabled:                 false,
		LinkFlairPosition:                "",
		MobileBannerImage:                "",
		NotificationLevel:                "low",
		OriginalContentTagEnabled:        false,
		PredictionLeaderboardEntryType:   "",
		PrimaryColor:                     "",
		PublicDescription:                "",
		PublicDescriptionHTML:            "",
		PublicTraffic:                    false,
		Quarantine:                       false,
		RestrictCommenting:               false,
		RestrictPosting:                  true,
		ShouldArchivePosts:               false,
		ShouldShowMediaInCommentsSetting: false,
		ShowMedia:                        true,
		ShowMediaPreview:                 true,
		SpoilersEnabled:                  true,
		SubmissionType:                   "any",
		SubmitLinkLabel:                  "",
		SubmitText:                       "",
		SubmitTextHTML:                   "",
		SubmitTextLabel:                  "",
		SuggestedCommentSort:             "",
		UserCanFlairInSr:                 interface{}(nil),
		UserFlairBackgroundColor:         interface{}(nil),
		UserFlairCSSClass:                interface{}(nil),
		UserFlairEnabledInSr:             false,
		UserFlairPosition:                "right",
		UserFlairRichtext:                []interface{}{},
		UserFlairTemplateID:              interface{}(nil),
		UserFlairText:                    interface{}(nil),
		UserFlairTextColor:               interface{}(nil),
		UserFlairType:                    "text",
		UserIsBanned:                     false,
		UserIsContributor:                false,
		UserIsMuted:                      false,
		UserSrFlairEnabled:               false,
		UserSrThemeEnabled:               true,
		WhitelistStatus:                  "all_ads",
		WikiEnabled:                      false,
		Wls:                              6,
	},
	{
		ID:      "2qh1i",
		Name:    "t5_2qh1i",
		Created: &Timestamp{time.Date(2008, time.January, 25, 11, 52, 15, 0, time.UTC)},

		URL:                 "/r/AskReddit/",
		DisplayName:         "AskReddit",
		DisplayNamePrefixed: "r/AskReddit",
		Title:               "Ask Reddit...",
		Description:         "###### [ [ SERIOUS ] ](http://www.reddit.com/r/askreddit/submit?selftext=true&amp;title=%5BSerious%5D)\n\n\n##### [Rules](https://www.reddit.com/r/AskReddit/wiki/index#wiki_rules):\n1. You must post a clear and direct question in the title. The title may contain two, short, necessary context sentences.\nNo text is allowed in the textbox. Your thoughts/responses to the question can go in the comments section. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_1-)\n\n2. Any post asking for advice should be generic and not specific to your situation alone. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_2-)\n\n3. Askreddit is for open-ended discussion questions. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_3-)\n\n4. Posting, or seeking, any identifying personal information, real or fake, will result in a ban without a prior warning. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_4-)\n\n5. Askreddit is not your soapbox, personal army, or advertising platform. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_5-)\n\n6. [Serious] tagged posts are off-limits to jokes or irrelevant replies. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_6-)\n\n7. Soliciting money, goods, services, or favours is not allowed. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_7-)\n\n8. Mods reserve the right to remove content or restrict users' posting privileges as necessary if it is deemed detrimental to the subreddit or to the experience of others. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_8-)\n\n9. Comment replies consisting solely of images will be removed. [more &gt;&gt;](https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_9-)\n\n##### If you think your post has disappeared, see spam or an inappropriate post, please do not hesitate to [contact the mods](https://www.reddit.com/message/compose?to=%2Fr%2FAskReddit), we're happy to help.\n\n---\n\n#### Tags to use:\n\n&gt; ## [[Serious]](https://www.reddit.com/r/AskReddit/wiki/mod_announcements#wiki_.5Bserious.5D_post_tags)\n\n### Use a **[Serious]** post tag to designate your post as a serious, on-topic-only thread.\n\n-\n\n#### Filter posts by subject:\n\n[Mod posts](http://ud.reddit.com/r/AskReddit/#ud)\n[Serious posts](http://dg.reddit.com/r/AskReddit/#dg)\n[Megathread](http://bu.reddit.com/r/AskReddit/#bu)\n[Breaking news](http://nr.reddit.com/r/AskReddit/#nr)\n[Unfilter](/r/AskReddit)\n\n\n-\n\n### Please use spoiler tags to hide spoilers. `&gt;!insert spoiler here!&lt;`\n\n-\n\n#### Other subreddits you might like:\nsome|header\n:---|:---\n[Ask Others](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_ask_others)|[Self &amp; Others](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_self_.26amp.3B_others)\n[Find a subreddit](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_find_a_subreddit)|[Learn something](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_learn_something)\n[Meta Subs](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_meta)|[What is this ___](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_what_is_this______)\n[AskReddit Offshoots](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_askreddit_offshoots)|[Offers &amp; Assistance](https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_offers_.26amp.3B_assistance)\n\n\n-\n\n### Ever read the reddiquette? [Take a peek!](/wiki/reddiquette)\n\n[](#/RES_SR_Config/NightModeCompatible)",
		SubredditType:       "public",

		Subscribers:      28449174,
		Over18:           false,
		UserIsModerator:  false,
		UserIsSubscriber: true,
		UserHasFavorited: true,

		// New
		AcceptFollowers:             false,
		AccountsActive:              0,
		AccountsActiveIsFuzzed:      false,
		ActiveUserCount:             0,
		AdvertiserCategory:          "Lifestyles",
		AllOriginalContent:          false,
		AllowChatPostCreation:       false,
		AllowDiscovery:              true,
		AllowGalleries:              true,
		AllowImages:                 false,
		AllowPolls:                  false,
		AllowPredictionContributors: false,
		AllowPredictions:            false,
		AllowPredictionsTournament:  false,
		AllowTalks:                  false,
		AllowVideogifs:              false,
		AllowVideos:                 false,
		AllowedMediaInComments:      false,
		BannerBackgroundColor:       "#f0f7fd",
		BannerBackgroundImage:       "",
		BannerImg:                   "https://b.thumbs.redditmedia.com/PXt8GnqdYu-9lgzb3iesJBLN21bXExRV1A45zdw4sYE.png",
		BannerSize: []float64{1280,
			384},
		CanAssignLinkFlair:          false,
		CanAssignUserFlair:          false,
		CollapseDeletedComments:     true,
		CommentContributionSettings: struct{}{},
		CommentScoreHideMins:        60,
		CommunityIcon:               "https://styles.redditmedia.com/t5_2qh1i/styles/communityIcon_tijjpyw1qe201.png?width=256&amp;s=4e76eadc662b8155a93d4d7487a6d3acb35f4334",
		CommunityReviewed:           false,
		CreatedUtc:                  &Timestamp{time.Date(2008, time.January, 25, 3, 52, 15, 0, time.UTC)},
		DescriptionHTML:             "&lt;!-- SC_OFF --&gt;&lt;div class=\"md\"&gt;&lt;h6&gt;&lt;a href=\"http://www.reddit.com/r/askreddit/submit?selftext=true&amp;amp;title=%5BSerious%5D\"&gt; [ SERIOUS ] &lt;/a&gt;&lt;/h6&gt;\n\n&lt;h5&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_rules\"&gt;Rules&lt;/a&gt;:&lt;/h5&gt;\n\n&lt;ol&gt;\n&lt;li&gt;&lt;p&gt;You must post a clear and direct question in the title. The title may contain two, short, necessary context sentences.\nNo text is allowed in the textbox. Your thoughts/responses to the question can go in the comments section. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_1-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;Any post asking for advice should be generic and not specific to your situation alone. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_2-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;Askreddit is for open-ended discussion questions. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_3-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;Posting, or seeking, any identifying personal information, real or fake, will result in a ban without a prior warning. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_4-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;Askreddit is not your soapbox, personal army, or advertising platform. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_5-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;[Serious] tagged posts are off-limits to jokes or irrelevant replies. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_6-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;Soliciting money, goods, services, or favours is not allowed. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_7-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;Mods reserve the right to remove content or restrict users&amp;#39; posting privileges as necessary if it is deemed detrimental to the subreddit or to the experience of others. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_8-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;Comment replies consisting solely of images will be removed. &lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/index#wiki_-rule_9-\"&gt;more &amp;gt;&amp;gt;&lt;/a&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;/ol&gt;\n\n&lt;h5&gt;If you think your post has disappeared, see spam or an inappropriate post, please do not hesitate to &lt;a href=\"https://www.reddit.com/message/compose?to=%2Fr%2FAskReddit\"&gt;contact the mods&lt;/a&gt;, we&amp;#39;re happy to help.&lt;/h5&gt;\n\n&lt;hr/&gt;\n\n&lt;h4&gt;Tags to use:&lt;/h4&gt;\n\n&lt;blockquote&gt;\n&lt;h2&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/mod_announcements#wiki_.5Bserious.5D_post_tags\"&gt;[Serious]&lt;/a&gt;&lt;/h2&gt;\n&lt;/blockquote&gt;\n\n&lt;h3&gt;Use a &lt;strong&gt;[Serious]&lt;/strong&gt; post tag to designate your post as a serious, on-topic-only thread.&lt;/h3&gt;\n\n&lt;h2&gt;&lt;/h2&gt;\n\n&lt;h4&gt;Filter posts by subject:&lt;/h4&gt;\n\n&lt;p&gt;&lt;a href=\"http://ud.reddit.com/r/AskReddit/#ud\"&gt;Mod posts&lt;/a&gt;\n&lt;a href=\"http://dg.reddit.com/r/AskReddit/#dg\"&gt;Serious posts&lt;/a&gt;\n&lt;a href=\"http://bu.reddit.com/r/AskReddit/#bu\"&gt;Megathread&lt;/a&gt;\n&lt;a href=\"http://nr.reddit.com/r/AskReddit/#nr\"&gt;Breaking news&lt;/a&gt;\n&lt;a href=\"/r/AskReddit\"&gt;Unfilter&lt;/a&gt;&lt;/p&gt;\n\n&lt;h2&gt;&lt;/h2&gt;\n\n&lt;h3&gt;Please use spoiler tags to hide spoilers. &lt;code&gt;&amp;gt;!insert spoiler here!&amp;lt;&lt;/code&gt;&lt;/h3&gt;\n\n&lt;h2&gt;&lt;/h2&gt;\n\n&lt;h4&gt;Other subreddits you might like:&lt;/h4&gt;\n\n&lt;table&gt;&lt;thead&gt;\n&lt;tr&gt;\n&lt;th align=\"left\"&gt;some&lt;/th&gt;\n&lt;th align=\"left\"&gt;header&lt;/th&gt;\n&lt;/tr&gt;\n&lt;/thead&gt;&lt;tbody&gt;\n&lt;tr&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_ask_others\"&gt;Ask Others&lt;/a&gt;&lt;/td&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_self_.26amp.3B_others\"&gt;Self &amp;amp; Others&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_find_a_subreddit\"&gt;Find a subreddit&lt;/a&gt;&lt;/td&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_learn_something\"&gt;Learn something&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_meta\"&gt;Meta Subs&lt;/a&gt;&lt;/td&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_what_is_this______\"&gt;What is this ___&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_askreddit_offshoots\"&gt;AskReddit Offshoots&lt;/a&gt;&lt;/td&gt;\n&lt;td align=\"left\"&gt;&lt;a href=\"https://www.reddit.com/r/AskReddit/wiki/sidebarsubs#wiki_offers_.26amp.3B_assistance\"&gt;Offers &amp;amp; Assistance&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;/tbody&gt;&lt;/table&gt;\n\n&lt;h2&gt;&lt;/h2&gt;\n\n&lt;h3&gt;Ever read the reddiquette? &lt;a href=\"/wiki/reddiquette\"&gt;Take a peek!&lt;/a&gt;&lt;/h3&gt;\n\n&lt;p&gt;&lt;a href=\"#/RES_SR_Config/NightModeCompatible\"&gt;&lt;/a&gt;&lt;/p&gt;\n&lt;/div&gt;&lt;!-- SC_ON --&gt;",
		DisableContributorRequests:  false,
		EmojisCustomSize:            interface{}(nil),
		EmojisEnabled:               true,
		FreeFormReports:             true,
		HasMenuWidget:               false,
		HeaderImg:                   "https://a.thumbs.redditmedia.com/IrfPJGuWzi_ewrDTBlnULeZsJYGz81hsSQoQJyw6LD8.png",
		HeaderSize: []int{125,
			73},
		HeaderTitle: "Ass Credit",
		HideAds:     false,
		IconImg:     "https://b.thumbs.redditmedia.com/EndDxMGB-FTZ2MGtjepQ06cQEkZw_YQAsOUudpb9nSQ.png",
		IconSize: []float64{256,
			256},
		IsChatPostFeatureEnabled:         false,
		IsCrosspostableSubreddit:         false,
		IsEnrolledInNewModmail:           false,
		KeyColor:                         "#222222",
		Lang:                             "es",
		LinkFlairEnabled:                 true,
		LinkFlairPosition:                "right",
		MobileBannerImage:                "",
		NotificationLevel:                "low",
		OriginalContentTagEnabled:        false,
		PredictionLeaderboardEntryType:   "",
		PrimaryColor:                     "#646d73",
		PublicDescription:                "r/AskReddit is the place to ask and answer thought-provoking questions.",
		PublicDescriptionHTML:            "&lt;!-- SC_OFF --&gt;&lt;div class=\"md\"&gt;&lt;p&gt;&lt;a href=\"/r/AskReddit\"&gt;r/AskReddit&lt;/a&gt; is the place to ask and answer thought-provoking questions.&lt;/p&gt;\n&lt;/div&gt;&lt;!-- SC_ON --&gt;",
		PublicTraffic:                    false,
		Quarantine:                       false,
		RestrictCommenting:               false,
		RestrictPosting:                  true,
		ShouldArchivePosts:               false,
		ShouldShowMediaInCommentsSetting: false,
		ShowMedia:                        false,
		ShowMediaPreview:                 true,
		SpoilersEnabled:                  true,
		SubmissionType:                   "self",
		SubmitLinkLabel:                  "",
		SubmitText:                       "**AskReddit is all about DISCUSSION. Your post needs to inspire discussion, ask an open-ended question that prompts redditors to share ideas or opinions.**\n\n**Questions need to be neutral and the question alone.** Any opinion or answer must go as a reply to your question, this includes examples or any kind of story about you. This is so that all responses will be to your question, and there's nothing else to respond to. Opinionated posts are forbidden.\n\n* If your question has a factual answer, try r/answers.\n* If you are trying to find out about something or get an explanation, try r/explainlikeimfive\n* If your question has a limited number of responses, then it's not suitable.\n* If you're asking for any kind of advice, then it's not suitable.\n* If you feel the need to add an example in order for your question to make sense then you need to re-word your question.\n* If you're explaining why you're asking the question, you need to stop.\n\nYou can always ask where to post in r/findareddit.",
		SubmitTextHTML:                   "&lt;!-- SC_OFF --&gt;&lt;div class=\"md\"&gt;&lt;p&gt;&lt;strong&gt;AskReddit is all about DISCUSSION. Your post needs to inspire discussion, ask an open-ended question that prompts redditors to share ideas or opinions.&lt;/strong&gt;&lt;/p&gt;\n\n&lt;p&gt;&lt;strong&gt;Questions need to be neutral and the question alone.&lt;/strong&gt; Any opinion or answer must go as a reply to your question, this includes examples or any kind of story about you. This is so that all responses will be to your question, and there&amp;#39;s nothing else to respond to. Opinionated posts are forbidden.&lt;/p&gt;\n\n&lt;ul&gt;\n&lt;li&gt;If your question has a factual answer, try &lt;a href=\"/r/answers\"&gt;r/answers&lt;/a&gt;.&lt;/li&gt;\n&lt;li&gt;If you are trying to find out about something or get an explanation, try &lt;a href=\"/r/explainlikeimfive\"&gt;r/explainlikeimfive&lt;/a&gt;&lt;/li&gt;\n&lt;li&gt;If your question has a limited number of responses, then it&amp;#39;s not suitable.&lt;/li&gt;\n&lt;li&gt;If you&amp;#39;re asking for any kind of advice, then it&amp;#39;s not suitable.&lt;/li&gt;\n&lt;li&gt;If you feel the need to add an example in order for your question to make sense then you need to re-word your question.&lt;/li&gt;\n&lt;li&gt;If you&amp;#39;re explaining why you&amp;#39;re asking the question, you need to stop.&lt;/li&gt;\n&lt;/ul&gt;\n\n&lt;p&gt;You can always ask where to post in &lt;a href=\"/r/findareddit\"&gt;r/findareddit&lt;/a&gt;.&lt;/p&gt;\n&lt;/div&gt;&lt;!-- SC_ON --&gt;",
		SubmitTextLabel:                  "Ask a question",
		SuggestedCommentSort:             "",
		UserCanFlairInSr:                 interface{}(nil),
		UserFlairBackgroundColor:         interface{}(nil),
		UserFlairCSSClass:                interface{}(nil),
		UserFlairEnabledInSr:             false,
		UserFlairPosition:                "right",
		UserFlairRichtext:                []interface{}{},
		UserFlairTemplateID:              interface{}(nil),
		UserFlairText:                    interface{}(nil),
		UserFlairTextColor:               interface{}(nil),
		UserFlairType:                    "text",
		UserIsBanned:                     false,
		UserIsContributor:                false,
		UserIsMuted:                      false,
		UserSrFlairEnabled:               false,
		UserSrThemeEnabled:               true,
		WhitelistStatus:                  "all_ads",
		WikiEnabled:                      true,
		Wls:                              6,
	},
	{
		ID:      "2qh0u",
		Name:    "t5_2qh0u",
		Created: &Timestamp{time.Date(2008, time.January, 25, 8, 31, 9, 0, time.UTC)},

		URL:                 "/r/pics/",
		DisplayName:         "pics",
		DisplayNamePrefixed: "r/pics",
		Title:               "Reddit Pics",
		Description:         "A place to share photographs and pictures. Feel free to post your own, but please **read the rules first** (see below), and note that we are *not a catch-all* for ALL images (of screenshots, comics, etc.).\n\n---\n\n#Spoiler code#\n\nPlease mark spoilers like this:  \n`&gt;!text here!&lt;`\n\nClick/tap to &gt;!read!&lt;.\n\n---\nCheck out http://nt.reddit.com/r/pics!\n\nCheck out /r/pics/wiki/v2/resources/takedown for help with taking down posts due to copyright or personal identifiable information reasons. \n\n---\n#[Posting Rules](/r/pics/wiki/index)#\n\n1. (1A) **No screenshots or pics where the only focus is a screen.**\n\n (1B) No pictures with added or superimposed **digital text, emojis, and \"MS Paint\"-like scribbles.** Exceptions to this rule include watermarks serving to credit the original author, and blurring/boxing out of personal information. \"Photoshopped\" or otherwise manipulated images are allowed.\n\n1. **No porn or gore.** Artistic nudity is allowed. NSFW comments must be tagged. Posting gratuitous materials may result in an immediate and permanent ban.\n\n1. **No personal information, in posts or comments.** No direct links to any Social Media. No subreddit-related meta-drama or witch-hunts. No Missing/Found posts for people or property.  A license plate is not PI. [**Reddit Policy**](https://www.reddithelp.com/en/categories/rules-reporting/account-and-community-restrictions/posting-someones-private-or-personal) \n\n **Stalking, harassment, witch hunting, or doxxing** will not be tolerated and will result in a ban.\n\n **No subreddit-related meta-drama or witch-hunts.**\n\n1. **Titles must follow all [title guidelines](https://www.reddit.com/r/pics/wiki/titles).**\n\n1. **Submissions must link directly to a specific image file or to an image hosting website with minimal ads.** *We do not allow blog hosting of images (\"blogspam\"), but links to albums on image hosting websites are okay. URL shorteners are prohibited. URLs in image or album descriptions are prohibited.* \n\n1. **No animated images.** *Please submit them to /r/gif, /r/gifs, or /r/reactiongifs instead.*\n\n1. We enforce a standard of common decency and civility here. **Please be respectful to others.** Personal attacks, bigotry, fighting words, otherwise inappropriate behavior or content, comments that insult or demean a specific user or group of users will be removed. Regular or egregious violations will result in a ban.\n**Optimally**, the level of discourse here should be at the level you'd find between you and your teacher, or between you and professional colleagues.  Obviously we're going to allow various types of humor here, but if it would make someone you respect lose respect for you, then you're best off avoiding it.\n\n1.  **No submissions featuring before-and-after depictions of personal health progress or achievement. Standalone images of medals, tokens, certificates, and awards are similarly disallowed, save for when the items are being presented as historical curiosities.**\n\n\n1. **No false claims of ownership (FCoO) or flooding.** False claims of ownership (FCoO) and/or flooding (*more than four posts in twenty-four hours*) will result in a ban.\n\n\n1. **Reposts of images on the front page, or within the set limit of /r/pics/top, will be removed.** \n\n (10A) Reposts of images currently on the front page of /r/Pics will be removed.\n\n (10B) Reposts of the top 25 images this year, and top 50 of \"all time\" will be removed.\n\n1. **Only one self-promotional link per post.** Content creators are only allowed one link per post. Anything more may result in temporary or permanent bans. Accounts that exist solely to advertise or promote will be banned.\n\n---\n\n**Loose-ends**\n\n* Serial reposters may be filtered or banned. \n\n---\n\nIf you come across any rule violations please report the submission or  [message the mods](http://www.reddit.com/message/compose?to=%23pics) and one of us will remove it!\n\n  \nIf your submission appears to be filtered, but **definitely** meets the above rules, [please send us a message](/message/compose?to=%23pics) with a link to the **comments section** of your post (not a direct link to the image). **Don't delete it**  as that just makes the filter hate you! \n\n---\n\n\n#Links#\nIf your post doesn't meet the above rules, consider submitting it on one of these other subreddits:\n\n#Subreddits\nBelow is a table of subreddits that you might want to check out!\n\nScreenshots | Advice Animals\n-----------|--------------\n/r/images | /r/adviceanimals\n/r/screenshots | /r/memes\n/r/desktops | /r/memesIRL\n/r/amoledbackgrounds | /r/wholesomememes \n**Animals** | **More Animals**\n/r/aww | /r/fawns\n/r/dogs | /r/rabbits\n/r/cats | /r/RealLifePokemon\n/r/foxes | /r/BeforeNAfterAdoption\n**GIFS** | **HQ / Curated**\n/r/gifs | /r/pic\n/r/catgifs | /r/earthporn\n/r/reactiongifs | /r/spaceporn\n\n##Topic subreddits\n\nEvery now and then, we choose 2 new topics, and find some subreddits about that topic to feature!\n\nOne Word | Art\n-----|----------\n/r/catsstandingup | /r/Art\n/r/nocontextpics | /r/ImaginaryBestOf\n&amp;nbsp; | /r/IDAP", DescriptionHTML: "&lt;!-- SC_OFF --&gt;&lt;div class=\"md\"&gt;&lt;p&gt;A place to share photographs and pictures. Feel free to post your own, but please &lt;strong&gt;read the rules first&lt;/strong&gt; (see below), and note that we are &lt;em&gt;not a catch-all&lt;/em&gt; for ALL images (of screenshots, comics, etc.).&lt;/p&gt;\n\n&lt;hr/&gt;\n\n&lt;h1&gt;Spoiler code&lt;/h1&gt;\n\n&lt;p&gt;Please mark spoilers like this:&lt;br/&gt;\n&lt;code&gt;&amp;gt;!text here!&amp;lt;&lt;/code&gt;&lt;/p&gt;\n\n&lt;p&gt;Click/tap to &lt;span class=\"md-spoiler-text\"&gt;read&lt;/span&gt;.&lt;/p&gt;\n\n&lt;hr/&gt;\n\n&lt;p&gt;Check out &lt;a href=\"http://nt.reddit.com/r/pics\"&gt;http://nt.reddit.com/r/pics&lt;/a&gt;!&lt;/p&gt;\n\n&lt;p&gt;Check out &lt;a href=\"/r/pics/wiki/v2/resources/takedown\"&gt;/r/pics/wiki/v2/resources/takedown&lt;/a&gt; for help with taking down posts due to copyright or personal identifiable information reasons. &lt;/p&gt;\n\n&lt;hr/&gt;\n\n&lt;h1&gt;&lt;a href=\"/r/pics/wiki/index\"&gt;Posting Rules&lt;/a&gt;&lt;/h1&gt;\n\n&lt;ol&gt;\n&lt;li&gt;&lt;p&gt;(1A) &lt;strong&gt;No screenshots or pics where the only focus is a screen.&lt;/strong&gt;&lt;/p&gt;\n\n&lt;p&gt;(1B) No pictures with added or superimposed &lt;strong&gt;digital text, emojis, and &amp;quot;MS Paint&amp;quot;-like scribbles.&lt;/strong&gt; Exceptions to this rule include watermarks serving to credit the original author, and blurring/boxing out of personal information. &amp;quot;Photoshopped&amp;quot; or otherwise manipulated images are allowed.&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;No porn or gore.&lt;/strong&gt; Artistic nudity is allowed. NSFW comments must be tagged. Posting gratuitous materials may result in an immediate and permanent ban.&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;No personal information, in posts or comments.&lt;/strong&gt; No direct links to any Social Media. No subreddit-related meta-drama or witch-hunts. No Missing/Found posts for people or property.  A license plate is not PI. &lt;a href=\"https://www.reddithelp.com/en/categories/rules-reporting/account-and-community-restrictions/posting-someones-private-or-personal\"&gt;&lt;strong&gt;Reddit Policy&lt;/strong&gt;&lt;/a&gt; &lt;/p&gt;\n\n&lt;p&gt;&lt;strong&gt;Stalking, harassment, witch hunting, or doxxing&lt;/strong&gt; will not be tolerated and will result in a ban.&lt;/p&gt;\n\n&lt;p&gt;&lt;strong&gt;No subreddit-related meta-drama or witch-hunts.&lt;/strong&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;Titles must follow all &lt;a href=\"https://www.reddit.com/r/pics/wiki/titles\"&gt;title guidelines&lt;/a&gt;.&lt;/strong&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;Submissions must link directly to a specific image file or to an image hosting website with minimal ads.&lt;/strong&gt; &lt;em&gt;We do not allow blog hosting of images (&amp;quot;blogspam&amp;quot;), but links to albums on image hosting websites are okay. URL shorteners are prohibited. URLs in image or album descriptions are prohibited.&lt;/em&gt; &lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;No animated images.&lt;/strong&gt; &lt;em&gt;Please submit them to &lt;a href=\"/r/gif\"&gt;/r/gif&lt;/a&gt;, &lt;a href=\"/r/gifs\"&gt;/r/gifs&lt;/a&gt;, or &lt;a href=\"/r/reactiongifs\"&gt;/r/reactiongifs&lt;/a&gt; instead.&lt;/em&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;We enforce a standard of common decency and civility here. &lt;strong&gt;Please be respectful to others.&lt;/strong&gt; Personal attacks, bigotry, fighting words, otherwise inappropriate behavior or content, comments that insult or demean a specific user or group of users will be removed. Regular or egregious violations will result in a ban.\n&lt;strong&gt;Optimally&lt;/strong&gt;, the level of discourse here should be at the level you&amp;#39;d find between you and your teacher, or between you and professional colleagues.  Obviously we&amp;#39;re going to allow various types of humor here, but if it would make someone you respect lose respect for you, then you&amp;#39;re best off avoiding it.&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;No submissions featuring before-and-after depictions of personal health progress or achievement. Standalone images of medals, tokens, certificates, and awards are similarly disallowed, save for when the items are being presented as historical curiosities.&lt;/strong&gt;&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;No false claims of ownership (FCoO) or flooding.&lt;/strong&gt; False claims of ownership (FCoO) and/or flooding (&lt;em&gt;more than four posts in twenty-four hours&lt;/em&gt;) will result in a ban.&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;Reposts of images on the front page, or within the set limit of &lt;a href=\"/r/pics/top\"&gt;/r/pics/top&lt;/a&gt;, will be removed.&lt;/strong&gt; &lt;/p&gt;\n\n&lt;p&gt;(10A) Reposts of images currently on the front page of &lt;a href=\"/r/Pics\"&gt;/r/Pics&lt;/a&gt; will be removed.&lt;/p&gt;\n\n&lt;p&gt;(10B) Reposts of the top 25 images this year, and top 50 of &amp;quot;all time&amp;quot; will be removed.&lt;/p&gt;&lt;/li&gt;\n&lt;li&gt;&lt;p&gt;&lt;strong&gt;Only one self-promotional link per post.&lt;/strong&gt; Content creators are only allowed one link per post. Anything more may result in temporary or permanent bans. Accounts that exist solely to advertise or promote will be banned.&lt;/p&gt;&lt;/li&gt;\n&lt;/ol&gt;\n\n&lt;hr/&gt;\n\n&lt;p&gt;&lt;strong&gt;Loose-ends&lt;/strong&gt;&lt;/p&gt;\n\n&lt;ul&gt;\n&lt;li&gt;Serial reposters may be filtered or banned. &lt;/li&gt;\n&lt;/ul&gt;\n\n&lt;hr/&gt;\n\n&lt;p&gt;If you come across any rule violations please report the submission or  &lt;a href=\"http://www.reddit.com/message/compose?to=%23pics\"&gt;message the mods&lt;/a&gt; and one of us will remove it!&lt;/p&gt;\n\n&lt;p&gt;If your submission appears to be filtered, but &lt;strong&gt;definitely&lt;/strong&gt; meets the above rules, &lt;a href=\"/message/compose?to=%23pics\"&gt;please send us a message&lt;/a&gt; with a link to the &lt;strong&gt;comments section&lt;/strong&gt; of your post (not a direct link to the image). &lt;strong&gt;Don&amp;#39;t delete it&lt;/strong&gt;  as that just makes the filter hate you! &lt;/p&gt;\n\n&lt;hr/&gt;\n\n&lt;h1&gt;Links&lt;/h1&gt;\n\n&lt;p&gt;If your post doesn&amp;#39;t meet the above rules, consider submitting it on one of these other subreddits:&lt;/p&gt;\n\n&lt;h1&gt;Subreddits&lt;/h1&gt;\n\n&lt;p&gt;Below is a table of subreddits that you might want to check out!&lt;/p&gt;\n\n&lt;table&gt;&lt;thead&gt;\n&lt;tr&gt;\n&lt;th&gt;Screenshots&lt;/th&gt;\n&lt;th&gt;Advice Animals&lt;/th&gt;\n&lt;/tr&gt;\n&lt;/thead&gt;&lt;tbody&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/images\"&gt;/r/images&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/adviceanimals\"&gt;/r/adviceanimals&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/screenshots\"&gt;/r/screenshots&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/memes\"&gt;/r/memes&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/desktops\"&gt;/r/desktops&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/memesIRL\"&gt;/r/memesIRL&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/amoledbackgrounds\"&gt;/r/amoledbackgrounds&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/wholesomememes\"&gt;/r/wholesomememes&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;strong&gt;Animals&lt;/strong&gt;&lt;/td&gt;\n&lt;td&gt;&lt;strong&gt;More Animals&lt;/strong&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/aww\"&gt;/r/aww&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/fawns\"&gt;/r/fawns&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/dogs\"&gt;/r/dogs&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/rabbits\"&gt;/r/rabbits&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/cats\"&gt;/r/cats&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/RealLifePokemon\"&gt;/r/RealLifePokemon&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/foxes\"&gt;/r/foxes&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/BeforeNAfterAdoption\"&gt;/r/BeforeNAfterAdoption&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;strong&gt;GIFS&lt;/strong&gt;&lt;/td&gt;\n&lt;td&gt;&lt;strong&gt;HQ / Curated&lt;/strong&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/gifs\"&gt;/r/gifs&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/pic\"&gt;/r/pic&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/catgifs\"&gt;/r/catgifs&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/earthporn\"&gt;/r/earthporn&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/reactiongifs\"&gt;/r/reactiongifs&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/spaceporn\"&gt;/r/spaceporn&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;/tbody&gt;&lt;/table&gt;\n\n&lt;h2&gt;Topic subreddits&lt;/h2&gt;\n\n&lt;p&gt;Every now and then, we choose 2 new topics, and find some subreddits about that topic to feature!&lt;/p&gt;\n\n&lt;table&gt;&lt;thead&gt;\n&lt;tr&gt;\n&lt;th&gt;One Word&lt;/th&gt;\n&lt;th&gt;Art&lt;/th&gt;\n&lt;/tr&gt;\n&lt;/thead&gt;&lt;tbody&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/catsstandingup\"&gt;/r/catsstandingup&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/Art\"&gt;/r/Art&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&lt;a href=\"/r/nocontextpics\"&gt;/r/nocontextpics&lt;/a&gt;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/ImaginaryBestOf\"&gt;/r/ImaginaryBestOf&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;tr&gt;\n&lt;td&gt;&amp;nbsp;&lt;/td&gt;\n&lt;td&gt;&lt;a href=\"/r/IDAP\"&gt;/r/IDAP&lt;/a&gt;&lt;/td&gt;\n&lt;/tr&gt;\n&lt;/tbody&gt;&lt;/table&gt;\n&lt;/div&gt;&lt;!-- SC_ON --&gt;",
		SubredditType: "public",

		Subscribers:      24987753,
		Over18:           false,
		UserIsModerator:  false,
		UserIsSubscriber: false,
		UserHasFavorited: false,

		// New
		AcceptFollowers:             false,
		AccountsActive:              0,
		AccountsActiveIsFuzzed:      false,
		ActiveUserCount:             0,
		AdvertiserCategory:          "Lifestyles",
		AllOriginalContent:          false,
		AllowChatPostCreation:       false,
		AllowDiscovery:              true,
		AllowGalleries:              true,
		AllowImages:                 true,
		AllowPolls:                  false,
		AllowPredictionContributors: false,
		AllowPredictions:            false,
		AllowPredictionsTournament:  false,
		AllowTalks:                  false,
		AllowVideogifs:              true,
		AllowVideos:                 false,
		AllowedMediaInComments:      false,
		BannerBackgroundColor:       "#5a74cc",
		BannerBackgroundImage:       "",
		BannerImg:                   "",
		BannerSize:                  []float64(nil),
		CanAssignLinkFlair:          false,
		CanAssignUserFlair:          false,
		CollapseDeletedComments:     true,
		CommentContributionSettings: struct{}{},
		CommentScoreHideMins:        60,
		CommunityIcon:               "",
		CommunityReviewed:           false,
		CreatedUtc:                  &Timestamp{time.Date(2008, time.January, 25, 0, 31, 9, 0, time.UTC)},
		DisableContributorRequests:  false,
		EmojisCustomSize:            interface{}(nil),
		EmojisEnabled:               false,
		FreeFormReports:             true,
		HasMenuWidget:               false,
		HeaderImg:                   "https://b.thumbs.redditmedia.com/1zT3FeN8pCAFIooNVuyuZ0ObU0x1ro4wPfArGHl3KjM.png",
		HeaderSize: []int{160,
			64},
		HeaderTitle: "Something Clever",
		HideAds:     false,
		IconImg:     "https://b.thumbs.redditmedia.com/VZX_KQLnI1DPhlEZ07bIcLzwR1Win808RIt7zm49VIQ.png",
		IconSize: []float64{256,
			256},
		IsChatPostFeatureEnabled:         false,
		IsCrosspostableSubreddit:         true,
		IsEnrolledInNewModmail:           false,
		KeyColor:                         "#222222",
		Lang:                             "en",
		LinkFlairEnabled:                 true,
		LinkFlairPosition:                "left",
		MobileBannerImage:                "",
		NotificationLevel:                interface{}(nil),
		OriginalContentTagEnabled:        true,
		PredictionLeaderboardEntryType:   "",
		PrimaryColor:                     "#cee3f8",
		PublicDescription:                "A place for pictures and photographs.",
		PublicDescriptionHTML:            "&lt;!-- SC_OFF --&gt;&lt;div class=\"md\"&gt;&lt;p&gt;A place for pictures and photographs.&lt;/p&gt;\n&lt;/div&gt;&lt;!-- SC_ON --&gt;",
		PublicTraffic:                    false,
		Quarantine:                       false,
		RestrictCommenting:               false,
		RestrictPosting:                  true,
		ShouldArchivePosts:               false,
		ShouldShowMediaInCommentsSetting: false,
		ShowMedia:                        true,
		ShowMediaPreview:                 true,
		SpoilersEnabled:                  true,
		SubmissionType:                   "link",
		SubmitLinkLabel:                  "Submit an image",
		SubmitText:                       "Please read [the sidebar](/r/pics/about/sidebar) before submitting, and know that by posting you are agreeing to follow those rules.\nLimit: 100 characters", SubmitTextHTML: "&lt;!-- SC_OFF --&gt;&lt;div class=\"md\"&gt;&lt;p&gt;Please read &lt;a href=\"/r/pics/about/sidebar\"&gt;the sidebar&lt;/a&gt; before submitting, and know that by posting you are agreeing to follow those rules.\nLimit: 100 characters&lt;/p&gt;\n&lt;/div&gt;&lt;!-- SC_ON --&gt;",
		SubmitTextLabel:          "",
		SuggestedCommentSort:     "",
		UserCanFlairInSr:         interface{}(nil),
		UserFlairBackgroundColor: interface{}(nil),
		UserFlairCSSClass:        interface{}(nil),
		UserFlairEnabledInSr:     false,
		UserFlairPosition:        "right",
		UserFlairRichtext:        []interface{}{},
		UserFlairTemplateID:      interface{}(nil),
		UserFlairText:            interface{}(nil),
		UserFlairTextColor:       interface{}(nil),
		UserFlairType:            "text",
		UserIsBanned:             false,
		UserIsContributor:        false,
		UserIsMuted:              false,
		UserSrFlairEnabled:       false,
		UserSrThemeEnabled:       true,
		WhitelistStatus:          "all_ads",
		WikiEnabled:              true,
		Wls:                      6,
	},
}

var expectedSubredditNames = []string{
	"golang",
	"golang_infosec",
	"GolangJobOfferings",
	"golanguage",
	"golang2",
	"GolangUnofficial",
	"golang_jobs",
}

var expectedSearchPosts = []*Post{
	{
		ID:      "hybow9",
		FullID:  "t3_hybow9",
		Created: &Timestamp{time.Date(2020, 7, 26, 18, 14, 24, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/WatchPeopleDieInside/comments/hybow9/pregnancy_test/",
		URL:       "https://v.redd.it/ra4qnt8bt8d51",

		Title: "Pregnancy test",

		Score:            103829,
		UpvoteRatio:      0.88,
		NumberOfComments: 3748,

		SubredditName:         "WatchPeopleDieInside",
		SubredditNamePrefixed: "r/WatchPeopleDieInside",
		SubredditID:           "t5_3h4zq",
		SubredditSubscribers:  2599948,

		Author:   "chocolat_ice_cream",
		AuthorID: "t2_3p32m02",
	},
	{
		ID:      "hmwhd7",
		FullID:  "t3_hmwhd7",
		Created: &Timestamp{time.Date(2020, 7, 7, 15, 19, 42, 0, time.UTC)},
		Edited:  &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},

		Permalink: "/r/worldnews/comments/hmwhd7/brazilian_president_jair_bolsonaro_tests_positive/",
		URL:       "https://www.theguardian.com/world/2020/jul/07/jair-bolsonaro-coronavirus-positive-test-brazil-president",

		Title: "Brazilian president Jair Bolsonaro tests positive for coronavirus",

		Score:            149238,
		UpvoteRatio:      0.94,
		NumberOfComments: 7415,

		SubredditName:         "worldnews",
		SubredditNamePrefixed: "r/worldnews",
		SubredditID:           "t5_2qh13",
		SubredditSubscribers:  24651441,

		Author:   "Jeremy_Martin",
		AuthorID: "t2_wgrkg",
	},
}

var expectedRandomSubreddit = &Subreddit{
	Name:                "t5_2wi4l",
	Created:             &Timestamp{time.Date(2013, time.March, 1, 12, 4, 18, 0, time.UTC)},
	URL:                 "/r/GalaxyS8/",
	DisplayName:         "GalaxyS8",
	DisplayNamePrefixed: "r/GalaxyS8",
	Title:               "Samsung Galaxy S8",
	Description:         "### Rules\n\n* Posts and comments must be relevant to the Galaxy S8.\n* Do not post any referral codes.\n* No trolling.\n* No buying/selling/trading.\n* Do not editorialize submission titles.\n* No spamming or blog-spam.\n* Any photos/videos taken with the S8 should be posted in the weekly photography thread.\n\n### Link flair must be used\n\n* News\n* Rumor\n* Discussion\n* Help\n* Tricks\n* Creative\n* Other\n\nFlair can also be added by putting it in brackets before the post title, for example:\n&gt; [Help] I need help\n\n### Related Subreddits\n\n* [Samsung](https://www.reddit.com/r/Samsung)\n* [Galaxy Photography](https://www.reddit.com/r/galaxyphotography)\n* [Amoled Backgrounds](https://www.reddit.com/r/Amoledbackgrounds)\n\n### Discord Server\n\n* [Click Here to Join](https://discord.gg/4uxusu8)",
	SubredditType:       "public",
	Subscribers:         52357,

	// New
	AcceptFollowers:             false,
	AccountsActive:              0,
	AccountsActiveIsFuzzed:      false,
	ActiveUserCount:             0,
	AdvertiserCategory:          "",
	AllOriginalContent:          false,
	AllowChatPostCreation:       false,
	AllowDiscovery:              false,
	AllowGalleries:              false,
	AllowImages:                 false,
	AllowPolls:                  false,
	AllowPredictionContributors: false,
	AllowPredictions:            false,
	AllowPredictionsTournament:  false,
	AllowTalks:                  false,
	AllowVideogifs:              false,
	AllowVideos:                 false,
	AllowedMediaInComments:      false,
	BannerBackgroundColor:       "",
	BannerBackgroundImage:       "",
	BannerImg:                   "",
	BannerSize:                  []float64(nil),
	CanAssignLinkFlair:          false,
	CanAssignUserFlair:          false,
	CollapseDeletedComments:     false,
	CommentContributionSettings: struct{}{},
	CommentScoreHideMins:        0,
	CommunityIcon:               "",
	CommunityReviewed:           false,
	CreatedUtc:                  &Timestamp{time.Date(2013, time.March, 1, 4, 4, 18, 0, time.UTC)},
	DescriptionHTML:             "",
	DisableContributorRequests:  false,
	EmojisCustomSize:            interface{}(nil),
	EmojisEnabled:               false,
	FreeFormReports:             true,
	HasMenuWidget:               false,
	HeaderImg:                   "https://b.thumbs.redditmedia.com/AfySt3BMPjuq79LOh84X4uomahu0JE8DLaJZMenG-5I.png",
	HeaderSize: []int{1,
		1},
	HeaderTitle: "",
	HideAds:     false,
	IconImg:     "https://b.thumbs.redditmedia.com/4hg41g2_X1R5S_HTUscWCK_7iAo6SPdag_oOlSx7WAM.png",
	IconSize: []float64{256,
		256},
	ID:                               "",
	IsChatPostFeatureEnabled:         true,
	IsCrosspostableSubreddit:         false,
	IsEnrolledInNewModmail:           false,
	KeyColor:                         "",
	Lang:                             "",
	LinkFlairEnabled:                 true,
	LinkFlairPosition:                "left",
	MobileBannerImage:                "",
	NotificationLevel:                interface{}(nil),
	OriginalContentTagEnabled:        false,
	Over18:                           false,
	PredictionLeaderboardEntryType:   "",
	PrimaryColor:                     "#373c3f",
	PublicDescription:                "The only place for news, discussion, photos, and everything else Samsung Galaxy S8.",
	PublicDescriptionHTML:            "",
	PublicTraffic:                    false,
	Quarantine:                       false,
	RestrictCommenting:               false,
	RestrictPosting:                  true,
	ShouldArchivePosts:               false,
	ShouldShowMediaInCommentsSetting: false,
	ShowMedia:                        true,
	ShowMediaPreview:                 false,
	SpoilersEnabled:                  false,
	SubmissionType:                   "",
	SubmitLinkLabel:                  "",
	SubmitText:                       "",
	SubmitTextHTML:                   "",
	SubmitTextLabel:                  "",
	SuggestedCommentSort:             "",
	UserCanFlairInSr:                 interface{}(nil),
	UserFlairBackgroundColor:         interface{}(nil),
	UserFlairCSSClass:                interface{}(nil),
	UserFlairEnabledInSr:             false,
	UserFlairPosition:                "",
	UserFlairRichtext:                []interface{}(nil),
	UserFlairTemplateID:              interface{}(nil),
	UserFlairText:                    interface{}(nil),
	UserFlairTextColor:               interface{}(nil),
	UserFlairType:                    "",
	UserHasFavorited:                 false,
	UserIsBanned:                     false,
	UserIsContributor:                false,
	UserIsModerator:                  false,
	UserIsMuted:                      false,
	UserIsSubscriber:                 false,
	UserSrFlairEnabled:               false,
	UserSrThemeEnabled:               false,
	WhitelistStatus:                  "",
	WikiEnabled:                      false,
	Wls:                              0,
}

var expectedRelationships3 = []*Relationship{
	{
		ID:      "rel_id1",
		Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 2, 0, time.UTC)},
		User:    "testuser1",
		UserID:  "t2_user1",
	},
	{
		ID:      "rel_id2",
		Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 0, 0, time.UTC)},
		User:    "testuser2",
		UserID:  "t2_user2",
	},
}

var expectedBans = []*Ban{
	{
		Relationship: &Relationship{
			ID:      "rb_123",
			Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 2, 0, time.UTC)},

			User:   "testuser1",
			UserID: "t2_user1",
		},

		DaysLeft: Int(43),
		Note:     "Spam",
	},
	{
		Relationship: &Relationship{
			ID:      "rb_456",
			Created: &Timestamp{time.Date(2020, 8, 11, 2, 35, 0, 0, time.UTC)},

			User:   "testuser2",
			UserID: "t2_user2",
		},

		DaysLeft: nil,
		Note:     "Spam",
	},
}

var expectedModerators = []*Moderator{
	{
		Relationship: &Relationship{
			ID:      "rb_tmatb9",
			User:    "testuser1",
			UserID:  "t2_test1",
			Created: &Timestamp{time.Date(2013, 7, 29, 20, 44, 27, 0, time.UTC)},
		},
		Permissions: []string{"all"},
	},
	{
		Relationship: &Relationship{
			ID:      "rb_5c9s4d",
			User:    "testuser2",
			UserID:  "t2_test2",
			Created: &Timestamp{time.Date(2014, 3, 1, 18, 13, 53, 0, time.UTC)},
		},
		Permissions: []string{"all"},
	},
}

var expectedRules = []*SubredditRule{
	{
		Kind:            "link",
		Name:            "Read the Rules Before Posting",
		ViolationReason: "Read the Rules Before Posting",
		Description:     "https://www.reddit.com/r/Fitness/wiki/rules",
		Priority:        0,
		Created:         &Timestamp{time.Date(2019, 5, 22, 5, 32, 58, 0, time.UTC)},
	},
	{
		Kind:            "link",
		Name:            "Read the Wiki Before Posting",
		ViolationReason: "Read the Wiki Before Posting",
		Description:     "https://thefitness.wiki",
		Priority:        1,
		Created:         &Timestamp{time.Date(2019, 11, 9, 7, 56, 33, 0, time.UTC)},
	},
}

var expectedDayTraffic = []*SubredditTrafficStats{
	{&Timestamp{time.Date(2020, 9, 13, 0, 0, 0, 0, time.UTC)}, 0, 0, 0},
	{&Timestamp{time.Date(2020, 9, 12, 0, 0, 0, 0, time.UTC)}, 1, 12, 0},
	{&Timestamp{time.Date(2020, 9, 11, 0, 0, 0, 0, time.UTC)}, 5, 85, 0},
	{&Timestamp{time.Date(2020, 9, 10, 0, 0, 0, 0, time.UTC)}, 4, 20, 0},
	{&Timestamp{time.Date(2020, 9, 9, 0, 0, 0, 0, time.UTC)}, 2, 64, 0},
	{&Timestamp{time.Date(2020, 9, 8, 0, 0, 0, 0, time.UTC)}, 2, 95, 0},
	{&Timestamp{time.Date(2020, 9, 7, 0, 0, 0, 0, time.UTC)}, 3, 41, 0},
}

var expectedHourTraffic = []*SubredditTrafficStats{
	{&Timestamp{time.Date(2020, 9, 12, 20, 0, 0, 0, time.UTC)}, 1, 12, 0},
	{&Timestamp{time.Date(2020, 9, 11, 3, 0, 0, 0, time.UTC)}, 4, 57, 0},
	{&Timestamp{time.Date(2020, 9, 11, 2, 0, 0, 0, time.UTC)}, 4, 28, 0},
}

var expectedMonthTraffic = []*SubredditTrafficStats{
	{&Timestamp{time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC)}, 7, 481, 0},
	{&Timestamp{time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC)}, 5, 346, 0},
	{&Timestamp{time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)}, 4, 264, 0},
}

var expectedStyleSheet = &SubredditStyleSheet{
	SubredditID: "t5_2rc7j",
	Images: []*SubredditImage{
		{
			Name: "gopher",
			Link: "url(%%gopher%%)",
			URL:  "http://b.thumbs.redditmedia.com/q5Wb6hTPm2Bd6Of9_xMrTu4n5qgAljJNqtnbE3Tging.png",
		},
	},
	StyleSheet: `.flair-gopher {
    background: url(%%gopher%%) no-repeat;
    border: 0;
    padding: 0;
    width: 16px;
    height: 16px;
}`,
}

var expectedSubredditSettings = &SubredditSettings{
	ID: "t5_test",

	Type: String("private"),

	Language: String("en"),

	Title:                 String("hello!"),
	Description:           String("description"),
	Sidebar:               String("sidebar"),
	SubmissionText:        String(""),
	WelcomeMessage:        String(""),
	WelcomeMessageEnabled: Bool(false),

	AllowCrossposts:            Bool(false),
	AllowChatPosts:             Bool(true),
	AllowPollPosts:             Bool(false),
	AllowFreeFormReports:       Bool(true),
	AllowOriginalContent:       Bool(false),
	AllowImages:                Bool(true),
	AllowMultipleImagesPerPost: Bool(true),

	ExcludeSitewideBannedUsersContent: Bool(false),

	CrowdControlChalLevel: Int(2),

	AllOriginalContent: Bool(false),

	SuggestedCommentSort: nil,

	SubmitLinkPostLabel: String("submit a link!"),
	SubmitTextPostLabel: String("submit a post!"),

	PostType: String("any"),

	SpamFilterStrengthLinkPosts: String("low"),
	SpamFilterStrengthTextPosts: String("low"),
	SpamFilterStrengthComments:  String("low"),

	ShowContentThumbnails:              Bool(false),
	ExpandMediaPreviewsOnCommentsPages: Bool(true),

	CollapseDeletedComments:    Bool(false),
	MinutesToHideCommentScores: Int(0),

	SpoilersEnabled: Bool(true),

	HeaderMouseoverText: String("hello!"),

	MobileColour: String(""),

	HideAds: Bool(false),
	NSFW:    Bool(false),

	AllowDiscoveryInHighTrafficFeeds: Bool(true),
	AllowDiscoveryByIndividualUsers:  Bool(true),

	WikiMode:              String("modonly"),
	WikiMinimumAccountAge: Int(0),
	WikiMinimumKarma:      Int(0),
}

var expectedSubredditPostRequirements = &SubredditPostRequirements{
	Guidelines:              "test",
	GuidelinesDisplayPolicy: "",

	TitleMinLength: 50,
	TitleMaxLength: 200,

	BodyMinLength: 50,
	BodyMaxLength: 2000,

	TitleBlacklistedStrings: []string{"no"},
	BodyBlacklistedStrings:  []string{"no"},

	TitleRequiredStrings: []string{"yes"},
	BodyRequiredStrings:  []string{"yes"},

	DomainBlacklist: []string{"example.com"},
	DomainWhitelist: []string{},

	BodyRestrictionPolicy: "none",
	LinkRestrictionPolicy: "none",

	GalleryMinItems:            2,
	GalleryMaxItems:            20,
	GalleryCaptionsRequirement: "none",
	GalleryURLsRequirement:     "none",

	LinkRepostAge: 2,
	FlairRequired: false,

	TitleRegexes: []string{},
	BodyRegexes:  []string{},
}

func TestSubredditService_HotPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/hot", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.HotPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_NewPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/new", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.NewPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_RisingPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/rising", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.RisingPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_ControversialPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/controversial", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.ControversialPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_TopPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/top", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.TopPosts(ctx, "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)
	require.Equal(t, "t3_hyhquk", resp.After)
}

func TestSubredditService_Get(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/about.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/golang/about", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	_, _, err = client.Subreddit.Get(ctx, "")
	require.EqualError(t, err, "name: cannot be empty")

	subreddit, _, err := client.Subreddit.Get(ctx, "golang")
	require.NoError(t, err)
	require.Equal(t, expectedSubreddit, subreddit)
}

func TestSubredditService_Popular(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/popular", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Popular(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_New(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/new", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.New(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Gold(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/gold", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Gold(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Default(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/default", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Default(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Subscribed(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/mine/subscriber", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Subscribed(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Approved(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/mine/contributor", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Approved(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_Moderated(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/mine/moderator", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Moderated(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_GetSticky1(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "1", r.Form.Get("num"))

		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Subreddit.GetSticky1(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestSubredditService_GetSticky2(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/post/post.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/sticky", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "2", r.Form.Get("num"))

		fmt.Fprint(w, blob)
	})

	postAndComments, _, err := client.Subreddit.GetSticky2(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, expectedPostAndComments, postAndComments)
}

func TestSubredditService_Subscribe(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "sub")
		form.Set("sr_name", "test,golang,nba")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Subscribe(ctx, "test", "golang", "nba")
	require.NoError(t, err)
}

func TestSubredditService_SubscribeByID(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "sub")
		form.Set("sr", "t5_test1,t5_test2,t5_test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.SubscribeByID(ctx, "t5_test1", "t5_test2", "t5_test3")
	require.NoError(t, err)
}

func TestSubredditService_Unsubscribe(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "unsub")
		form.Set("sr_name", "test,golang,nba")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Unsubscribe(ctx, "test", "golang", "nba")
	require.NoError(t, err)
}

func TestSubredditService_UnsubscribeByID(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("action", "unsub")
		form.Set("sr", "t5_test1,t5_test2,t5_test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.UnsubscribeByID(ctx, "t5_test1", "t5_test2", "t5_test3")
	require.NoError(t, err)
}

func TestSubredditService_Favorite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/favorite", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("sr_name", "testsubreddit")
		form.Set("make_favorite", "true")
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Favorite(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_Unfavorite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/favorite", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("sr_name", "testsubreddit")
		form.Set("make_favorite", "false")
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Unfavorite(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_Search(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/list.json")
	require.NoError(t, err)

	mux.HandleFunc("/subreddits/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "golang")
		form.Set("limit", "10")
		form.Set("sort", "activity")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	subreddits, resp, err := client.Subreddit.Search(ctx, "golang", &ListSubredditOptions{
		ListOptions: ListOptions{
			Limit: 10,
		},
		Sort: "activity",
	})
	require.NoError(t, err)
	require.Equal(t, expectedSubreddits, subreddits)
	require.Equal(t, "t5_2qh0u", resp.After)
}

func TestSubredditService_SearchNames(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-names.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/search_reddit_names", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("query", "golang")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	names, _, err := client.Subreddit.SearchNames(ctx, "golang")
	require.NoError(t, err)
	require.Equal(t, expectedSubredditNames, names)
}

func TestSubredditService_SearchPosts(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/all/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.SearchPosts(ctx, "test", "", nil)
	require.NoError(t, err)
	require.Equal(t, expectedSearchPosts, posts)
	require.Equal(t, "t3_hmwhd7", resp.After)
}

func TestSubredditService_SearchPosts_InSubreddit(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "test")
		form.Set("restrict_sr", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.SearchPosts(ctx, "test", "test", nil)
	require.NoError(t, err)
	require.Equal(t, expectedSearchPosts, posts)
	require.Equal(t, "t3_hmwhd7", resp.After)
}

func TestSubredditService_SearchPosts_InSubreddits(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/search-posts.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test+golang+nba/search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("q", "test")
		form.Set("restrict_sr", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Subreddit.SearchPosts(ctx, "test", "test+golang+nba", nil)
	require.NoError(t, err)
	require.Equal(t, expectedSearchPosts, posts)
	require.Equal(t, "t3_hmwhd7", resp.After)
}

func TestSubredditService_Random(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/random.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/random", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "true", r.Form.Get("sr_detail"))
		require.Equal(t, "1", r.Form.Get("limit"))

		fmt.Fprint(w, blob)
	})

	subreddit, _, err := client.Subreddit.Random(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedRandomSubreddit, subreddit)
}

func TestSubredditService_RandomNSFW(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/random.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/randnsfw", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, "true", r.Form.Get("sr_detail"))
		require.Equal(t, "1", r.Form.Get("limit"))

		fmt.Fprint(w, blob)
	})

	subreddit, _, err := client.Subreddit.RandomNSFW(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedRandomSubreddit, subreddit)
}

func TestSubredditService_SubmissionText(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/test/api/submit_text", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `{
			"submit_text": "this is a test",
			"submit_text_html": ""
		}`)
	})

	text, _, err := client.Subreddit.SubmissionText(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, "this is a test", text)
}

func TestSubredditService_Banned(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/banned-users.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/banned", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("after", "testafter")
		form.Set("limit", "10")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	bans, _, err := client.Subreddit.Banned(ctx, "test", &ListOptions{After: "testafter", Limit: 10})
	require.NoError(t, err)
	require.Equal(t, expectedBans, bans)
}

func TestSubredditService_Muted(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/relationships.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/muted", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("before", "testbefore")
		form.Set("limit", "50")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	mutes, _, err := client.Subreddit.Muted(ctx, "test", &ListOptions{Before: "testbefore", Limit: 50})
	require.NoError(t, err)
	require.Equal(t, expectedRelationships3, mutes)
}

func TestSubredditService_WikiBanned(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/banned-users.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/wikibanned", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("after", "testafter")
		form.Set("limit", "15")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	bans, _, err := client.Subreddit.WikiBanned(ctx, "test", &ListOptions{After: "testafter", Limit: 15})
	require.NoError(t, err)
	require.Equal(t, expectedBans, bans)
}

func TestSubredditService_Contributors(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/relationships.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/contributors", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "5")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	contributors, _, err := client.Subreddit.Contributors(ctx, "test", &ListOptions{Limit: 5})
	require.NoError(t, err)
	require.Equal(t, expectedRelationships3, contributors)
}

func TestSubredditService_WikiContributors(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/relationships.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/wikicontributors", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("limit", "99")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	contributors, _, err := client.Subreddit.WikiContributors(ctx, "test", &ListOptions{Limit: 99})
	require.NoError(t, err)
	require.Equal(t, expectedRelationships3, contributors)
}

func TestSubredditService_Moderators(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/moderators.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/test/about/moderators", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	moderators, _, err := client.Subreddit.Moderators(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, expectedModerators, moderators)
}

func TestSubredditService_Rules(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/rules.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/rules", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	rules, _, err := client.Subreddit.Rules(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedRules, rules)
}

func TestSubredditService_CreateRule(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/add_subreddit_rule", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("kind", "all")
		form.Set("short_name", "testname")
		form.Set("violation_reason", "testreason")
		form.Set("description", "testdescription")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{
		Kind:            "all",
		Name:            "testname",
		ViolationReason: "testreason",
		Description:     "testdescription",
	})
	require.NoError(t, err)
}

func TestSubredditService_CreateRule_Error(t *testing.T) {
	client, _ := setup(t)

	_, err := client.Subreddit.CreateRule(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "*SubredditRuleCreateRequest: cannot be nil")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{Kind: "invalid"})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).Kind: must be one of: comment, link, all")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{Kind: "all", Name: ""})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).Name: must be between 1-100 characters")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{
		Kind:            "all",
		Name:            "testname",
		ViolationReason: strings.Repeat("x", 101),
	})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).ViolationReason: cannot be longer than 100 characters")

	_, err = client.Subreddit.CreateRule(ctx, "testsubreddit", &SubredditRuleCreateRequest{
		Kind:        "all",
		Name:        "testname",
		Description: strings.Repeat("x", 501),
	})
	require.EqualError(t, err, "(*SubredditRuleCreateRequest).Description: cannot be longer than 500 characters")
}

func TestSubredditService_Traffic(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/traffic.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/traffic", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	dayTraffic, hourTraffic, monthTraffic, _, err := client.Subreddit.Traffic(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedDayTraffic, dayTraffic)
	require.Equal(t, expectedHourTraffic, hourTraffic)
	require.Equal(t, expectedMonthTraffic, monthTraffic)
}

func TestSubredditService_StyleSheet(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/stylesheet.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/stylesheet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	styleSheet, _, err := client.Subreddit.StyleSheet(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedStyleSheet, styleSheet)
}

func TestSubredditService_StyleSheetRaw(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/stylesheet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, "* { box-sizing: border-box; }")
	})

	styleSheet, _, err := client.Subreddit.StyleSheetRaw(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, "* { box-sizing: border-box; }", styleSheet)
}

func TestSubredditService_UpdateStyleSheet(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/subreddit_stylesheet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("op", "save")
		form.Set("stylesheet_contents", "* { box-sizing: border-box; }")
		form.Set("reason", "testreason")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.UpdateStyleSheet(ctx, "testsubreddit", "* { box-sizing: border-box; }", "testreason")
	require.NoError(t, err)
}

func TestSubredditService_RemoveImage(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("img_name", "testimage")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveImage(ctx, "testsubreddit", "testimage")
	require.NoError(t, err)
}

func TestSubredditService_RemoveHeader(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_header", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveHeader(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_RemoveMobileHeader(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_banner", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveMobileHeader(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_RemoveMobileIcon(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/delete_sr_icon", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.RemoveMobileIcon(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestSubredditService_UploadImage(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.png")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "img")
		form.Set("name", "testname")
		form.Set("img_type", "png")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.png"
		}`)
	})

	link, _, err := client.Subreddit.UploadImage(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.png", link)
}

func TestSubredditService_UploadHeader(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.png")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "header")
		form.Set("name", "testname")
		form.Set("img_type", "png")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.png"
		}`)
	})

	link, _, err := client.Subreddit.UploadHeader(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.png", link)
}

func TestSubredditService_UploadMobileHeader(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.png")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "banner")
		form.Set("name", "testname")
		form.Set("img_type", "png")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.png"
		}`)
	})

	link, _, err := client.Subreddit.UploadMobileHeader(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.png", link)
}

func TestSubredditService_UploadMobileIcon(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.jpg")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	_, err = imageFile.WriteString("this is a test")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		require.NoError(t, err)

		rdr, err := file.Open()
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		require.NoError(t, err)
		require.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("upload_type", "icon")
		form.Set("name", "testname")
		form.Set("img_type", "jpg")

		err = r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, `{
			"img_src": "https://example.com/test.jpg"
		}`)
	})

	link, _, err := client.Subreddit.UploadMobileIcon(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.NoError(t, err)
	require.Equal(t, "https://example.com/test.jpg", link)
}

func TestSubredditService_UploadImage_Error(t *testing.T) {
	client, mux := setup(t)

	imageFile, err := ioutil.TempFile("/tmp", "emoji*.jpg")
	require.NoError(t, err)
	defer func() {
		imageFile.Close()
		os.Remove(imageFile.Name())
	}()

	mux.HandleFunc("/r/testsubreddit/api/upload_sr_img", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		fmt.Fprint(w, `{
			"errors_values": [
				"error one",
				"error two"
			]
		}`)
	})

	_, _, err = client.Subreddit.UploadImage(ctx, "testsubreddit", "does-not-exist.jpg", "testname")
	require.EqualError(t, err, "open does-not-exist.jpg: no such file or directory")

	_, _, err = client.Subreddit.UploadImage(ctx, "testsubreddit", imageFile.Name(), "testname")
	require.EqualError(t, err, "could not upload image: error one; error two")
}

func TestSubredditService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/site_admin", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testsubreddit")
		form.Set("type", "private")
		form.Set("lang", "en")
		form.Set("title", "hello!")
		form.Set("public_description", "description")
		form.Set("description", "sidebar")
		form.Set("submit_text", "")
		form.Set("welcome_message_text", "")
		form.Set("welcome_message_enabled", "false")
		form.Set("allow_post_crossposts", "false")
		form.Set("allow_chat_post_creation", "true")
		form.Set("allow_polls", "false")
		form.Set("free_form_reports", "true")
		form.Set("original_content_tag_enabled", "false")
		form.Set("allow_images", "true")
		form.Set("allow_galleries", "true")
		form.Set("exclude_banned_modqueue", "false")
		form.Set("crowd_control_chat_level", "2")
		form.Set("all_original_content", "false")
		form.Set("submit_link_label", "submit a link!")
		form.Set("submit_text_label", "submit a post!")
		form.Set("link_type", "any")
		form.Set("spam_links", "low")
		form.Set("spam_selfposts", "low")
		form.Set("spam_comments", "low")
		form.Set("show_media", "false")
		form.Set("show_media_preview", "true")
		form.Set("collapse_deleted_comments", "false")
		form.Set("comment_score_hide_mins", "0")
		form.Set("spoilers_enabled", "true")
		form.Set("header-title", "hello!")
		form.Set("key_color", "")
		form.Set("hide_ads", "false")
		form.Set("over_18", "false")
		form.Set("allow_top", "true")
		form.Set("allow_discovery", "true")
		form.Set("wikimode", "modonly")
		form.Set("wiki_edit_age", "0")
		form.Set("wiki_edit_karma", "0")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Create(ctx, "testsubreddit", nil)
	require.EqualError(t, err, "*SubredditSettings: cannot be nil")

	_, err = client.Subreddit.Create(ctx, "testsubreddit", expectedSubredditSettings)
	require.NoError(t, err)
}

func TestSubredditService_Edit(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/site_admin", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("sr", "t5_test")
		form.Set("type", "private")
		form.Set("lang", "en")
		form.Set("title", "hello!")
		form.Set("public_description", "description")
		form.Set("description", "sidebar")
		form.Set("submit_text", "")
		form.Set("welcome_message_text", "")
		form.Set("welcome_message_enabled", "false")
		form.Set("allow_post_crossposts", "false")
		form.Set("allow_chat_post_creation", "true")
		form.Set("allow_polls", "false")
		form.Set("free_form_reports", "true")
		form.Set("original_content_tag_enabled", "false")
		form.Set("allow_images", "true")
		form.Set("allow_galleries", "true")
		form.Set("exclude_banned_modqueue", "false")
		form.Set("crowd_control_chat_level", "2")
		form.Set("all_original_content", "false")
		form.Set("submit_link_label", "submit a link!")
		form.Set("submit_text_label", "submit a post!")
		form.Set("link_type", "any")
		form.Set("spam_links", "low")
		form.Set("spam_selfposts", "low")
		form.Set("spam_comments", "low")
		form.Set("show_media", "false")
		form.Set("show_media_preview", "true")
		form.Set("collapse_deleted_comments", "false")
		form.Set("comment_score_hide_mins", "0")
		form.Set("spoilers_enabled", "true")
		form.Set("header-title", "hello!")
		form.Set("key_color", "")
		form.Set("hide_ads", "false")
		form.Set("over_18", "false")
		form.Set("allow_top", "true")
		form.Set("allow_discovery", "true")
		form.Set("wikimode", "modonly")
		form.Set("wiki_edit_age", "0")
		form.Set("wiki_edit_karma", "0")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Subreddit.Edit(ctx, "t5_test", nil)
	require.EqualError(t, err, "*SubredditSettings: cannot be nil")

	_, err = client.Subreddit.Edit(ctx, "t5_test", expectedSubredditSettings)
	require.NoError(t, err)
}

func TestSubredditService_GetSettings(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/settings.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/edit", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	subredditSettings, _, err := client.Subreddit.GetSettings(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedSubredditSettings, subredditSettings)
}

func TestSubredditService_PostRequirements(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/subreddit/post-requirements.json")
	require.NoError(t, err)

	mux.HandleFunc("/api/v1/testsubreddit/post_requirements", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	postRequirements, _, err := client.Subreddit.PostRequirements(ctx, "testsubreddit")
	require.NoError(t, err)
	require.Equal(t, expectedSubredditPostRequirements, postRequirements)
}
