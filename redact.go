package jsonutil

import (
	"encoding/json"
	"strings"
)

var DefaultSensitiveKeys = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"access_token",
	"refresh_token",
	"api_key",
	"apikey",
	"authorization",
	"cookie",
	"session",
	"ssn",
}

func Redact(v any, extraKeys ...string) any {
	keys := make(map[string]struct{}, len(DefaultSensitiveKeys)+len(extraKeys))

	for _, key := range DefaultSensitiveKeys {
		keys[strings.ToLower(key)] = struct{}{}
	}

	for _, key := range extraKeys {
		keys[strings.ToLower(key)] = struct{}{}
	}

	return redactValue(v, keys)
}

func redactValue(v any, keys map[string]struct{}) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))

		for k, val := range x {
			if _, ok := keys[strings.ToLower(k)]; ok {
				out[k] = "[REDACTED]"
			} else {
				out[k] = redactValue(val, keys)
			}
		}

		return out

	case []any:
		out := make([]any, len(x))
		for i, val := range x {
			out[i] = redactValue(val, keys)
		}
		return out

	default:
		return v
	}
}

func SafePretty(v any) (string, error) {
	b, err := json.Marshal(Redact(v))
	if err != nil {
		return "", err
	}

	return PrettyJSON(b)
}