# go-reddit

[![Actions Status](https://github.com/vartanbeno/go-reddit/workflows/tests/badge.svg)](https://github.com/vartanbeno/go-reddit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/vartanbeno/go-reddit)](https://goreportcard.com/report/github.com/vartanbeno/go-reddit)

go-reddit is a Go client library for accessing the Reddit API.

You can view Reddit's official API documentation [here](https://www.reddit.com/dev/api/).

## Install

To get a specific version from the list of [versions](https://github.com/vartanbeno/go-reddit/releases):

```sh
go get github.com/vartanbeno/go-reddit@vX.Y.Z
```

Or for the latest version:

```sh
go get github.com/vartanbeno/go-reddit
```

## Usage

Make sure to have a Reddit app with a valid client id and secret. [Here](https://github.com/reddit-archive/reddit/wiki/OAuth2-Quick-Start-Example#first-steps) is a quick guide on how to create an app and get credentials.

```go
package main

import "github.com/vartanbeno/go-reddit/reddit"

func main() {
    credentials := &reddit.Credentials{
        ID:       "id",
        Secret:   "secret",
        Username: "username",
        Password: "password",
    }
    client, _ := reddit.NewClient(credentials)
}
```

You can pass in a number of options to `NewClient` that further configure the client. For example, to use a custom HTTP client:

```go
httpClient := &http.Client{Timeout: time.Second * 30}
client, _ := reddit.NewClient(credentials, reddit.WithHTTPClient(httpClient))
```

If this option is not used, it will be set to `&http.Client{}` by default. More options are available in the [reddit/reddit-options.go](reddit/reddit-options.go) file.

## Examples

<details>
    <summary>Submit a comment.</summary>

```go
comment, _, err := client.Comment.Submit(context.Background(), "t3_postid", "comment body")
if err != nil {
    return err
}
fmt.Printf("Comment permalink: %s\n", comment.Permalink)
```
</details>

<details>
    <summary>Upvote a post.</summary>

```go
_, err := client.Post.Upvote(context.Background(), "t3_postid")
if err != nil {
    return err
}
```
</details>

<details>
    <summary>Get r/golang's top 5 posts of all time.</summary>

```go
result, _, err := client.Subreddit.TopPosts(context.Background(), "golang", &reddit.ListPostOptions{
    ListOptions: reddit.ListOptions{
        Limit: 5,
    },
    Time: "all",
})
if err != nil {
    return err
}
fmt.Printf("Received %d posts.\n", len(result.Posts))
```
</details>

More examples are available in the [examples](examples) folder.

## Design

The package design is heavily inspired from [Google's GitHub API client](https://github.com/google/go-github) and [DigitalOcean's API client](https://github.com/digitalocean/godo).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
