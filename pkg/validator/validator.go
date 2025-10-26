package validator

import (
	"context"
	"fmt"
)

type Problems map[string][]string

type Validator interface {
	Valid(ctx context.Context) Problems
}

type ValidationError struct {
	Problems Problems
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %d problems", len(e.Problems))
}
