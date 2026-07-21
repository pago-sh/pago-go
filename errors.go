package pago

import (
	"fmt"
	"net/http"
)

// Error is implemented by every error returned by this SDK, which makes it
// possible to tell an SDK failure apart from any other error with errors.As.
type Error interface {
	error
	isPagoError()
}

type baseError struct{}

func (baseError) isPagoError() {}

// EncodingError is returned when a request body cannot be encoded, or when a
// successful response body cannot be decoded into the expected type.
type EncodingError struct {
	baseError
	Op  string
	Err error
}

func (e *EncodingError) Error() string {
	return fmt.Sprintf("pago: failed to %s: %v", e.Op, e.Err)
}

func (e *EncodingError) Unwrap() error { return e.Err }

// NetworkError is returned when the request never produced an HTTP response.
type NetworkError struct {
	baseError
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("pago: API network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error { return e.Err }

// ServerError is returned for any 5xx response.
type ServerError struct {
	baseError
	StatusCode int
	Body       []byte
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("pago: API returned a server error: %d - %s", e.StatusCode, e.Body)
}

// APIError is returned for any 4xx response that has no generated error type.
// Every generated error type embeds it, so errors.As with a *APIError target
// matches those too.
type APIError struct {
	baseError
	StatusCode int
	Body       []byte
	// Header holds the response headers, useful for request tracing.
	Header http.Header
}

func (e *APIError) Error() string {
	return fmt.Sprintf("pago: API returned an error: %d - %s", e.StatusCode, e.Body)
}

// RateLimitError is returned for a 429 response.
type RateLimitError struct {
	APIError
	// RetryAfter is the value of the Retry-After header in seconds, or -1 when
	// the header is absent or malformed.
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("pago: API rate limit exceeded: %d - retry after %ds", e.StatusCode, e.RetryAfter)
}

func (e *RateLimitError) Unwrap() error { return &e.APIError }
