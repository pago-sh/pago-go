package pago

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type testPayload struct {
	Name string `json:"name"`
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *ClientBase {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return NewClientBase("2026-04", WithBaseURL(server.URL), WithAccessToken("token"))
}

func TestDoSendsAuthenticatedRequest(t *testing.T) {
	var received *http.Request
	var body testPayload
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		received = r
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"pong"}`))
	})

	var out testPayload
	err := client.Do(context.Background(), Request{
		Method:       "POST",
		Path:         "/v1/things/{id}/",
		PathParams:   map[string]any{"id": "abc def"},
		Query:        url.Values{"limit": []string{"10"}},
		Body:         testPayload{Name: "ping"},
		ResponseType: ResponseJSON,
	}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Name != "pong" {
		t.Fatalf("expected decoded response, got %q", out.Name)
	}
	if body.Name != "ping" {
		t.Fatalf("expected encoded request body, got %q", body.Name)
	}
	if got := received.URL.Path; got != "/v1/things/abc def/" {
		t.Fatalf("expected the path parameter to be interpolated, got %q", got)
	}
	if got := received.URL.Query().Get("limit"); got != "10" {
		t.Fatalf("expected the query string to be sent, got %q", got)
	}
	if got := received.Header.Get("Authorization"); got != "Bearer token" {
		t.Fatalf("expected a bearer token, got %q", got)
	}
	if got := received.Header.Get("Pago-Version"); got != "2026-04" {
		t.Fatalf("expected the version header, got %q", got)
	}
	if got := received.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected a JSON content type, got %q", got)
	}
}

func TestDoWithoutBodyOmitsContentType(t *testing.T) {
	var received *http.Request
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		received = r
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Do(context.Background(), Request{
		Method:       "DELETE",
		Path:         "/v1/things/1/",
		ResponseType: ResponseNone,
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := received.Header.Get("Content-Type"); got != "" {
		t.Fatalf("expected no content type, got %q", got)
	}
}

func TestDoDecodesTextResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("plain body"))
	})

	var out string
	err := client.Do(context.Background(), Request{
		Method:       "GET",
		Path:         "/v1/text/",
		ResponseType: ResponseText,
	}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "plain body" {
		t.Fatalf("expected the raw body, got %q", out)
	}
}

type notFoundError struct {
	APIError
	Detail string
}

func (e *notFoundError) Error() string { return "not found: " + e.Detail }

func TestDoUsesGeneratedErrorDecoder(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"missing"}`))
	})

	err := client.Do(context.Background(), Request{
		Method:       "GET",
		Path:         "/v1/things/1/",
		ResponseType: ResponseJSON,
		Errors: map[int]ErrorDecoder{
			http.StatusNotFound: func(statusCode int, body []byte, header http.Header) error {
				var payload struct {
					Detail string `json:"detail"`
				}
				_ = json.Unmarshal(body, &payload)
				return &notFoundError{
					APIError: APIError{StatusCode: statusCode, Body: body, Header: header},
					Detail:   payload.Detail,
				}
			},
		},
	}, &testPayload{})

	var typed *notFoundError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *notFoundError, got %v", err)
	}
	if typed.Detail != "missing" {
		t.Fatalf("expected the decoded error body, got %q", typed.Detail)
	}
	if typed.StatusCode != http.StatusNotFound {
		t.Fatalf("expected the status code, got %d", typed.StatusCode)
	}
}

func TestDoReturnsAPIErrorWithoutDecoder(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	})

	err := client.Do(context.Background(), Request{
		Method:       "GET",
		Path:         "/v1/things/",
		ResponseType: ResponseJSON,
	}, &testPayload{})

	var typed *APIError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *APIError, got %v", err)
	}
	if typed.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", typed.StatusCode)
	}
}

func TestDoReturnsRateLimitError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "42")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	err := client.Do(context.Background(), Request{
		Method:       "GET",
		Path:         "/v1/things/",
		ResponseType: ResponseJSON,
	}, &testPayload{})

	var typed *RateLimitError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *RateLimitError, got %v", err)
	}
	if typed.RetryAfter != 42 {
		t.Fatalf("expected Retry-After to be parsed, got %d", typed.RetryAfter)
	}

	var apiError *APIError
	if !errors.As(err, &apiError) {
		t.Fatal("expected a rate limit error to unwrap to *APIError")
	}
}

func TestDoReturnsServerError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	})

	err := client.Do(context.Background(), Request{
		Method:       "GET",
		Path:         "/v1/things/",
		ResponseType: ResponseJSON,
	}, &testPayload{})

	var typed *ServerError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *ServerError, got %v", err)
	}
	if string(typed.Body) != "boom" {
		t.Fatalf("expected the raw body, got %q", typed.Body)
	}
}

func TestDoReturnsNetworkError(t *testing.T) {
	client := NewClientBase("2026-04", WithBaseURL("http://127.0.0.1:1"))

	err := client.Do(context.Background(), Request{
		Method:       "GET",
		Path:         "/v1/things/",
		ResponseType: ResponseJSON,
	}, &testPayload{})

	var typed *NetworkError
	if !errors.As(err, &typed) {
		t.Fatalf("expected a *NetworkError, got %v", err)
	}
}

func TestErrorsImplementTheSDKErrorInterface(t *testing.T) {
	values := []error{
		&NetworkError{},
		&ServerError{},
		&APIError{},
		&RateLimitError{},
		&EncodingError{},
		&WebhookError{},
		&WebhookVerificationError{},
		&WebhookUnknownTypeError{},
	}
	for _, err := range values {
		var sdkError Error
		if !errors.As(err, &sdkError) {
			t.Fatalf("%T does not implement pago.Error", err)
		}
	}
}

func TestExpandPathLeavesUnknownPlaceholders(t *testing.T) {
	if got := expandPath("/v1/{a}/{b}", map[string]any{"a": 1}); got != "/v1/1/{b}" {
		t.Fatalf("unexpected path: %q", got)
	}
}
