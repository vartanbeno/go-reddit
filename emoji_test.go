package reddit

import (
	"fmt"
	"net/http"
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
