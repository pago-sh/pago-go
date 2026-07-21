package pago

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
)

// AddQueryValue appends value to the query string under key.
//
// It is the single encoding point used by every generated parameter struct:
//
//   - nil values and nil pointers are skipped, which is how an unset optional
//     parameter stays out of the query string;
//   - pointers and interfaces are followed;
//   - slices repeat the key once per element;
//   - maps use the deepObject form key[subkey]=value, sorted by key so the
//     query string is deterministic;
//   - anything else is rendered with Stringify.
func AddQueryValue(query url.Values, key string, value any) {
	if value == nil {
		return
	}
	addQueryValue(query, key, reflect.ValueOf(value))
}

func addQueryValue(query url.Values, key string, value reflect.Value) {
	switch value.Kind() {
	case reflect.Invalid:
		return
	case reflect.Pointer, reflect.Interface:
		if value.IsNil() {
			return
		}
		addQueryValue(query, key, value.Elem())
	case reflect.Slice, reflect.Array:
		if value.Kind() == reflect.Slice && value.IsNil() {
			return
		}
		for index := 0; index < value.Len(); index++ {
			addQueryValue(query, key, value.Index(index))
		}
	case reflect.Map:
		if value.IsNil() {
			return
		}
		keys := value.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return stringifyValue(keys[i]) < stringifyValue(keys[j])
		})
		for _, mapKey := range keys {
			addQueryValue(query, key+"["+stringifyValue(mapKey)+"]", value.MapIndex(mapKey))
		}
	default:
		// Generated union types carry their value behind a custom marshaller.
		// Encoding the raw struct would quote or brace it, so the JSON is
		// decoded back into a plain value and re-dispatched.
		if decoded, ok := decodeMarshaler(value); ok {
			if decoded == nil {
				return
			}
			addQueryValue(query, key, reflect.ValueOf(decoded))
			return
		}
		query.Add(key, stringifyValue(value))
	}
}

func decodeMarshaler(value reflect.Value) (any, bool) {
	if !value.CanInterface() {
		return nil, false
	}
	marshaler, ok := value.Interface().(json.Marshaler)
	if !ok {
		return nil, false
	}
	encoded, err := marshaler.MarshalJSON()
	if err != nil {
		return nil, false
	}
	var decoded any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, false
	}
	return decoded, true
}

// Stringify renders a path or query parameter value as a string.
func Stringify(value any) string {
	if value == nil {
		return ""
	}
	return stringifyValue(reflect.ValueOf(value))
}

func stringifyValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Invalid:
		return ""
	case reflect.Pointer, reflect.Interface:
		if value.IsNil() {
			return ""
		}
		return stringifyValue(value.Elem())
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	default:
		if !value.CanInterface() {
			return ""
		}
		if encoded, err := json.Marshal(value.Interface()); err == nil {
			return string(encoded)
		}
		return fmt.Sprintf("%v", value.Interface())
	}
}
