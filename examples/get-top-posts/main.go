package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rgood/go-reddit/reddit"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	// Let's get the top 200 posts of r/golang.
	// Reddit returns a maximum of 100 posts at a time,
	// so we'll need to separate this into 2 requests.
	posts, resp, err := reddit.DefaultClient().Subreddit.TopPosts(ctx, "golang", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 100,
		},
		Time: "all",
	})
	if err != nil {
		return
	}

	for _, post := range posts {
		fmt.Println(post.Title)
	}

	// The After option sets the id of an item that Reddit
	// will use as an anchor point for the returned listing.
	posts, _, err = reddit.DefaultClient().Subreddit.TopPosts(ctx, "golang", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 100,
			After: resp.After,
		},
		Time: "all",
	})
	if err != nil {
		return
	}

	for _, post := range posts {
		fmt.Println(post.Title)
	}

	return
}
