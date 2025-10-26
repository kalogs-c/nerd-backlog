package httpjson

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kalogs-c/nerd-backlog/pkg/validator"
)

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, payload T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	return v, nil
}

func DecodeValid[T validator.Validator](r *http.Request) (T, error) {
	v, err := Decode[T](r)
	if err != nil {
		return v, err
	}

	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, validator.ValidationError{Problems: problems}
	}

	return v, nil
}
