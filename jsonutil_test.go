package jsonutil

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestDecodeStrict(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	got, err := DecodeStrict[payload](strings.NewReader(`{"name":"alice"}`))
	if err != nil {
		t.Fatalf("DecodeStrict returned error: %v", err)
	}

	if got.Name != "alice" {
		t.Fatalf("got name %q, want %q", got.Name, "alice")
	}

	_, err = DecodeStrict[payload](strings.NewReader(`{"name":"alice","extra":true}`))
	if err == nil {
		t.Fatal("expected unknown field error")
	}

	_, err = DecodeStrict[payload](strings.NewReader(`{"name":"alice"} {"name":"bob"}`))
	if err == nil {
		t.Fatal("expected trailing data error")
	}
}

func TestDecodeUseNumber(t *testing.T) {
	got, err := DecodeUseNumber([]byte(`{"count":123,"name":"alice"}`))
	if err != nil {
		t.Fatalf("DecodeUseNumber returned error: %v", err)
	}

	n, ok := got["count"].(json.Number)
	if !ok {
		t.Fatalf("count type = %T, want json.Number", got["count"])
	}

	if n.String() != "123" {
		t.Fatalf("count = %q, want %q", n.String(), "123")
	}
}

func TestPrettyCompactAndValidJSON(t *testing.T) {
	raw := []byte(`{"name":"alice","count":1}`)

	if !ValidJSON(raw) {
		t.Fatal("ValidJSON returned false for valid JSON")
	}

	if ValidJSON([]byte(`{"name":`)) {
		t.Fatal("ValidJSON returned true for invalid JSON")
	}

	pretty, err := PrettyJSON(raw)
	if err != nil {
		t.Fatalf("PrettyJSON returned error: %v", err)
	}

	if !strings.Contains(pretty, "\n") {
		t.Fatalf("PrettyJSON result does not appear formatted: %q", pretty)
	}

	compact, err := CompactJSON([]byte(pretty))
	if err != nil {
		t.Fatalf("CompactJSON returned error: %v", err)
	}

	if compact != string(raw) {
		t.Fatalf("compact = %q, want %q", compact, raw)
	}
}

func TestMapGetters(t *testing.T) {
	m := map[string]any{
		"name":   "alice",
		"active": true,
		"age":    float64(30),
		"score":  json.Number("99.5"),
		"meta":   map[string]any{"role": "admin"},
		"items":  []any{"a", "b"},
	}

	if got, ok := GetString(m, "name"); !ok || got != "alice" {
		t.Fatalf("GetString = %q, %v; want alice, true", got, ok)
	}

	if got, ok := GetBool(m, "active"); !ok || got != true {
		t.Fatalf("GetBool = %v, %v; want true, true", got, ok)
	}

	if got, ok := GetInt(m, "age"); !ok || got != 30 {
		t.Fatalf("GetInt = %d, %v; want 30, true", got, ok)
	}

	if got, ok := GetFloat64(m, "score"); !ok || got != 99.5 {
		t.Fatalf("GetFloat64 = %v, %v; want 99.5, true", got, ok)
	}

	if got, ok := GetMap(m, "meta"); !ok || got["role"] != "admin" {
		t.Fatalf("GetMap = %#v, %v; want role admin, true", got, ok)
	}

	if got, ok := GetSlice(m, "items"); !ok || len(got) != 2 {
		t.Fatalf("GetSlice = %#v, %v; want len 2, true", got, ok)
	}
}

func TestMapFallbacks(t *testing.T) {
	m := map[string]any{
		"name":   "alice",
		"active": true,
		"age":    "30",
	}

	if got := StringOr(m, "name", "fallback"); got != "alice" {
		t.Fatalf("StringOr existing = %q, want alice", got)
	}

	if got := StringOr(m, "missing", "fallback"); got != "fallback" {
		t.Fatalf("StringOr fallback = %q, want fallback", got)
	}

	if got := BoolOr(m, "active", false); got != true {
		t.Fatalf("BoolOr existing = %v, want true", got)
	}

	if got := BoolOr(m, "missing", true); got != true {
		t.Fatalf("BoolOr fallback = %v, want true", got)
	}

	if got := IntOr(m, "age", 0); got != 30 {
		t.Fatalf("IntOr existing = %d, want 30", got)
	}

	if got := IntOr(m, "missing", 7); got != 7 {
		t.Fatalf("IntOr fallback = %d, want 7", got)
	}
}

func TestRequireString(t *testing.T) {
	m := map[string]any{
		"name": "alice",
		"age":  30,
	}

	got, err := RequireString(m, "name")
	if err != nil {
		t.Fatalf("RequireString returned error: %v", err)
	}

	if got != "alice" {
		t.Fatalf("RequireString = %q, want alice", got)
	}

	_, err = RequireString(m, "missing")
	if !errors.Is(err, ErrMissingKey) {
		t.Fatalf("missing key error = %v, want ErrMissingKey", err)
	}

	_, err = RequireString(m, "age")
	if !errors.Is(err, ErrWrongType) {
		t.Fatalf("wrong type error = %v, want ErrWrongType", err)
	}
}

func TestRedact(t *testing.T) {
	input := map[string]any{
		"username": "alice",
		"password": "secret",
		"nested": map[string]any{
			"token": "abc123",
			"keep":  "value",
		},
		"items": []any{
			map[string]any{
				"api_key": "key",
				"name":    "service",
			},
		},
	}

	got := Redact(input).(map[string]any)

	if got["password"] != "[REDACTED]" {
		t.Fatalf("password was not redacted: %#v", got["password"])
	}

	nested := got["nested"].(map[string]any)
	if nested["token"] != "[REDACTED]" {
		t.Fatalf("nested token was not redacted: %#v", nested["token"])
	}

	if nested["keep"] != "value" {
		t.Fatalf("non-sensitive value changed: %#v", nested["keep"])
	}

	items := got["items"].([]any)
	first := items[0].(map[string]any)

	if first["api_key"] != "[REDACTED]" {
		t.Fatalf("slice item api_key was not redacted: %#v", first["api_key"])
	}

	if first["name"] != "service" {
		t.Fatalf("slice item name changed: %#v", first["name"])
	}
}

func TestRedactExtraKeysAndCaseInsensitive(t *testing.T) {
	input := map[string]any{
		"CustomSecret": "value",
	}

	got := Redact(input, "customsecret").(map[string]any)

	if got["CustomSecret"] != "[REDACTED]" {
		t.Fatalf("custom key was not redacted: %#v", got["CustomSecret"])
	}
}

func TestSafePretty(t *testing.T) {
	input := map[string]any{
		"password": "secret",
		"name":     "alice",
	}

	got, err := SafePretty(input)
	if err != nil {
		t.Fatalf("SafePretty returned error: %v", err)
	}

	if strings.Contains(got, "secret") {
		t.Fatalf("SafePretty leaked sensitive value: %s", got)
	}

	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("SafePretty did not include redaction marker: %s", got)
	}
}

func TestEqualJSON(t *testing.T) {
	a := []byte(`{"name":"alice","age":30}`)
	b := []byte(`{
		"age": 30,
		"name": "alice"
	}`)

	if !EqualJSON(a, b) {
		t.Fatal("EqualJSON returned false for equivalent JSON")
	}

	if EqualJSON(a, []byte(`{"name":"bob","age":30}`)) {
		t.Fatal("EqualJSON returned true for different JSON")
	}

	if EqualJSON(a, []byte(`{"name":`)) {
		t.Fatal("EqualJSON returned true for invalid JSON")
	}
}

func TestNormalizeJSON(t *testing.T) {
	got, err := NormalizeJSON([]byte(`{
		"count": 123,
		"name": "alice"
	}`))
	if err != nil {
		t.Fatalf("NormalizeJSON returned error: %v", err)
	}

	wantA := `{"count":123,"name":"alice"}`
	wantB := `{"name":"alice","count":123}`

	if got != wantA && got != wantB {
		t.Fatalf("NormalizeJSON = %q, want %q or %q", got, wantA, wantB)
	}
}
