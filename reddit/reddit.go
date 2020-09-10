package reddit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
)

const (
	libraryName    = "github.com/vartanbeno/go-reddit"
	libraryVersion = "1.0.0"

	defaultBaseURL         = "https://oauth.reddit.com"
	defaultBaseURLReadonly = "https://reddit.com"
	defaultTokenURL        = "https://www.reddit.com/api/v1/access_token"

	mediaTypeJSON = "application/json"
	mediaTypeForm = "application/x-www-form-urlencoded"

	headerContentType = "Content-Type"
	headerAccept      = "Accept"
	headerUserAgent   = "User-Agent"
)

// DefaultClient is a readonly client with limited access to the Reddit API.
var DefaultClient, _ = NewReadonlyClient()

// RequestCompletionCallback defines the type of the request callback function.
type RequestCompletionCallback func(*http.Request, *http.Response)

// Credentials used to authenticate to make requests to the Reddit API.
type Credentials struct {
	ID       string
	Secret   string
	Username string
	Password string
}

// Client manages communication with the Reddit API.
type Client struct {
	// HTTP client used to communicate with the Reddit API.
	client *http.Client

	BaseURL  *url.URL
	TokenURL *url.URL

	userAgent string

	ID       string
	Secret   string
	Username string
	Password string

	// This is the client's user ID in Reddit's database.
	redditID string

	Account    *AccountService
	Collection *CollectionService
	Comment    *CommentService
	Emoji      *EmojiService
	Flair      *FlairService
	Gold       *GoldService
	Listings   *ListingsService
	Message    *MessageService
	Moderation *ModerationService
	Multi      *MultiService
	Post       *PostService
	Stream     *StreamService
	Subreddit  *SubredditService
	User       *UserService
	Wiki       *WikiService

	oauth2Transport *oauth2.Transport

	onRequestCompleted RequestCompletionCallback
}

// OnRequestCompleted sets the client's request completion callback.
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

func newClient() *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	tokenURL, _ := url.Parse(defaultTokenURL)

	client := &Client{BaseURL: baseURL, TokenURL: tokenURL}

	client.Account = &AccountService{client: client}
	client.Collection = &CollectionService{client: client}
	client.Emoji = &EmojiService{client: client}
	client.Flair = &FlairService{client: client}
	client.Gold = &GoldService{client: client}
	client.Listings = &ListingsService{client: client}
	client.Message = &MessageService{client: client}
	client.Moderation = &ModerationService{client: client}
	client.Multi = &MultiService{client: client}
	client.Stream = &StreamService{client: client}
	client.Subreddit = &SubredditService{client: client}
	client.User = &UserService{client: client}
	client.Wiki = &WikiService{client: client}

	postAndCommentService := &postAndCommentService{client: client}
	client.Comment = &CommentService{client: client, postAndCommentService: postAndCommentService}
	client.Post = &PostService{client: client, postAndCommentService: postAndCommentService}

	return client
}

// NewClient returns a new Reddit API client.
// Use an Opt to configure the client credentials, such as WithCredentials or FromEnv.
func NewClient(opts ...Opt) (*Client, error) {
	client := newClient()

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	if client.client == nil {
		client.client = &http.Client{}
	}

	userAgentTransport := &userAgentTransport{
		userAgent: client.UserAgent(),
		Base:      client.client.Transport,
	}
	client.client.Transport = userAgentTransport

	if client.client.CheckRedirect == nil {
		client.client.CheckRedirect = client.redirect
	}

	oauthTransport := oauthTransport(client)
	client.client.Transport = oauthTransport

	return client, nil
}

// NewReadonlyClient returns a new read-only Reddit API client.
// The client will have limited access to the Reddit API.
// Options that modify credentials (such as WithCredentials or FromEnv) won't have any effect on this client.
func NewReadonlyClient(opts ...Opt) (*Client, error) {
	client := newClient()
	client.BaseURL, _ = url.Parse(defaultBaseURLReadonly)

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	if client.client == nil {
		client.client = &http.Client{}
	}

	userAgentTransport := &userAgentTransport{
		userAgent: client.UserAgent(),
		Base:      client.client.Transport,
	}
	client.client.Transport = userAgentTransport

	return client, nil
}

// todo...
// Some endpoints (notably the ones to get random subreddits/posts) redirect to a
// reddit.com url, which returns a 403 Forbidden for some reason, unless the url's
// host is changed to oauth.reddit.com
func (c *Client) redirect(req *http.Request, via []*http.Request) error {
	redirectURL := req.URL.String()
	redirectURL = strings.Replace(redirectURL, "https://www.reddit.com", defaultBaseURL, 1)

	reqURL, err := url.Parse(redirectURL)
	if err != nil {
		return err
	}
	req.URL = reqURL

	return nil
}

// The readonly Reddit url needs .json at the end of its path to return responses in JSON instead of HTML.
func (c *Client) appendJSONExtensionToRequestURLPath(req *http.Request) {
	readonlyURL, err := url.Parse(defaultBaseURLReadonly)
	if err != nil {
		return
	}

	if req.URL.Host != readonlyURL.Host {
		return
	}

	req.URL.Path += ".json"
}

// UserAgent returns the client's user agent.
func (c *Client) UserAgent() string {
	if c.userAgent == "" {
		userAgent := fmt.Sprintf("golang:%s:v%s", libraryName, libraryVersion)
		if c.Username != "" {
			userAgent += fmt.Sprintf(" (by /u/%s)", c.Username)
		}
		c.userAgent = userAgent
	}
	return c.userAgent
}

// NewRequest creates an API request with form data as the body.
// The path is the relative URL which will be resolves to the BaseURL of the Client.
// It should always be specified without a preceding slash.
func (c *Client) NewRequest(method string, path string, form url.Values) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	c.appendJSONExtensionToRequestURLPath(req)
	req.Header.Add(headerContentType, mediaTypeForm)
	req.Header.Add(headerAccept, mediaTypeJSON)

	return req, nil
}

// NewJSONRequest creates an API request with a JSON body.
// The path is the relative URL which will be resolved to the BaseURL of the Client.
// It should always be specified without a preceding slash.
func (c *Client) NewJSONRequest(method string, path string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	reqBody := bytes.NewReader(buf.Bytes())
	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	c.appendJSONExtensionToRequestURLPath(req)
	req.Header.Add(headerContentType, mediaTypeJSON)
	req.Header.Add(headerAccept, mediaTypeJSON)

	return req, nil
}

// Response is a Reddit response. This wraps the standard http.Response returned from Reddit.
type Response struct {
	*http.Response

	// Pagination anchor indicating there are more results after this id.
	After string
	// Pagination anchor indicating there are more results before this id.
	// todo: not sure yet if responses ever contain this
	Before string
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	return &response
}

func (r *Response) populateAnchors(a anchor) {
	r.After = a.After()
	r.Before = a.Before()
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, response.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(response.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}

		if anchor, ok := v.(anchor); ok {
			response.populateAnchors(anchor)
		}
	}

	return response, nil
}

// id returns the client's Reddit ID.
func (c *Client) id(ctx context.Context) (string, *Response, error) {
	if c.redditID != "" {
		return c.redditID, nil, nil
	}

	self, resp, err := c.User.Get(ctx, c.Username)
	if err != nil {
		return "", resp, err
	}

	c.redditID = fmt.Sprintf("%s_%s", kindUser, self.ID)
	return c.redditID, resp, nil
}

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
// Reddit also sometimes sends errors with 200 codes; we check for those too.
func CheckResponse(r *http.Response) error {
	jsonErrorResponse := &JSONErrorResponse{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, jsonErrorResponse)
		if len(jsonErrorResponse.JSON.Errors) > 0 {
			return jsonErrorResponse
		}
	}

	// reset response body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err = ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

func (c *Client) getListing(ctx context.Context, path string, opts interface{}) (*listing, *Response, error) {
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(thing)
	resp, err := c.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	listing, _ := root.Listing()
	return listing, resp, nil
}

// ListOptions specifies the optional parameters to various API calls that return a listing.
type ListOptions struct {
	// Maximum number of items to be returned.
	// Generally, the default is 25 and max is 100.
	Limit int `url:"limit,omitempty"`

	// The full ID of an item in the listing to use
	// as the anchor point of the list. Only items
	// appearing after it will be returned.
	After string `url:"after,omitempty"`

	// The full ID of an item in the listing to use
	// as the anchor point of the list. Only items
	// appearing before it will be returned.
	Before string `url:"before,omitempty"`
}

// ListSubredditOptions defines possible options used when searching for subreddits.
type ListSubredditOptions struct {
	ListOptions
	// One of: relevance, activity.
	Sort string `url:"sort,omitempty"`
}

// ListPostOptions defines possible options used when getting posts from a subreddit.
type ListPostOptions struct {
	ListOptions
	// One of: hour, day, week, month, year, all.
	Time string `url:"t,omitempty"`
}

// ListPostSearchOptions defines possible options used when searching for posts within a subreddit.
type ListPostSearchOptions struct {
	ListPostOptions
	// One of: relevance, hot, top, new, comments.
	Sort string `url:"sort,omitempty"`
}

// ListUserOverviewOptions defines possible options used when getting a user's post and/or comments.
type ListUserOverviewOptions struct {
	ListOptions
	// One of: hot, new, top, controversial.
	Sort string `url:"sort,omitempty"`
	// One of: hour, day, week, month, year, all.
	Time string `url:"t,omitempty"`
}

// ListDuplicatePostOptions defines possible options used when getting duplicates of a post, i.e.
// other submissions of the same URL.
type ListDuplicatePostOptions struct {
	ListOptions
	// If empty, it'll search for duplicates in all subreddits.
	Subreddit string `url:"sr,omitempty"`
	// One of: num_comments, new.
	Sort string `url:"sort,omitempty"`
	// If true, the search will only return duplicates that are
	// crossposts of the original post.
	CrosspostsOnly bool `url:"crossposts_only,omitempty"`
}

// ListModActionOptions defines possible options used when getting moderation actions in a subreddit.
type ListModActionOptions struct {
	// The max for the limit parameter here is 500.
	ListOptions
	// If empty, the search will return all action types.
	// One of: banuser, unbanuser, spamlink, removelink, approvelink, spamcomment, removecomment,
	// approvecomment, addmoderator, showcomment, invitemoderator, uninvitemoderator, acceptmoderatorinvite,
	// removemoderator, addcontributor, removecontributor, editsettings, editflair, distinguish, marknsfw,
	// wikibanned, wikicontributor, wikiunbanned, wikipagelisted, removewikicontributor, wikirevise,
	// wikipermlevel, ignorereports, unignorereports, setpermissions, setsuggestedsort, sticky, unsticky,
	// setcontestmode, unsetcontestmode, lock, unlock, muteuser, unmuteuser, createrule, editrule,
	// reorderrules, deleterule, spoiler, unspoiler, modmail_enrollment, community_styling, community_widgets,
	// markoriginalcontent, collections, events, hidden_award, add_community_topics, remove_community_topics,
	// create_scheduled_post, edit_scheduled_post, delete_scheduled_post, submit_scheduled_post,
	// edit_post_requirements, invitesubscriber, submit_content_rating_survey.
	Type string `url:"type,omitempty"`
	// If provided, only return the actions of this moderator.
	Moderator string `url:"mod,omitempty"`
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}
