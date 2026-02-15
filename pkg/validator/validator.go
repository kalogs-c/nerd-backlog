package validator

import (
	"context"
	"fmt"
)

type Problems map[string][]string

func (p Problems) Add(field string, message string) {
	p[field] = append(p[field], message)
}

type Validator interface {
	Valid(ctx context.Context) Problems
}

type ValidationError struct {
	Problems Problems
}

func (e ValidationError) Error() string {
	errorMessage := fmt.Sprintf("validation failed: %d problems", len(e.Problems))

	for key, p := range e.Problems {
		for _, msg := range p {
			errorMessage = fmt.Sprintf("%s\n\t[%s]: %s", errorMessage, key, msg)
		}
	}

	return errorMessage
}
