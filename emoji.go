package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// EmojiService handles communication with the emoji
// related methods of the Reddit API.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_collections
type EmojiService struct {
	client *Client
}

// Emoji is a graphic element you can include in a post flair or user flair.
type Emoji struct {
	Name             string `json:"name,omitempty"`
	URL              string `json:"url,omitempty"`
	UserFlairAllowed bool   `json:"user_flair_allowed,omitempty"`
	PostFlairAllowed bool   `json:"post_flair_allowed,omitempty"`
	ModFlairOnly     bool   `json:"mod_flair_only,omitempty"`
	// ID of the user who created this emoji.
	CreatedBy string `json:"created_by,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// func (e *Emoji) UnmarshalJSON(data []byte) (err error) {
// 	fmt.Println("===", string(data))
// 	return nil
// }

type emojis []*Emoji

func (e *emojis) UnmarshalJSON(data []byte) (err error) {
	emojiMap := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &emojiMap)
	if err != nil {
		return
	}

	for emojiName, emojiValue := range emojiMap {
		emoji := new(Emoji)
		err = json.Unmarshal(emojiValue, emoji)
		if err != nil {
			return
		}
		emoji.Name = emojiName
		*e = append(*e, emoji)
	}

	return
}

// Get returns the default set of Reddit emojis, and those of the subreddit, respectively.
func (s *EmojiService) Get(ctx context.Context, subreddit string) ([]*Emoji, []*Emoji, *Response, error) {
	path := fmt.Sprintf("api/v1/%s/emojis/all", subreddit)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	root := make(map[string]emojis)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, nil, resp, err
	}

	/*
		The response to this request is something like:
		{
			"snoomojis": { ... },
			"t5_subredditId": { ... }
		}
	*/
	defaultEmojis := root["snoomojis"]
	var subredditEmojis []*Emoji

	for k := range root {
		if strings.HasPrefix(k, kindSubreddit) {
			subredditEmojis = root[k]
			break
		}
	}

	return defaultEmojis, subredditEmojis, resp, nil
}

// Delete deletes the emoji from the subreddit.
func (s *EmojiService) Delete(ctx context.Context, subreddit string, emoji string) (*Response, error) {
	path := fmt.Sprintf("api/v1/%s/emoji/%s", subreddit, emoji)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

//todo: fav subreddits

// SetSize sets the custom emoji size in the subreddit.
// Both height and width must be between 1 and 40 (inclusive).
func (s *EmojiService) SetSize(ctx context.Context, subreddit string, height, width int) (*Response, error) {
	path := fmt.Sprintf("api/v1/%s/emoji_custom_size", subreddit)

	form := url.Values{}
	form.Set("height", fmt.Sprint(height))
	form.Set("width", fmt.Sprint(width))

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DisableCustomSize disables the custom emoji size in the subreddit.
func (s *EmojiService) DisableCustomSize(ctx context.Context, subreddit string) (*Response, error) {
	path := fmt.Sprintf("api/v1/%s/emoji_custom_size", subreddit)

	req, err := s.client.NewRequestWithForm(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
