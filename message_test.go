package reddit

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageService_ReadAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/read_all_messages", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusAccepted)
	})

	res, err := client.Message.ReadAll(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)
}

func TestMessageService_Read(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/read_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.Message.Read(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Read(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Unread(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/unread_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.Message.Unread(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Unread(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Block(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/block", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.Message.Block(ctx, "test")
	require.NoError(t, err)
}

func TestMessageService_Collapse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/collapse_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.Message.Collapse(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Collapse(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Uncollapse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/uncollapse_message", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test1,test2,test3")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.Message.Uncollapse(ctx)
	require.EqualError(t, err, "must provide at least 1 id")

	_, err = client.Message.Uncollapse(ctx, "test1", "test2", "test3")
	require.NoError(t, err)
}

func TestMessageService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/del_msg", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)
	})

	_, err := client.Message.Delete(ctx, "test")
	require.NoError(t, err)
}
