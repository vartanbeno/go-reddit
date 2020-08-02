package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	client.OnRequestCompleted(logResponse)

	client.Subreddit.Search(ctx, "programming", reddit.SetLimit(10))
	client.Subreddit.SearchNames(ctx, "monitor")
	client.Subreddit.SearchPosts(ctx, "react", "webdev", reddit.SortByNumberOfComments)
	client.User.Posts(ctx, reddit.SetLimit(50))

	return
}

func logResponse(req *http.Request, res *http.Response) {
	fmt.Printf("%s %s %s\n", req.Method, req.URL, res.Status)
}
