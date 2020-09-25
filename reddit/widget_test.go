package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedWidgets = []Widget{
	&TextAreaWidget{
		widget: widget{
			ID:   "widget_15p7borvnnw5a",
			Kind: "textarea",
			Style: &WidgetStyle{
				HeaderColor:     "#373c3f",
				BackgroundColor: "#cc5289",
			},
		},
		Name: "test title",
		Text: "test text",
	},

	&ButtonWidget{
		widget: widget{
			ID:    "widget_15paxrbiodp8v",
			Kind:  "button",
			Style: &WidgetStyle{},
		},
		Name:        "test text",
		Description: "test description",
		Buttons: []*WidgetButton{
			{
				Text:        "test text",
				URL:         "https://example.com",
				TextColor:   "#ff66ac",
				FillColor:   "#014980",
				StrokeColor: "#73ad34",
				HoverState: &WidgetButtonHoverState{
					Text:        "test text",
					TextColor:   "#000000",
					FillColor:   "#00a6a5",
					StrokeColor: "#000000",
				},
			},
		},
	},

	&ImageWidget{
		widget: widget{
			ID:    "widget_15p7o01nqr5tu",
			Kind:  "image",
			Style: &WidgetStyle{},
		},
		Name: "test title",
		Images: []*WidgetImageLink{
			{
				URL:     "https://www.redditstatic.com/image-processing.png",
				LinkURL: "https://example.com",
			},
		},
	},

	&CommunityListWidget{
		widget: widget{
			ID:   "widget_15p7qwb2kxc6j",
			Kind: "community-list",
			Style: &WidgetStyle{
				HeaderColor: "#ffb000",
			},
		},
		Name: "test title",
		Communities: []*WidgetCommunity{
			{
				Name:        "nba",
				Subscribers: 3571840,
				Subscribed:  true,
				NSFW:        false,
			},
			{
				Name:        "golang",
				Subscribers: 125961,
				Subscribed:  true,
				NSFW:        false,
			},
		},
	},

	&SubredditRulesWidget{
		widget: widget{
			ID:    "widget_rules-2uquw1",
			Kind:  "subreddit-rules",
			Style: &WidgetStyle{},
		},
		Name:    "Subreddit Rules",
		Display: "compact",
		Rules:   []string{"be nice"},
	},

	&CommunityDetailsWidget{
		widget: widget{
			ID:    "widget_id-card-2uquw1",
			Kind:  "id-card",
			Style: &WidgetStyle{},
		},
		Name:                 "Community Details",
		Description:          "Community Description",
		Subscribers:          2,
		CurrentlyViewing:     3,
		SubscribersText:      "subscriberz",
		CurrentlyViewingText: "viewerz",
	},

	&MenuWidget{
		widget: widget{
			ID:    "widget_15owrhqvgfhke",
			Kind:  "menu",
			Style: &WidgetStyle{},
		},
		ShowWiki: true,
		Links: []WidgetLink{
			&WidgetLinkSingle{
				Text: "link1",
				URL:  "https://example.com",
			},
			&WidgetLinkMultiple{
				Text: "test",
				URLs: []*WidgetLinkSingle{
					{
						Text: "link2",
						URL:  "https://example.com",
					},
					{
						Text: "link3",
						URL:  "https://example.com",
					},
				},
			},
		},
	},

	&ModeratorsWidget{
		widget: widget{
			ID:    "widget_moderators-2uquw1",
			Kind:  "moderators",
			Style: &WidgetStyle{},
		},
		Mods:  []string{"testuser"},
		Total: 1,
	},

	&CustomWidget{
		widget: widget{
			ID:    "widget_15osq4jms4tdo",
			Kind:  "custom",
			Style: &WidgetStyle{},
		},
		Name:          "custom image widget",
		Text:          "some image",
		StyleSheet:    "* {}",
		StyleSheetURL: "https://styles.redditmedia.com/t5_2uquw1/styles/customWidget-stylesheet-n2q86gjf04o51.css",
		Images: []*WidgetImage{
			{
				Name: "test",
				URL:  "https://www.redditstatic.com/image-processing.png",
			},
		},
	},
}

func TestWidgetService_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	blob, err := readFileContents("../testdata/widget/widgets.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/widgets", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("progressive_images", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	widgets, _, err := client.Widget.Get(ctx, "testsubreddit")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedWidgets, widgets)
}

func TestWidgetService_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/r/testsubreddit/api/widget/abc123", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
	})

	_, err := client.Widget.Delete(ctx, "testsubreddit", "abc123")
	require.NoError(t, err)
}
