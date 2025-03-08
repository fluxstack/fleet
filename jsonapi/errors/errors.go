package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func FromError(err error) *APIError {
	return &APIError{
		Code:    -1,
		Message: err.Error(),
		err:     err,
	}
}

type APIError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	err     error
}

func (e *APIError) Wrap(err error) {
	e.err = err
}

func (e *APIError) Unwrap() error {
	return e.err
}

type ErrorItem struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "fastapi: Error %d: ", e.Code)
	if e.Message != "" {
		fmt.Fprintf(&buf, "%s", e.Message)
	}
	if len(e.Details) > 0 {
		var detailBuf bytes.Buffer
		enc := json.NewEncoder(&detailBuf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(e.Details); err == nil {
			fmt.Fprint(&buf, "\nDetails:")
			fmt.Fprintf(&buf, "\n%s", detailBuf.String())

		}
	}

	fmt.Fprintln(&buf, "\nMore details:")

	return buf.String()
}

var _ error = new(APIError)
