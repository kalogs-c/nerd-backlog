package httpjson

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ErrorResponse struct {
	Title         string         `json:"title"`
	Detail        string         `json:"detail,omitempty"`
	Status        int            `json:"status"`
	InvalidParams []InvalidParam `json:"invalid-params,omitempty"`
}

func EncodeError(w http.ResponseWriter, r *http.Request, status int, title string, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Title:  title,
		Detail: detail,
		Status: status,
	}

	json.NewEncoder(w).Encode(resp)
}

func EncodeValidationErrors(w http.ResponseWriter, r *http.Request, problems map[string][]string) {
	invalidParams := make([]InvalidParam, 0, len(problems))
	for field, reasons := range problems {
		for _, reason := range reasons {
			invalidParams = append(invalidParams, InvalidParam{
				Name:   field,
				Reason: reason,
			})
		}
	}

	resp := ErrorResponse{
		Title:         "Validation failed",
		Detail:        "Some fields are invalid",
		Status:        http.StatusUnprocessableEntity,
		InvalidParams: invalidParams,
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(resp)
}

func NotifyError(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	code int,
	title,
	detail string,
	err error,
) {
	if err != nil {
		logger.ErrorContext(ctx, title, "err", err.Error())
	} else {
		logger.ErrorContext(ctx, title)
	}

	EncodeError(w, r, code, title, detail)
}
