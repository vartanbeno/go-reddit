package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	posts, errs, stop := client.Stream.Posts("AskReddit", reddit.StreamInterval(time.Second*3), reddit.StreamDiscardInitial)
	defer stop()

	go func() {
		for {
			select {
			case post, ok := <-posts:
				if !ok {
					return
				}
				fmt.Printf("Received post: %s\n", post.Title)
			case err, ok := <-errs:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "Error! %v\n", err)
			}
		}
	}()

	<-time.After(time.Minute)
	return
}
