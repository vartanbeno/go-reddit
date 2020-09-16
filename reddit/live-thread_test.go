package reddit

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedLiveThread = &LiveThread{
	ID:      "15nevtv8e54dh",
	FullID:  "LiveUpdateEvent_15nevtv8e54dh",
	Created: &Timestamp{time.Date(2020, 9, 16, 1, 20, 27, 0, time.UTC)},

	Title:       "test",
	Description: "test",
	Resources:   "",

	State:             "live",
	ViewerCount:       6,
	ViewerCountFuzzed: true,

	WebSocketURL: "wss://ws-078adc7cb2099a9df.wss.redditmedia.com/live/15nevtv8e54dh?m=AQAA7rxiX6EpLYFCFZ0KJD4lVAPaMt0A1z2-xJ1b2dWCmxNIfMwL",

	Announcement: false,
	NSFW:         false,
}

func TestLiveThreadService_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	blob, err := readFileContents("../testdata/live-thread/live-thread.json")
	require.NoError(t, err)

	mux.HandleFunc("/live/id123/about", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	liveThread, _, err := client.LiveThread.Get(ctx, "id123")
	require.NoError(t, err)
	require.Equal(t, expectedLiveThread, liveThread)
}
