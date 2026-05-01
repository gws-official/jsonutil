package jsonutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

func DecodeStrict[T any](r io.Reader) (T, error) {
	var out T

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&out); err != nil {
		return out, err
	}

	if dec.Decode(&struct{}{}) != io.EOF {
		return out, errors.New("jsonutil: trailing data after JSON value")
	}

	return out, nil
}

func DecodeUseNumber(data []byte) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var out map[string]any
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

func PrettyJSON(data []byte) (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return "", err
	}
	return out.String(), nil
}

func CompactJSON(data []byte) (string, error) {
	var out bytes.Buffer
	if err := json.Compact(&out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}

func ValidJSON(data []byte) bool {
	return json.Valid(data)
}