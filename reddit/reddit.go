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
	libraryVersion = "0.0.1"

	defaultBaseURL  = "https://oauth.reddit.com"
	defaultTokenURL = "https://www.reddit.com/api/v1/access_token"

	mediaTypeJSON = "application/json"
	mediaTypeForm = "application/x-www-form-urlencoded"

	headerContentType = "Content-Type"
	headerAccept      = "Accept"
	headerUserAgent   = "User-Agent"
)

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map,
// since we'll only be modify the headers.
// Per the specification of http.RoundTripper, we should not directly modify a request.
func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

// Sets the User-Agent header for requests.
type userAgentTransport struct {
	userAgent string
	Base      http.RoundTripper
}

func (t *userAgentTransport) setUserAgent(req *http.Request) *http.Request {
	req2 := cloneRequest(req)
	req2.Header.Set(headerUserAgent, t.userAgent)
	return req2
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := t.setUserAgent(req)
	return t.base().RoundTrip(req2)
}

func (t *userAgentTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// RequestCompletionCallback defines the type of the request callback function.
type RequestCompletionCallback func(*http.Request, *http.Response)

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

	oauth2Transport *oauth2.Transport

	onRequestCompleted RequestCompletionCallback
}

// OnRequestCompleted sets the client's request completion callback.
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

func newClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// todo...
	// Some endpoints (notably the ones to get random subreddits/posts) redirect to a
	// reddit.com url, which returns a 403 Forbidden for some reason, unless the url's
	// host is changed to oauth.reddit.com
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectURL := req.URL.String()
		redirectURL = strings.Replace(redirectURL, "https://www.reddit.com", defaultBaseURL, 1)

		reqURL, err := url.Parse(redirectURL)
		if err != nil {
			return err
		}
		req.URL = reqURL

		return nil
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	tokenURL, _ := url.Parse(defaultTokenURL)

	c := &Client{client: httpClient, BaseURL: baseURL, TokenURL: tokenURL}

	c.Account = &AccountService{client: c}
	c.Collection = &CollectionService{client: c}
	c.Emoji = &EmojiService{client: c}
	c.Flair = &FlairService{client: c}
	c.Gold = &GoldService{client: c}
	c.Listings = &ListingsService{client: c}
	c.Message = &MessageService{client: c}
	c.Moderation = &ModerationService{client: c}
	c.Multi = &MultiService{client: c}
	c.Stream = &StreamService{client: c}
	c.Subreddit = &SubredditService{client: c}
	c.User = &UserService{client: c}

	postAndCommentService := &postAndCommentService{client: c}
	c.Comment = &CommentService{client: c, postAndCommentService: postAndCommentService}
	c.Post = &PostService{client: c, postAndCommentService: postAndCommentService}

	return c
}

// NewClient returns a client that can make requests to the Reddit API.
func NewClient(httpClient *http.Client, opts ...Opt) (c *Client, err error) {
	c = newClient(httpClient)

	for _, opt := range opts {
		if err = opt(c); err != nil {
			return
		}
	}

	oauthTransport := oauthTransport(c)
	c.client.Transport = oauthTransport

	return
}

// UserAgent returns the client's user agent.
func (c *Client) UserAgent() string {
	if c.userAgent == "" {
		c.userAgent = fmt.Sprintf("golang:%s:v%s (by /u/%s)", libraryName, libraryVersion, c.Username)
	}
	return c.userAgent
}

// NewRequest creates an API request.
// The path is the relative URL which will be resolves to the BaseURL of the Client.
// It should always be specified without a preceding slash.
func (c *Client) NewRequest(method string, path string, body interface{}) (*http.Request, error) {
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

	req.Header.Add(headerContentType, mediaTypeJSON)
	req.Header.Add(headerAccept, mediaTypeJSON)

	return req, nil
}

// NewRequestWithForm creates an API request with form data.
// The path is the relative URL which will be resolves to the BaseURL of the Client.
// It should always be specified without a preceding slash.
func (c *Client) NewRequestWithForm(method string, path string, form url.Values) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add(headerContentType, mediaTypeForm)
	req.Header.Add(headerAccept, mediaTypeJSON)

	return req, nil
}

// Response is a PlayNetwork response. This wraps the standard http.Response returned from PlayNetwork.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	return &response
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

	c.redditID = fmt.Sprintf("%s_%s", kindAccount, self.ID)
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
