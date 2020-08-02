package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vartanbeno/go-reddit"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	withCredentials := reddit.WithCredentials("id", "secret", "username", "password")

	client, err := reddit.NewClient(nil, withCredentials)
	if err != nil {
		return
	}

	// Let's get the top 200 posts of r/golang.
	// Reddit returns a maximum of 100 posts at a time,
	// so we'll need to separate this into 2 requests.
	result, _, err := client.Subreddit.Top(ctx, "golang", reddit.SetLimit(100), reddit.FromAllTime)
	if err != nil {
		return
	}

	for _, post := range result.Posts {
		fmt.Println(post.Title)
	}

	// The SetAfter option sets the id of an item that Reddit
	// will use as an anchor point for the returned listing.
	result, _, err = client.Subreddit.Top(ctx, "golang", reddit.SetLimit(100), reddit.FromAllTime, reddit.SetAfter(result.After))
	if err != nil {
		return
	}

	for _, post := range result.Posts {
		fmt.Println(post.Title)
	}

	return
}
