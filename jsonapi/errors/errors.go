package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
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
	Errors  []ErrorItem    `json:"errors,omitempty"`
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
	if len(e.Errors) == 0 && e.Message == "" {
		return fmt.Sprintf("fastapi: got HTTP response code %d", e.Code)
	}
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
	if len(e.Errors) == 0 {
		return strings.TrimSpace(buf.String())
	}
	if len(e.Errors) == 1 && e.Errors[0].Message == e.Message {
		fmt.Fprintf(&buf, ", %s", e.Errors[0].Reason)
		return buf.String()
	}
	fmt.Fprintln(&buf, "\nMore details:")
	for _, v := range e.Errors {
		fmt.Fprintf(&buf, "Reason: %s, Message: %s\n", v.Reason, v.Message)
	}
	return buf.String()
}

var _ error = new(APIError)
