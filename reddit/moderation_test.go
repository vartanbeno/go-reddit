package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var expectedModActions = []*ModAction{
	{
		ID:      "ModAction_b4e7979a-c4ad-11ea-8440-0ea1b7c2b8f9",
		Action:  "spamcomment",
		Created: &Timestamp{time.Date(2020, 7, 13, 2, 8, 14, 0, time.UTC)},

		Moderator:   "v_95",
		ModeratorID: "164ab8",

		TargetAuthor:    "testuser",
		TargetID:        "t1_fxw10aa",
		TargetPermalink: "/r/helloworldtestt/comments/hq6r3t/yo/fxw10aa/",
		TargetBody:      "hi",

		Subreddit:   "helloworldtestt",
		SubredditID: "2uquw1",
	},
	{
		ID:      "ModAction_a0408162-c4ad-11ea-8239-0e3b48262e8b",
		Action:  "sticky",
		Created: &Timestamp{time.Date(2020, 7, 13, 2, 7, 38, 0, time.UTC)},

		Moderator:   "v_95",
		ModeratorID: "164ab8",

		TargetAuthor:    "testuser",
		TargetID:        "t3_hq6r3t",
		TargetTitle:     "yo",
		TargetPermalink: "/r/helloworldtestt/comments/hq6r3t/yo/",

		Subreddit:   "helloworldtestt",
		SubredditID: "2uquw1",
	},
}

func TestModerationService_Actions(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/moderation/actions.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/log", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		form := url.Values{}
		form.Set("type", "testtype")
		form.Set("mod", "testmod")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.Form)

		fmt.Fprint(w, blob)
	})

	modActions, resp, err := client.Moderation.Actions(ctx, "testsubreddit", &ListModActionOptions{Type: "testtype", Moderator: "testmod"})
	require.NoError(t, err)
	require.Equal(t, expectedModActions, modActions)
	require.Equal(t, "ModAction_a0408162-c4ad-11ea-8239-0e3b48262e8b", resp.After)
}

func TestModerationService_AcceptInvite(t *testing.T) {
	client, mux := setup(t)

	blob, err := readFileContents("../testdata/moderation/actions.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/api/accept_moderator_invite", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)

		fmt.Fprint(w, blob)
	})

	_, err = client.Moderation.AcceptInvite(ctx, "testsubreddit")
	require.NoError(t, err)
}

func TestModerationService_Approve(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/approve", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Approve(ctx, "t3_test")
	require.NoError(t, err)
}

func TestModerationService_Remove(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/remove", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("spam", "false")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Remove(ctx, "t3_test")
	require.NoError(t, err)
}

func TestModerationService_RemoveSpam(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/remove", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")
		form.Set("spam", "true")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.RemoveSpam(ctx, "t3_test")
	require.NoError(t, err)
}

func TestModerationService_Leave(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/leavemoderator", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t5_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Leave(ctx, "t5_test")
	require.NoError(t, err)
}

func TestModerationService_LeaveContributor(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/leavecontributor", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t5_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.LeaveContributor(ctx, "t5_test")
	require.NoError(t, err)
}

func TestModerationService_Reported(t *testing.T) {
	client, mux := setup(t)

	// contains posts and comments
	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/reports", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, resp, err := client.Moderation.Reported(ctx, "testsubreddit", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t1_f0zsa37", resp.After)

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])
	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestModerationService_Spam(t *testing.T) {
	client, mux := setup(t)

	// contains posts and comments
	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/spam", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, resp, err := client.Moderation.Spam(ctx, "testsubreddit", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t1_f0zsa37", resp.After)

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])
	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestModerationService_Queue(t *testing.T) {
	client, mux := setup(t)

	// contains posts and comments
	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/modqueue", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, resp, err := client.Moderation.Queue(ctx, "testsubreddit", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t1_f0zsa37", resp.After)

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])
	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestModerationService_Unmoderated(t *testing.T) {
	client, mux := setup(t)

	// contains posts and comments
	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/unmoderated", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, resp, err := client.Moderation.Unmoderated(ctx, "testsubreddit", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestModerationService_Edited(t *testing.T) {
	client, mux := setup(t)

	// contains posts and comments
	blob, err := readFileContents("../testdata/user/overview.json")
	require.NoError(t, err)

	mux.HandleFunc("/r/testsubreddit/about/edited", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, blob)
	})

	posts, comments, resp, err := client.Moderation.Edited(ctx, "testsubreddit", nil)
	require.NoError(t, err)

	require.Len(t, posts, 1)
	require.Equal(t, expectedPost, posts[0])
	require.Equal(t, "t1_f0zsa37", resp.After)

	require.Len(t, comments, 1)
	require.Equal(t, expectedComment, comments[0])
	require.Equal(t, "t1_f0zsa37", resp.After)
}

func TestModerationService_IgnoreReports(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/ignore_reports", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.IgnoreReports(ctx, "t3_test")
	require.NoError(t, err)
}

func TestModerationService_UnignoreReports(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/unignore_reports", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("id", "t3_test")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.UnignoreReports(ctx, "t3_test")
	require.NoError(t, err)
}

func TestModerationService_Invite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/friend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "moderator_invite")
		form.Set("permissions", "+all")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Invite(ctx, "testsubreddit", "testuser", nil)
	require.NoError(t, err)
}

func TestModerationService_Invite_Permissions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/friend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "moderator_invite")
		form.Set("permissions", "-all,-access,-chat_config,+chat_operator,+config,-flair,-mail,-posts,+wiki")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Invite(ctx, "testsubreddit", "testuser", &ModPermissions{ChatOperator: true, Config: true, Wiki: true})
	require.NoError(t, err)
}

func TestModerationService_Uninvite(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "moderator_invite")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Uninvite(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_SetPermissions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/setpermissions", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "moderator_invite")
		form.Set("permissions", "-all,+access,-chat_config,-chat_operator,-config,+flair,-mail,+posts,-wiki")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.SetPermissions(ctx, "testsubreddit", "testuser", &ModPermissions{Access: true, Flair: true, Posts: true})
	require.NoError(t, err)
}

func TestModerationService_Ban(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/friend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "banned")
		form.Set("reason", "test reason")
		form.Set("note", "test mod note")
		form.Set("duration", "5")
		form.Set("ban_message", "test message")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Ban(ctx, "testsubreddit", "testuser", &BanConfig{
		Reason:  "test reason",
		ModNote: "test mod note",
		Days:    Int(5),
		Message: "test message",
	})
	require.NoError(t, err)
}

func TestModerationService_Unban(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "banned")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Unban(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_BanWiki(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/friend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "wikibanned")
		form.Set("reason", "test reason")
		form.Set("note", "test mod note")
		form.Set("duration", "5")
		form.Set("ban_message", "test message")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.BanWiki(ctx, "testsubreddit", "testuser", &BanConfig{
		Reason:  "test reason",
		ModNote: "test mod note",
		Days:    Int(5),
		Message: "test message",
	})
	require.NoError(t, err)
}

func TestModerationService_UnbanWiki(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "wikibanned")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.UnbanWiki(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_Mute(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/friend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "muted")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Mute(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_Unmuted(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "muted")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Unmute(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_ApproveUser(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/friend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "contributor")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.ApproveUser(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_UnapproveUser(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "contributor")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.UnapproveUser(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_ApproveUserWiki(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/friend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "wikicontributor")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.ApproveUserWiki(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_UnapproveUserWiki(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/r/testsubreddit/api/unfriend", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("name", "testuser")
		form.Set("type", "wikicontributor")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.UnapproveUserWiki(ctx, "testsubreddit", "testuser")
	require.NoError(t, err)
}

func TestModerationService_Distinguish(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/distinguish", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("how", "yes")
		form.Set("id", "t1_123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Distinguish(ctx, "t1_123")
	require.NoError(t, err)
}

func TestModerationService_DistinguishAndSticky(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/distinguish", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("how", "yes")
		form.Set("sticky", "true")
		form.Set("id", "t1_123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.DistinguishAndSticky(ctx, "t1_123")
	require.NoError(t, err)
}

func TestModerationService_Undistinguish(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/api/distinguish", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		form := url.Values{}
		form.Set("api_type", "json")
		form.Set("how", "no")
		form.Set("id", "t1_123")

		err := r.ParseForm()
		require.NoError(t, err)
		require.Equal(t, form, r.PostForm)
	})

	_, err := client.Moderation.Undistinguish(ctx, "t1_123")
	require.NoError(t, err)
}
