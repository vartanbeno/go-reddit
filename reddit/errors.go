package reddit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is an error coming from Reddit.
type APIError struct {
	Label  string
	Reason string
	Field  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("field %q caused %s: %s", e.Field, e.Label, e.Reason)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *APIError) UnmarshalJSON(data []byte) error {
	var info [3]string

	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}

	e.Label = info[0]
	e.Reason = info[1]
	e.Field = info[2]

	return nil
}

// JSONErrorResponse is an error response that sometimes gets returned with a 200 code.
type JSONErrorResponse struct {
	// HTTP response that caused this error.
	Response *http.Response `json:"-"`

	JSON struct {
		Errors []APIError `json:"errors,omitempty"`
	} `json:"json"`
}

func (r *JSONErrorResponse) Error() string {
	var message string
	if len(r.JSON.Errors) > 0 {
		message = r.JSON.Errors[0].Error()
	}

	return fmt.Sprintf(
		"%s %s: %d %s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, message,
	)
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response `json:"-"`

	// Error message
	Message string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"%s %s: %d %s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message,
	)
}

// todo: rate limit errors
