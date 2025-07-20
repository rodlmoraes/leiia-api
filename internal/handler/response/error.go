package response

// Error represents a uniform error response structure
type Error struct {
	Error string `json:"error"`
}

// NewError creates a new error response
func NewError(message string) Error {
	return Error{
		Error: message,
	}
}
