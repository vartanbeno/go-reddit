package reddit

import (
	"context"
	"sync"
	"time"
)

// StreamService allows streaming new content from Reddit as it appears.
type StreamService struct {
	client *Client
}

// Posts streams posts from the specified subreddit.
// It returns 2 channels and a function:
//   - a channel into which new posts will be sent
//   - a channel into which any errors will be sent
//   - a function that the client can call once to stop the streaming and close the channels
// Because of the 100 post limit imposed by Reddit when fetching posts, some high-traffic
// streams might drop submissions between API requests, such as when streaming r/all.
func (s *StreamService) Posts(subreddit string) (<-chan *Post, <-chan error, func()) {
	posts := make(chan *Post)
	errs := make(chan error)
	ticker := time.NewTicker(time.Second * 5)

	// originally used the "before" parameter, but if that post gets deleted, subsequent requests
	// would just return empty listings; easier to just keep track of all post ids encountered
	ids := set{}

	go func() {
		for ; ; <-ticker.C {
			result, err := s.getPosts(subreddit)
			if err != nil {
				errs <- err
				continue
			}

			for _, post := range result.Posts {
				id := post.FullID

				// if this post id is already part of the set, it means that it and the ones
				// after it in the list have already been streamed, so break out of the loop
				if ids.Exists(id) {
					break
				}

				ids.Add(id)
				posts <- post
			}
		}
	}()

	var once sync.Once
	return posts, errs, func() {
		once.Do(func() {
			ticker.Stop()
			close(posts)
			close(errs)
		})
	}
}

func (s *StreamService) getPosts(subreddit string) (*Posts, error) {
	opts := &ListOptions{
		Limit: 100,
	}

	result, _, err := s.client.Subreddit.NewPosts(context.Background(), subreddit, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type set map[string]struct{}

func (s set) Add(v string) {
	s[v] = struct{}{}
}

func (s set) Delete(v string) {
	delete(s, v)
}

func (s set) Len() int {
	return len(s)
}

func (s set) Exists(v string) bool {
	_, ok := s[v]
	return ok
}
