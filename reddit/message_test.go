package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedCommentMessages = []*Message{
	{
		ID:      "g1xi2m9",
		FullID:  "t1_g1xi2m9",
		Created: &Timestamp{time.Date(2020, 8, 18, 0, 24, 13, 0, time.UTC)},

		Subject:  "post reply",
		Text:     "u/testuser2 hello",
		ParentID: "t3_hs03f3",

		Author: "testuser1",
		To:     "testuser2",

		IsComment: true,
	},
}

var expectedMessages = []*Message{
	{
		ID:      "qwki97",
		FullID:  "t4_qwki97",
		Created: &Timestamp{time.Date(2020, 8, 18, 0, 16, 53, 0, time.UTC)},

		Subject:  "re: test",
		Text:     "test",
		ParentID: "t4_qwki4m",

		Author: "testuser1",
		To:     "testuser2",

		IsComment: false,
	},
}

func TestMessageService_ReadAll(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/read_all_messages", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.Message.ReadAll(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestMessageService_Read(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/read_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Message.Read(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Read(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Unread(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unread_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Message.Unread(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Unread(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Block(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/block", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Message.Block(ctx, "test")
	require.NoError(t, err)
}

func TestMessageService_Collapse(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/collapse_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Message.Collapse(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Collapse(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Uncollapse(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/uncollapse_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Message.Uncollapse(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Uncollapse(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Delete(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/del_msg", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Message.Delete(ctx, "test")
	require.NoError(t, err)
}

func TestMessageService_Send(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/compose", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("to", "test")
		form.Set("subject", "test subject")
		form.Set("text", "test text")
		form.Set("from_sr", "hello world")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Message.Send(ctx, nil)
	require.EqualError(t, err, "*SendMessageRequest: cannot be nil")

	_, err = client.Message.Send(ctx, &SendMessageRequest{
		To:            "test",
		Subject:       "test subject",
		Text:          "test text",
		FromSubreddit: "hello world",
	})
	require.NoError(t, err)
}

func TestMessageService_Inbox(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/message/inbox.json")
	require.NoError(t, err)

	mux.HandleFunc("/message/inbox", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	comments, messages, _, err := client.Message.Inbox(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedCommentMessages, comments)
	require.Equal(t, expectedMessages, messages)
}

func TestMessageService_InboxUnread(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/message/inbox.json")
	require.NoError(t, err)

	mux.HandleFunc("/message/unread", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	comments, messages, _, err := client.Message.InboxUnread(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedCommentMessages, comments)
	require.Equal(t, expectedMessages, messages)
}

func TestMessageService_Sent(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/message/inbox.json")
	require.NoError(t, err)

	mux.HandleFunc("/message/sent", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	messages, _, err := client.Message.Sent(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedMessages, messages)
}
