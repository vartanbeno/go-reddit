package reddit

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStreamService_Posts(t *testing.T) {
	client, mux := setup(t)

	var counter int
	mux.HandleFunc("/r/testsubreddit/new", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		defer func() { counter++ }()

		switch counter {
		case 0:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post1"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post2"
							}
						}
					]
				}
			}`)
		case 1:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post3"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post1"
							}
						}
					]
				}
			}`)
		case 2:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post4"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post5"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post6"
							}
						}
					]
				}
			}`)
		case 3:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post7"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post8"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post9"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post10"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post11"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post12"
							}
						}
					]
				}
			}`)
		default:
			fmt.Fprint(w, `{}`)
		}
	})

	posts, errs, stop := client.Stream.Posts("testsubreddit", StreamInterval(time.Millisecond*10), StreamMaxRequests(4))
	defer stop()

	expectedPostIDs := []string{"t3_post1", "t3_post2", "t3_post3", "t3_post4", "t3_post5", "t3_post6", "t3_post7", "t3_post8", "t3_post9", "t3_post10", "t3_post11", "t3_post12"}
	var i int

loop:
	for i != len(expectedPostIDs) {
		select {
		case post, ok := <-posts:
			if !ok {
				break loop
			}
			require.Equal(t, expectedPostIDs[i], post.FullID)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			require.NoError(t, err)
		}
		i++
	}

	require.Len(t, expectedPostIDs, i)
}

func TestStreamService_Posts_DiscardInitial(t *testing.T) {
	client, mux := setup(t)

	var counter int
	mux.HandleFunc("/r/testsubreddit/new", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		defer func() { counter++ }()

		switch counter {
		case 0:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post1"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post2"
							}
						}
					]
				}
			}`)
		case 1:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post3"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post1"
							}
						}
					]
				}
			}`)
		case 2:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post4"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post5"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post6"
							}
						}
					]
				}
			}`)
		case 3:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t3",
							"data": {
								"name": "t3_post7"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post8"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post9"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post10"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post11"
							}
						},
						{
							"kind": "t3",
							"data": {
								"name": "t3_post12"
							}
						}
					]
				}
			}`)
		default:
			fmt.Fprint(w, `{}`)
		}
	})

	posts, errs, stop := client.Stream.Posts("testsubreddit", StreamInterval(time.Millisecond*10), StreamMaxRequests(4), StreamDiscardInitial)
	defer stop()

	expectedPostIDs := []string{"t3_post3", "t3_post4", "t3_post5", "t3_post6", "t3_post7", "t3_post8", "t3_post9", "t3_post10", "t3_post11", "t3_post12"}
	var i int

loop:
	for i != len(expectedPostIDs) {
		select {
		case post, ok := <-posts:
			if !ok {
				break loop
			}
			require.Equal(t, expectedPostIDs[i], post.FullID)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			require.NoError(t, err)
		}
		i++
	}

	require.Len(t, expectedPostIDs, i)
}

func TestStreamService_Comments(t *testing.T) {
	client, mux := setup(t)

	var counter int
	mux.HandleFunc("/r/testsubreddit/comments", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		defer func() { counter++ }()

		switch counter {
		case 0:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment1"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment2"
							}
						}
					]
				}
			}`)
		case 1:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment3"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment1"
							}
						}
					]
				}
			}`)
		case 2:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment4"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment5"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment6"
							}
						}
					]
				}
			}`)
		case 3:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment7"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment8"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment9"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment10"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment11"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment12"
							}
						}
					]
				}
			}`)
		default:
			fmt.Fprint(w, `{}`)
		}
	})

	comments, errs, stop := client.Stream.Comments("testsubreddit", StreamInterval(time.Millisecond*10), StreamMaxRequests(4))
	defer stop()

	expectedCommentIds := []string{"t1_comment1", "t1_comment2", "t1_comment3", "t1_comment4", "t1_comment5", "t1_comment6", "t1_comment7", "t1_comment8", "t1_comment9", "t1_comment10", "t1_comment11", "t1_comment12"}
	var i int

loop:
	for i != len(expectedCommentIds) {
		select {
		case comment, ok := <-comments:
			if !ok {
				break loop
			}
			require.Equal(t, expectedCommentIds[i], comment.FullID)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			require.NoError(t, err)
		}
		i++
	}

	require.Len(t, expectedCommentIds, i)
}

func TestStreamService_CommentsDiscardInitial(t *testing.T) {
	client, mux := setup(t)

	var counter int
	mux.HandleFunc("/r/testsubreddit/comments", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		defer func() { counter++ }()

		switch counter {
		case 0:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment1"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment2"
							}
						}
					]
				}
			}`)
		case 1:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment3"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment1"
							}
						}
					]
				}
			}`)
		case 2:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment4"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment5"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment6"
							}
						}
					]
				}
			}`)
		case 3:
			fmt.Fprint(w, `{
				"kind": "Listing",
				"data": {
					"children": [
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment7"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment8"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment9"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment10"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment11"
							}
						},
						{
							"kind": "t1",
							"data": {
								"name": "t1_comment12"
							}
						}
					]
				}
			}`)
		default:
			fmt.Fprint(w, `{}`)
		}
	})

	comments, errs, stop := client.Stream.Comments("testsubreddit", StreamInterval(time.Millisecond*10), StreamMaxRequests(4), StreamDiscardInitial)
	defer stop()

	expectedCommentIds := []string{"t1_comment3", "t1_comment4", "t1_comment5", "t1_comment6", "t1_comment7", "t1_comment8", "t1_comment9", "t1_comment10", "t1_comment11", "t1_comment12"}
	var i int

loop:
	for i != len(expectedCommentIds) {
		select {
		case comment, ok := <-comments:
			if !ok {
				break loop
			}
			require.Equal(t, expectedCommentIds[i], comment.FullID)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			require.NoError(t, err)
		}
		i++
	}

	require.Len(t, expectedCommentIds, i)
}
