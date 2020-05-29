package geddit

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
	libraryName    = "github.com/vartanbeno/geddit"
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

// sets the User-Agent header for requests
type userAgentTransport struct {
	ua   string
	Base http.RoundTripper
}

func (t *userAgentTransport) setUserAgent(req *http.Request) *http.Request {
	req2 := cloneRequest(req)
	req2.Header.Set(headerUserAgent, t.ua)
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

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// Client manages communication with the Reddit API
type Client struct {
	// HTTP client used to communicate with the Reddit API
	client *http.Client

	BaseURL  *url.URL
	TokenURL *url.URL

	userAgent string

	ID       string
	Secret   string
	Username string
	Password string

	Comment   CommentService
	Flair     FlairService
	Link      LinkService
	Listings  ListingsService
	Search    SearchService
	Subreddit SubredditService
	User      UserService
	Vote      VoteService

	oauth2Transport *oauth2.Transport

	onRequestCompleted RequestCompletionCallback
}

// OnRequestCompleted sets the client's request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

func newClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	tokenURL, _ := url.Parse(defaultTokenURL)

	c := &Client{client: httpClient, BaseURL: baseURL, TokenURL: tokenURL}

	c.Comment = &CommentServiceOp{client: c}
	c.Flair = &FlairServiceOp{client: c}
	c.Link = &LinkServiceOp{client: c}
	c.Listings = &ListingsServiceOp{client: c}
	c.Search = &SearchServiceOp{client: c}
	c.Subreddit = &SubredditServiceOp{client: c}
	c.User = &UserServiceOp{client: c}
	c.Vote = &VoteServiceOp{client: c}

	return c
}

// NewClient returns a client that can make requests to the Reddit API
func NewClient(httpClient *http.Client, opts ...Opt) (c *Client, err error) {
	c = newClient(httpClient)

	for _, opt := range opts {
		if err = opt(c); err != nil {
			return
		}
	}

	c.userAgent = fmt.Sprintf("golang:%s:v%s (by /u/%s)", libraryName, libraryVersion, c.Username)
	userAgentTransport := &userAgentTransport{
		ua:   c.userAgent,
		Base: c.client.Transport,
	}

	oauth2Config := oauth2Config{
		id:                 c.ID,
		secret:             c.Secret,
		username:           c.Username,
		password:           c.Password,
		tokenURL:           c.TokenURL.String(),
		userAgentTransport: userAgentTransport,
	}
	c.client.Transport = oauth2Transport(oauth2Config)

	return
}

// NewRequest creates an API request
// The path is the relative URL which will be resolves to the BaseURL of the Client
// It should always be specified without a preceding slash
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
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

// NewPostForm creates an API request with a POST form
// The path is the relative URL which will be resolves to the BaseURL of the Client
// It should always be specified without a preceding slash
func (c *Client) NewPostForm(path string, form url.Values) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(form.Encode()))
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

// newResponse creates a new Response for the provided http.Response
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

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	response := newResponse(resp)
	defer func() {
		if rerr := response.Body.Close(); err == nil {
			err = rerr
		}
	}()

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

	return response, err
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
		if jsonErrorResponse.JSON != nil && len(jsonErrorResponse.JSON.Errors) > 0 {
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

// ListOptions are the optional parameters to the various endpoints that return lists
type ListOptions struct {
	// For getting submissions
	// all, year, month, week, day, hour
	Timespan string `url:"t,omitempty"`

	// Common for all listing endpoints
	After  string `url:"after,omitempty"`
	Before string `url:"before,omitempty"`
	Limit  int    `url:"limit,omitempty"` // default: 25, max: 100
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

func addQuery(url string, query url.Values) string {
	if query == nil || len(query) == 0 {
		return url
	}
	return url + "?" + query.Encode()
}
