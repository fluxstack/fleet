package errors

type APIError struct {
	Code     int            `json:"code"`
	Message  string         `json:"message"`
	Reason   string         `json:"reason"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
