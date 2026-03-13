package models

type ErrorResponse struct {
	Error string
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Error: err.Error()}
}
