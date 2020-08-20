package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/vartanbeno/go-reddit/reddit"
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

	client.OnRequestCompleted(logResponse)

	client.Subreddit.Search(ctx, "programming", nil)
	client.Subreddit.SearchNames(ctx, "monitor")
	client.Subreddit.SearchPosts(ctx, "react", "webdev", nil)
	client.User.Posts(ctx, &reddit.ListUserOverviewOptions{
		ListOptions: reddit.ListOptions{
			Limit: 50,
		},
		Sort: "top",
		Time: "month",
	})

	return
}

func logResponse(req *http.Request, res *http.Response) {
	fmt.Printf("%s %s %s\n", req.Method, req.URL, res.Status)
}
