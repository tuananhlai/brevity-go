package shared

type ErrorResponse struct {
	// Code is a short error code. Examples: "invalid_credentials", "user_already_exists".
	Code Code `json:"errorCode"`
	// Message is a human-readable error message.
	Message string `json:"message,omitempty"`
	Details any    `json:"details,omitempty"`
}

// Code is a string representing a type of server error.
type Code string

const (
	CodeUnknown             Code = "unknown"
	CodeBindingRequestError Code = "binding_request_error"
	CodeUnauthorized        Code = "unauthorized"
)

// Clamp returns value if it is between min and max, otherwise returns min or max.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
