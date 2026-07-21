// Package pago is the runtime shared by every generated Pago API version
// package. It carries the HTTP client, the error hierarchy and the webhook
// signature verification.
//
// Use it through a versioned client, for example:
//
//	client := v2026_04.New(pago.WithAccessToken(os.Getenv("PAGO_ACCESS_TOKEN")))
package pago

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DefaultBaseURL is the base URL of the Pago API.
const DefaultBaseURL = "https://api.pago.sh"

// ResponseType describes how the body of a successful response is decoded.
type ResponseType string

const (
	// ResponseJSON decodes the response body as JSON.
	ResponseJSON ResponseType = "json"
	// ResponseText decodes the response body as plain text.
	ResponseText ResponseType = "text"
	// ResponseNone discards the response body.
	ResponseNone ResponseType = "none"
)

// ErrorDecoder builds a typed error from an error response.
type ErrorDecoder func(statusCode int, body []byte, header http.Header) error

// Option configures a ClientBase.
type Option func(*ClientBase)

// WithBaseURL overrides the API base URL, which is useful for testing against
// a local server.
func WithBaseURL(baseURL string) Option {
	return func(c *ClientBase) { c.baseURL = strings.TrimSuffix(baseURL, "/") }
}

// WithAccessToken sets the bearer token sent with every request.
func WithAccessToken(accessToken string) Option {
	return func(c *ClientBase) { c.accessToken = accessToken }
}

// WithHTTPClient injects the http.Client used to perform requests, which is
// the hook for custom timeouts, proxies, retries or instrumentation.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *ClientBase) { c.httpClient = httpClient }
}

// WithVersion overrides the API version header sent with every request.
func WithVersion(version string) Option {
	return func(c *ClientBase) { c.version = version }
}

// WithHeader adds a header to every request.
func WithHeader(key, value string) Option {
	return func(c *ClientBase) { c.headers[key] = value }
}

// ClientBase performs the HTTP requests of a versioned client.
type ClientBase struct {
	baseURL     string
	accessToken string
	version     string
	httpClient  *http.Client
	headers     map[string]string
}

// NewClientBase builds the runtime client used by a versioned client.
func NewClientBase(version string, options ...Option) *ClientBase {
	c := &ClientBase{
		baseURL:    DefaultBaseURL,
		version:    version,
		httpClient: http.DefaultClient,
		headers:    map[string]string{},
	}
	for _, option := range options {
		option(c)
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}

// BaseURL returns the base URL requests are sent to.
func (c *ClientBase) BaseURL() string { return c.baseURL }

// Version returns the API version sent with every request.
func (c *ClientBase) Version() string { return c.version }

// HTTPClient returns the underlying http.Client.
func (c *ClientBase) HTTPClient() *http.Client { return c.httpClient }

// Request describes a single API call.
type Request struct {
	Method       string
	Path         string
	PathParams   map[string]any
	Query        url.Values
	Body         any
	ResponseType ResponseType
	// Errors maps an HTTP status code to the decoder of its documented error
	// body.
	Errors map[int]ErrorDecoder
}

// Do performs the request and decodes a successful response into out, which
// must be a non-nil pointer unless the response type is ResponseNone.
func (c *ClientBase) Do(ctx context.Context, request Request, out any) error {
	var payload io.Reader
	if request.Body != nil {
		encoded, err := json.Marshal(request.Body)
		if err != nil {
			return &EncodingError{Op: "encode request body", Err: err}
		}
		payload = bytes.NewReader(encoded)
	}

	target := c.baseURL + expandPath(request.Path, request.PathParams)
	if len(request.Query) > 0 {
		target += "?" + request.Query.Encode()
	}

	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, target, payload)
	if err != nil {
		return &NetworkError{Err: err}
	}
	for key, value := range c.headers {
		httpRequest.Header.Set(key, value)
	}
	httpRequest.Header.Set("Accept", "application/json")
	if request.Body != nil {
		httpRequest.Header.Set("Content-Type", "application/json")
	}
	if c.version != "" {
		httpRequest.Header.Set("Pago-Version", c.version)
	}
	if c.accessToken != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return &NetworkError{Err: err}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &NetworkError{Err: err}
	}

	if err := c.checkResponse(response, body, request.Errors); err != nil {
		return err
	}

	if out == nil {
		return nil
	}
	switch request.ResponseType {
	case ResponseJSON:
		if len(bytes.TrimSpace(body)) == 0 {
			return nil
		}
		if err := json.Unmarshal(body, out); err != nil {
			return &EncodingError{Op: "decode response body", Err: err}
		}
	case ResponseText:
		text, ok := out.(*string)
		if !ok {
			return &EncodingError{
				Op:  "decode response body",
				Err: fmt.Errorf("expected *string, got %T", out),
			}
		}
		*text = string(body)
	}
	return nil
}

func (c *ClientBase) checkResponse(
	response *http.Response,
	body []byte,
	decoders map[int]ErrorDecoder,
) error {
	status := response.StatusCode
	switch {
	case status >= 500:
		return &ServerError{StatusCode: status, Body: body}
	case status == http.StatusTooManyRequests:
		retryAfter := -1
		if value, err := strconv.Atoi(response.Header.Get("Retry-After")); err == nil {
			retryAfter = value
		}
		return &RateLimitError{
			APIError:   APIError{StatusCode: status, Body: body, Header: response.Header},
			RetryAfter: retryAfter,
		}
	case status >= 400:
		if decoder, ok := decoders[status]; ok {
			return decoder(status, body, response.Header)
		}
		return &APIError{StatusCode: status, Body: body, Header: response.Header}
	}
	return nil
}

// expandPath substitutes the {name} placeholders of a path template.
func expandPath(path string, params map[string]any) string {
	if len(params) == 0 {
		return path
	}
	for name, value := range params {
		path = strings.ReplaceAll(path, "{"+name+"}", url.PathEscape(Stringify(value)))
	}
	return path
}
