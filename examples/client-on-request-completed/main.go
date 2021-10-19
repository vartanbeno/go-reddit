package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rgood/go-reddit/reddit"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	client, err := reddit.NewReadonlyClient()
	if err != nil {
		return
	}

	client.OnRequestCompleted(logResponse)

	client.Subreddit.Search(ctx, "programming", nil)
	client.Subreddit.SearchNames(ctx, "monitor")
	client.Subreddit.SearchPosts(ctx, "react", "webdev", nil)
	client.Subreddit.HotPosts(ctx, "golang", &reddit.ListOptions{Limit: 5})

	return
}

func logResponse(req *http.Request, res *http.Response) {
	fmt.Printf("%s %s %s\n", req.Method, req.URL, res.Status)
}
