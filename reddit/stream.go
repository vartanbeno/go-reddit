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
func (s *StreamService) Posts(subreddit string, opts ...StreamOpt) (<-chan *Post, <-chan error, func()) {
	streamConfig := &streamConfig{
		Interval:       defaultStreamInterval,
		DiscardInitial: false,
		MaxRequests:    0,
	}
	for _, opt := range opts {
		opt(streamConfig)
	}

	ticker := time.NewTicker(streamConfig.Interval)
	postsCh := make(chan *Post)
	errsCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	// originally used the "before" parameter, but if that post gets deleted, subsequent requests
	// would just return empty listings; easier to just keep track of all post ids encountered
	ids := NewOrderedMaxSet(2000)

	go func() {
		defer close(postsCh)
		defer close(errsCh)
		defer cancel()
		var wg sync.WaitGroup
		defer wg.Wait()
		var mutex sync.Mutex

		var n int
		infinite := streamConfig.MaxRequests == 0

		for ; ; <-ticker.C {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:
			}

			n++
			wg.Add(1)
			go s.getPosts(ctx, subreddit, func(posts []*Post, err error) {
				defer wg.Done()

				if err != nil {
					select {
					case <-ctx.Done():
					default:
						errsCh <- err
					}
					return
				}

				for _, post := range posts {
					id := post.FullID

					// if this post id is already part of the set, it means that it and the ones
					// after it in the list have already been streamed, so break out of the loop
					if ids.Exists(id) {
						break
					}
					ids.Add(id)

					if func() bool {
						mutex.Lock()
						toReturn := false
						if streamConfig.DiscardInitial {
							streamConfig.DiscardInitial = false
							toReturn = true
						}
						mutex.Unlock()
						return toReturn
					}() {
						break
					}

					select {
					case <-ctx.Done():
						return
					default:
						postsCh <- post
					}
				}
			})

			if !infinite && n >= streamConfig.MaxRequests {
				break
			}
		}
	}()

	return postsCh, errsCh, cancel
}

// Comments streams comments from the specified subreddit.
// It returns 2 channels and a function:
//   - a channel into which new comments will be sent
//   - a channel into which any errors will be sent
//   - a function that the client can call once to stop the streaming and close the channels
// Because of the 100 result limit imposed by Reddit when fetching posts, some high-traffic
// streams might drop submissions between API requests, such as when streaming r/all.
func (s *StreamService) Comments(subreddit string, opts ...StreamOpt) (<-chan *Comment, <-chan error, func()) {
	streamConfig := &streamConfig{
		Interval:       defaultStreamInterval,
		DiscardInitial: false,
		MaxRequests:    0,
	}
	for _, opt := range opts {
		opt(streamConfig)
	}

	ticker := time.NewTicker(streamConfig.Interval)
	commentsCh := make(chan *Comment)
	errsCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	ids := NewOrderedMaxSet(2000)

	go func() {
		defer close(commentsCh)
		defer close(errsCh)
		defer cancel()
		var wg sync.WaitGroup
		defer wg.Wait()
		var mutex sync.Mutex

		var n int
		infinite := streamConfig.MaxRequests == 0

		for ; ; <-ticker.C {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:
			}
			n++
			wg.Add(1)

			go s.getComments(ctx, subreddit, func(comments []*Comment, err error) {
				defer wg.Done()
				if err != nil {
					select {
					case <-ctx.Done():
					default:
						errsCh <- err
					}
					return
				}

				for _, comment := range comments {
					id := comment.FullID

					// certain comment streams are inconsistent about the completeness of returned comments
					// it's not enough to check if we've seen older comments, but we must check for every comment individually
					if !ids.Exists(id) {
						ids.Add(id)

						if func() bool {
							mutex.Lock()
							toReturn := false
							if streamConfig.DiscardInitial {
								streamConfig.DiscardInitial = false
								toReturn = true
							}
							mutex.Unlock()
							return toReturn
						}() {
							break
						}

						select {
						case <-ctx.Done():
							return
						default:
							commentsCh <- comment
						}
					}

				}
			})
			if !infinite && n >= streamConfig.MaxRequests {
				break
			}
		}
	}()

	return commentsCh, errsCh, cancel
}

func (s *StreamService) getPosts(ctx context.Context, subreddit string, cb func([]*Post, error)) {
	posts, _, err := s.client.Subreddit.NewPosts(ctx, subreddit, &ListOptions{Limit: 100})
	cb(posts, err)
}

func (s *StreamService) getComments(ctx context.Context, subreddit string, cb func([]*Comment, error)) {
	comments, _, err := s.client.Subreddit.Comments(ctx, subreddit, &ListOptions{Limit: 100})
	cb(comments, err)
}
