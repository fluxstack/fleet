package errors

import "github.com/lynx-go/x/json"

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

func (e *APIError) Wrap(err error) *APIError {
	e.err = err
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details["internal_error"] = err.Error()
	return e
}

func (e *APIError) Unwrap() error {
	return e.err
}

type ErrorItem struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return json.SafeMarshalString(e)
}

var _ error = new(APIError)
