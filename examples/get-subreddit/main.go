package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kmulvey/reddit/v2/reddit"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	sr, _, err := reddit.DefaultClient().Subreddit.Get(ctx, "golang")
	if err != nil {
		return
	}

	fmt.Printf("%s was created on %s and has %d subscribers.\n", sr.DisplayNamePrefixed, sr.Created.Local(), sr.Subscribers)
	return
}
