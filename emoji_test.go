package reddit

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

	blob, err := readFileContents("testdata/emoji/emojis.json")
	assert.NoError(t, err)

	mux.HandleFunc("/api/v1/test/emojis/all", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	defaultEmojis, subredditEmojis, _, err := client.Emoji.Get(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, defaultEmojis, 2)
	assert.Contains(t, expectedDefaultEmojis, defaultEmojis[0])
	assert.Contains(t, expectedDefaultEmojis, defaultEmojis[1])
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

func TestEmojiService_Upload(t *testing.T) {
	setup()
	defer teardown()

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)

	uploadURL := u.Host + "/api/emoji_upload"

	blob, err := readFileContents("testdata/emoji/lease.json")
	assert.NoError(t, err)
	blob = fmt.Sprintf(blob, uploadURL)

	emojiFile, err := ioutil.TempFile("/tmp", "emoji*.png")
	assert.NoError(t, err)
	defer func() {
		emojiFile.Close()
		os.Remove(emojiFile.Name())
	}()

	_, err = emojiFile.WriteString("this is a test")
	assert.NoError(t, err)

	mux.HandleFunc("/api/v1/testsubreddit/emoji_asset_upload_s3.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("filepath", emojiFile.Name())
		form.Set("mimetype", "image/png")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	mux.HandleFunc("/api/emoji_upload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		_, file, err := r.FormFile("file")
		assert.NoError(t, err)

		rdr, err := file.Open()
		assert.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rdr)
		assert.NoError(t, err)
		assert.Equal(t, "this is a test", buf.String())

		form := url.Values{}
		form.Set("key", "t5_2uquw1/t2_164ab8/a94a8f45ccb199a61c4c0873d391e98c982fabd3")
		form.Set("test name", "test value")

		// for some reason this has to come after the FormFile call
		err = r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	mux.HandleFunc("/api/v1/testsubreddit/emoji.json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("name", "testemoji")
		form.Set("user_flair_allowed", "false")
		form.Set("post_flair_allowed", "true")
		form.Set("mod_flair_only", "true")
		form.Set("s3_key", "t5_2uquw1/t2_164ab8/a94a8f45ccb199a61c4c0873d391e98c982fabd3")

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, form, r.Form)
	})

	_, err = client.Emoji.Upload(ctx, "testsubreddit", nil, emojiFile.Name())
	assert.EqualError(t, err, "createRequest: cannot be nil")

	_, err = client.Emoji.Upload(ctx, "testsubreddit", &EmojiCreateOrUpdateRequest{Name: ""}, emojiFile.Name())
	assert.EqualError(t, err, "name: cannot be empty")

	_, err = client.Emoji.Upload(ctx, "testsubreddit", &EmojiCreateOrUpdateRequest{
		Name:             "testemoji",
		UserFlairAllowed: Bool(false),
		PostFlairAllowed: Bool(true),
		ModFlairOnly:     Bool(true),
	}, emojiFile.Name())
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
