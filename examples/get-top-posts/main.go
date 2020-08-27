package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vartanbeno/go-reddit/reddit"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	credentials := &reddit.Credentials{
		ID:       "id",
		Secret:   "secret",
		Username: "username",
		Password: "password",
	}

	client, err := reddit.NewClient(credentials)
	if err != nil {
		return
	}

	// Let's get the top 200 posts of r/golang.
	// Reddit returns a maximum of 100 posts at a time,
	// so we'll need to separate this into 2 requests.
	result, _, err := client.Subreddit.TopPosts(ctx, "golang", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 100,
		},
		Time: "all",
	})
	if err != nil {
		return
	}

	for _, post := range result.Posts {
		fmt.Println(post.Title)
	}

	// The SetAfter option sets the id of an item that Reddit
	// will use as an anchor point for the returned listing.
	result, _, err = client.Subreddit.TopPosts(ctx, "golang", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 100,
			After: result.After,
		},
		Time: "all",
	})
	if err != nil {
		return
	}

	for _, post := range result.Posts {
		fmt.Println(post.Title)
	}

	return
}
