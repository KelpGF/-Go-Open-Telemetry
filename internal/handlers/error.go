package handlers

type ResponseError struct {
	Message string `json:"message"`
}

func newResponseError(message string) ResponseError {
	return ResponseError{Message: message}
}
