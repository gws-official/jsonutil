package jsonutil

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func EqualJSON(a, b []byte) bool {
	var av any
	var bv any

	if json.Unmarshal(a, &av) != nil {
		return false
	}

	if json.Unmarshal(b, &bv) != nil {
		return false
	}

	return reflect.DeepEqual(av, bv)
}

func NormalizeJSON(data []byte) (string, error) {
	var v any

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	if err := dec.Decode(&v); err != nil {
		return "", err
	}

	out, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(out), nil
}