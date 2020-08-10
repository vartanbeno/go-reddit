package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedDefaultEmojis = []*Emoji{
	{
		Name:             "cake",
		URL:              "https://emoji.redditmedia.com/46kel8lf1guz_t5_3nqvj/cake",
		UserFlairAllowed: true,
		PostFlairAllowed: true,
		ModFlairOnly:     false,
		CreatedBy:        "t2_6zfp6ii",
	},
	{
		Name:             "cat_blep",
		URL:              "https://emoji.redditmedia.com/p9sxc1zh1guz_t5_3nqvj/cat_blep",
		UserFlairAllowed: true,
		PostFlairAllowed: true,
		ModFlairOnly:     false,
		CreatedBy:        "t2_6zfp6ii",
	},
}

var expectedSubredditEmojis = []*Emoji{
	{
		Name:             "TestEmoji",
		URL:              "https://emoji.redditmedia.com/fxe5a674hpf51_t5_2uquw1/TestEmoji",
		UserFlairAllowed: true,
		PostFlairAllowed: true,
		ModFlairOnly:     false,
		CreatedBy:        "t2_164ab8",
	},
}

func TestEmojiService_Get(t *testing.T) {
	setup()
	defer teardown()

	blob, err := readFileContents("./testdata/emoji/emojis.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/v1/test/emojis/all", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	defaultEmojis, subredditEmojis, _, err := client.Emoji.Get(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedDefaultEmojis, defaultEmojis)
	assert.Equal(t, expectedSubredditEmojis, subredditEmojis)
}

func TestEmojiService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/testsubreddit/emoji/testemoji", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
	})

	_, err := client.Emoji.Delete(ctx, "testsubreddit", "testemoji")
	assert.NoError(t, err)
}

func TestEmojiService_SetSize(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/testsubreddit/emoji_custom_size", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("height", "20")
		form.Set("width", "20")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Emoji.SetSize(ctx, "testsubreddit", 20, 20)
	assert.NoError(t, err)
}

func TestEmojiService_DisableCustomSize(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/testsubreddit/emoji_custom_size", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Emoji.DisableCustomSize(ctx, "testsubreddit")
	assert.NoError(t, err)
}

func TestEmojiService_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/testsubreddit/emoji_permissions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("name", "testemoji")
		form.Set("post_flair_allowed", "false")
		form.Set("mod_flair_only", "true")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err := client.Emoji.Update(ctx, "testsubreddit", nil)
	assert.EqualError(t, err, "updateRequest: cannot be nil")

	_, err = client.Emoji.Update(ctx, "testsubreddit", &EmojiCreateOrUpdateRequest{Name: ""})
	assert.EqualError(t, err, "name: cannot be empty")

	_, err = client.Emoji.Update(ctx, "testsubreddit", &EmojiCreateOrUpdateRequest{
		Name:             "testemoji",
		PostFlairAllowed: Bool(false),
		ModFlairOnly:     Bool(true),
	})
	assert.NoError(t, err)
}
