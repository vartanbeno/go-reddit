package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// WidgetService handles communication with the widget
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_widgets
type WidgetService struct {
	client *Client
}

// Widget is a section of useful content on a subreddit.
// They can feature information such as rules, links, the origins of the subreddit, etc.
// Read about them here: https://mods.reddithelp.com/hc/en-us/articles/360010364372-Sidebar-Widgets
type Widget interface {
	// kind returns the widget kind.
	// having un unexported method on an exported interface means it cannot be implemented by a client.
	kind() string
}

const (
	widgetKindMenu             = "menu"
	widgetKindCommunityDetails = "id-card"
	widgetKindModerators       = "moderators"
	widgetKindSubredditRules   = "subreddit-rules"
	widgetKindCustom           = "custom"
)

// WidgetList is a list of widgets.
type WidgetList []Widget

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *WidgetList) UnmarshalJSON(data []byte) error {
	var widgetMap map[string]json.RawMessage
	err := json.Unmarshal(data, &widgetMap)
	if err != nil {
		return err
	}

	type widgetKind struct {
		Kind string `json:"kind"`
	}
	for _, w := range widgetMap {
		root := new(widgetKind)
		err = json.Unmarshal(w, root)
		if err != nil {
			return err
		}

		var widget Widget
		switch root.Kind {
		case widgetKindMenu:
			widget = new(MenuWidget)
		case widgetKindCommunityDetails:
			widget = new(CommunityDetailsWidget)
		case widgetKindModerators:
			widget = new(ModeratorsWidget)
		case widgetKindSubredditRules:
			widget = new(SubredditRulesWidget)
		case widgetKindCustom:
			widget = new(CustomWidget)
		default:
			continue
		}

		err = json.Unmarshal(w, widget)
		if err != nil {
			return err
		}

		*l = append(*l, widget)
	}

	return nil
}

// common widget fields
type widget struct {
	ID    string       `json:"id,omitempty"`
	Kind  string       `json:"kind,omitempty"`
	Style *WidgetStyle `json:"styles,omitempty"`
}

// MenuWidget displays tabs for your community's menu. These can be direct links or submenus that
// create a drop-down menu to multiple links.
type MenuWidget struct {
	widget

	ShowWiki bool           `json:"showWiki"`
	Links    WidgetLinkList `json:"data,omitempty"`
}

func (w *MenuWidget) kind() string {
	return widgetKindMenu
}

// CommunityDetailsWidget displays your subscriber count, users online, and community description,
// as defined in your subreddit settings. You can customize the displayed text for subscribers and
// users currently viewing the community.
type CommunityDetailsWidget struct {
	widget

	Name        string `json:"shortName,omitempty"`
	Description string `json:"description,omitempty"`

	Subscribers      int `json:"subscribersCount"`
	CurrentlyViewing int `json:"currentlyViewingCount"`

	SubscribersText      string `json:"subscribersText,omitempty"`
	CurrentlyViewingText string `json:"currentlyViewingText,omitempty"`
}

func (*CommunityDetailsWidget) kind() string {
	return widgetKindCommunityDetails
}

// ModeratorsWidget displays the list of moderators of the subreddit.
type ModeratorsWidget struct {
	widget

	Mods  []string `json:"mods"`
	Total int      `json:"totalMods"`
}

func (*ModeratorsWidget) kind() string {
	return widgetKindModerators
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (w *ModeratorsWidget) UnmarshalJSON(data []byte) error {
	root := new(struct {
		widget

		Mods []struct {
			Name string `json:"name"`
		} `json:"mods"`
		Total int `json:"totalMods"`
	})

	err := json.Unmarshal(data, root)
	if err != nil {
		return err
	}

	w.widget = root.widget
	w.Total = root.Total
	for _, mod := range root.Mods {
		w.Mods = append(w.Mods, mod.Name)
	}

	return nil
}

// SubredditRulesWidget displays your community rules.
type SubredditRulesWidget struct {
	widget

	Name string `json:"shortName,omitempty"`
	// One of: full (includes description), compact (rule is collapsed).
	Display string   `json:"display,omitempty"`
	Rules   []string `json:"rules,omitempty"`
}

func (*SubredditRulesWidget) kind() string {
	return widgetKindSubredditRules
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (w *SubredditRulesWidget) UnmarshalJSON(data []byte) error {
	root := new(struct {
		widget

		Name    string `json:"shortName"`
		Display string `json:"display"`
		Rules   []struct {
			Description string `json:"description"`
		} `json:"data"`
	})

	err := json.Unmarshal(data, root)
	if err != nil {
		return err
	}

	w.widget = root.widget
	w.Name = root.Name
	w.Display = root.Display
	for _, r := range root.Rules {
		w.Rules = append(w.Rules, r.Description)
	}

	return nil
}

// CustomWidget is a custom widget.
type CustomWidget struct {
	widget

	Name string `json:"shortName,omitempty"`
	Text string `json:"text,omitempty"`

	StyleSheet    string         `json:"css,omitempty"`
	StyleSheetURL string         `json:"stylesheetUrl,omitempty"`
	Images        []*WidgetImage `json:"imageData,omitempty"`
}

func (*CustomWidget) kind() string {
	return widgetKindCustom
}

// WidgetStyle contains style information for the widget.
type WidgetStyle struct {
	HeaderColor     string `json:"headerColor"`
	BackgroundColor string `json:"backgroundColor"`
}

// WidgetImage is an image in a widget.
type WidgetImage struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// WidgetLink is a link or a group of links that's part of a widget.
type WidgetLink interface {
	// single returns whether or not the widget holds just one single link.
	// having un unexported method on an exported interface means it cannot be implemented by a client.
	single() bool
}

// WidgetLinkSingle is a link that's part of a widget.
type WidgetLinkSingle struct {
	Text string `json:"text,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (l *WidgetLinkSingle) single() bool { return true }

// WidgetLinkMultiple is a dropdown of multiple links that's part of a widget.
type WidgetLinkMultiple struct {
	Text string              `json:"text,omitempty"`
	URLs []*WidgetLinkSingle `json:"children,omitempty"`
}

func (l *WidgetLinkMultiple) single() bool { return false }

// WidgetLinkList is a list of widgets links.
type WidgetLinkList []WidgetLink

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *WidgetLinkList) UnmarshalJSON(data []byte) error {
	var dataMap []json.RawMessage
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	for _, d := range dataMap {
		var widgetLinkDataMap map[string]json.RawMessage
		err = json.Unmarshal(d, &widgetLinkDataMap)
		if err != nil {
			return err
		}

		var wl WidgetLink
		if _, ok := widgetLinkDataMap["children"]; ok {
			wl = new(WidgetLinkMultiple)
		} else {
			wl = new(WidgetLinkSingle)
		}

		err = json.Unmarshal(d, wl)
		if err != nil {
			return err
		}

		*l = append(*l, wl)
	}

	return nil
}

// Get the subreddit's widgets.
func (s *WidgetService) Get(ctx context.Context, subreddit string) ([]Widget, *Response, error) {
	path := fmt.Sprintf("r/%s/api/widgets?progressive_images=true", subreddit)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(struct {
		Widgets WidgetList `json:"items"`
	})
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Widgets, resp, nil
}

// Delete a widget via its id.
func (s *WidgetService) Delete(ctx context.Context, subreddit, id string) (*Response, error) {
	path := fmt.Sprintf("r/%s/api/widget/%s", subreddit, id)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
