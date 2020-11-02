package reddit

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoldService_Gild(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/gold/gild/t1_test", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
	})

	_, err := client.Gold.Gild(ctx, "t1_test")
	require.NoError(t, err)
}

func TestGoldService_Give(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/v1/gold/give/testuser", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("months", "1")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Gold.Give(ctx, "testuser", 0)
	require.EqualError(t, err, "months: must be between 1 and 36 (inclusive)")

	_, err = client.Gold.Give(ctx, "testuser", 37)
	require.EqualError(t, err, "months: must be between 1 and 36 (inclusive)")

	_, err = client.Gold.Give(ctx, "testuser", 1)
	require.NoError(t, err)
}
