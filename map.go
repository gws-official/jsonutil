package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrMissingKey = errors.New("jsonutil: missing key")
	ErrWrongType  = errors.New("jsonutil: wrong type")
)

func GetString(m map[string]any, key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	return s, ok
}

func GetBool(m map[string]any, key string) (bool, bool) {
	v, ok := m[key]
	if !ok {
		return false, false
	}

	b, ok := v.(bool)
	return b, ok
}

func GetInt(m map[string]any, key string) (int, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}

	switch x := v.(type) {
	case int:
		return x, true
	case int64:
		return int(x), true
	case float64:
		if x == float64(int(x)) {
			return int(x), true
		}
	case json.Number:
		i, err := x.Int64()
		if err == nil {
			return int(i), true
		}
	case string:
		i, err := strconv.Atoi(x)
		if err == nil {
			return i, true
		}
	}

	return 0, false
}

func GetFloat64(m map[string]any, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}

	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case json.Number:
		f, err := x.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	}

	return 0, false
}

func GetMap(m map[string]any, key string) (map[string]any, bool) {
	v, ok := m[key]
	if !ok {
		return nil, false
	}

	out, ok := v.(map[string]any)
	return out, ok
}

func GetSlice(m map[string]any, key string) ([]any, bool) {
	v, ok := m[key]
	if !ok {
		return nil, false
	}

	out, ok := v.([]any)
	return out, ok
}

func StringOr(m map[string]any, key string, fallback string) string {
	if v, ok := GetString(m, key); ok {
		return v
	}
	return fallback
}

func BoolOr(m map[string]any, key string, fallback bool) bool {
	if v, ok := GetBool(m, key); ok {
		return v
	}
	return fallback
}

func IntOr(m map[string]any, key string, fallback int) int {
	if v, ok := GetInt(m, key); ok {
		return v
	}
	return fallback
}

func RequireString(m map[string]any, key string) (string, error) {
	v, exists := m[key]
	if !exists {
		return "", fmt.Errorf("%w: %s", ErrMissingKey, key)
	}

	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%w: %s is %T, not string", ErrWrongType, key, v)
	}

	return s, nil
}