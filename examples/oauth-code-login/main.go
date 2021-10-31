package main

import (
	"context"
	"net/http"

	"github.com/rgood/go-reddit/reddit"
)

func main() {
	client, _ := reddit.NewTokenClient(
		"YIB6CeMwOJ9fvsmh4fIZjw",
		"JsIpBV3Mz0V7Ag6xkVOAjAeoa12ZQw",
		"http://localhost:1337/auth/callback",
		[]string{"identity"},
		reddit.WithUserAgent("Go-Reddit Oauth2 (by /u/The1RGood)"),
	)

	url, _ := client.AuthorizeURL("test")
	println(url)

	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		// Log request received
		println("Handling request...")

		// Create a new client in the request context
		userClient, _ := reddit.NewTokenClient(
			"YIB6CeMwOJ9fvsmh4fIZjw",
			"JsIpBV3Mz0V7Ag6xkVOAjAeoa12ZQw",
			"http://localhost:1337/auth/callback",
			[]string{"identity"},
			reddit.WithUserAgent("Go-Reddit Oauth2 (by /u/The1RGood)"),
		)

		// Retrieve code from request
		code := r.FormValue("code")

		// Authorize client with the code
		userClient.Authorize(code)

		// Fetch user info of authorized client
		user, _, _ := userClient.Account.Info(context.Background())

		// Write the user info in the response
		w.Write([]byte("Authorized as: " + user.Name))

		// Log complete
		println("Request completed.")
	})

	http.ListenAndServe(":1337", nil)
}
