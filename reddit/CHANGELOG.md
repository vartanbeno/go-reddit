# Change Log

## [v2.0.0] - 2021-01-24

- The underlying `*http.Client` is now passed as an option when initializing a client.
- Add `DefaultClient()` method which returns a valid, read-only client with limited access to the Reddit API.
- Remove the `before` anchor from responses.
- The decoding process for a listing response has been revamped, and the `after` anchor is now included in the `*Response` object. For example, when fetching a user's comments, instead of getting a struct containing the list of comments, and the `after` value, you now get the list of comments directly, and the `after` can be obtained from the `*Response`.
- Create `WikiService`, `LiveThreadService`, and `WidgetService`.
- Add more methods to `SubredditService`, `ModerationService`, and `FlairService`.
- Add error handling for rate limit errors (`*RateLimitError`).

## [v1.0.0] - 2020-08-25

- Most endpoints outlined in the [official Reddit API documentation](https://www.reddit.com/dev/api/) have been implemented!
