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
    withCredentials := reddit.WithCredentials("id", "secret", "username", "password")
    client, _ := reddit.NewClient(nil, withCredentials)
}
```

The first argument (the one set to `nil`) is of type `*http.Client`. It will be used to make the requests. If nil, it will be set to `&http.Client{}`.

The `WithCredentials` option sets the credentials used to make requests to the Reddit API.

## Examples

<details>
    <summary>Configure the client from environment variables.</summary>

```go
client, _ := reddit.NewClient(nil, reddit.FromEnv)
```
</details>

<details>
    <summary>Upvote a post.</summary>

```go
_, err := client.Post.Upvote(context.Background(), "t3_postid")
if err != nil {
    fmt.Printf("Something bad happened: %v\n", err)
    return err
}
```
</details>

<details>
    <summary>Get r/golang's top 5 posts of all time.</summary>

```go
result, _, err := client.Subreddit.Top(context.Background(), "golang", &reddit.ListPostOptions{
    ListOptions: reddit.ListOptions{
        Limit: 5,
    },
    Time: "all",
})
if err != nil {
    fmt.Printf("Something bad happened: %v\n", err)
    return err
}
fmt.Printf("Received %d posts.\n", len(result.Posts))
```
</details>

More examples are available in the [examples](examples) folder.

## Design

The package design and structure are heavily inspired from [Google's GitHub API client](https://github.com/google/go-github) and [DigitalOcean's API client](https://github.com/digitalocean/godo).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
