package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rgood/go-reddit/reddit"
)

var ctx = context.Background()

func main() {
	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	posts, errs, stop := reddit.DefaultClient().Stream.Posts("AskReddit", reddit.StreamInterval(time.Second*3), reddit.StreamDiscardInitial)
	defer stop()

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

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
		case rcvSig, ok := <-sig:
			if !ok {
				return
			}
			fmt.Printf("Stopping due to %s signal.\n", rcvSig)
			return
		case <-timer.C:
			return
		}
	}
}
