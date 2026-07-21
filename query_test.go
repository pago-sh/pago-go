package pago

import (
	"net/url"
	"testing"
)

type queryEnum string

func TestAddQueryValueSkipsUnsetValues(t *testing.T) {
	query := url.Values{}
	var missing *int
	var emptySlice []string
	AddQueryValue(query, "a", nil)
	AddQueryValue(query, "b", missing)
	AddQueryValue(query, "c", emptySlice)

	if len(query) != 0 {
		t.Fatalf("expected an empty query, got %v", query)
	}
}

func TestAddQueryValueEncodesScalars(t *testing.T) {
	query := url.Values{}
	limit := 10
	ratio := 1.5
	enabled := true
	AddQueryValue(query, "limit", &limit)
	AddQueryValue(query, "ratio", ratio)
	AddQueryValue(query, "enabled", enabled)
	AddQueryValue(query, "status", queryEnum("active"))

	expected := "enabled=true&limit=10&ratio=1.5&status=active"
	if got := query.Encode(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestAddQueryValueRepeatsSliceElements(t *testing.T) {
	query := url.Values{}
	AddQueryValue(query, "id", []string{"a", "b"})

	if got := query["id"]; len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("expected the key to repeat, got %v", got)
	}
}

func TestAddQueryValueUsesDeepObjectForMaps(t *testing.T) {
	query := url.Values{}
	AddQueryValue(query, "metadata", map[string]any{"b": 2, "a": []string{"x", "y"}})

	expected := "metadata%5Ba%5D=x&metadata%5Ba%5D=y&metadata%5Bb%5D=2"
	if got := query.Encode(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

// unionValue stands in for a generated union type: its value lives behind a
// custom marshaller.
type unionValue struct {
	encoded string
}

func (u unionValue) MarshalJSON() ([]byte, error) { return []byte(u.encoded), nil }

func TestAddQueryValueUnwrapsCustomMarshalers(t *testing.T) {
	query := url.Values{}
	AddQueryValue(query, "text", unionValue{`"raw"`})
	AddQueryValue(query, "number", unionValue{`12`})
	AddQueryValue(query, "list", unionValue{`["a","b"]`})
	AddQueryValue(query, "empty", unionValue{`null`})

	expected := "list=a&list=b&number=12&text=raw"
	if got := query.Encode(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestAddQueryValueUnwrapsMarshalersInsideMaps(t *testing.T) {
	query := url.Values{}
	AddQueryValue(query, "metadata", map[string]unionValue{"a": {`"x"`}})

	expected := "metadata%5Ba%5D=x"
	if got := query.Encode(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestStringify(t *testing.T) {
	value := 7
	cases := []struct {
		input    any
		expected string
	}{
		{nil, ""},
		{"raw", "raw"},
		{7, "7"},
		{&value, "7"},
		{true, "true"},
		{queryEnum("active"), "active"},
	}
	for _, testCase := range cases {
		if got := Stringify(testCase.input); got != testCase.expected {
			t.Fatalf("Stringify(%v) = %q, expected %q", testCase.input, got, testCase.expected)
		}
	}
}
