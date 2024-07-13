package httpresponse

import "net/http"

type HttpError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

func (h *HttpError) Error() string {
	return h.Message
}

// WithStatusCode is function to set httpError statusCode
func (h *HttpError) WithStatusCode(statusCode int) *HttpError {
	h.StatusCode = statusCode
	return h
}

// WithMessage is function to set message
func (h *HttpError) WithMessage(message string) *HttpError {
	h.Message = message
	return h
}

// WithCode is function to set code
func (h *HttpError) WithCode(code string) *HttpError {
	h.Code = code
	return h
}

var (
	// 400
	ErrorBadRequest = &HttpError{StatusCode: http.StatusBadRequest, Code: "STU4000001", Message: "Bad Request"}

	// 404
	ErrorStudentNotFound = &HttpError{StatusCode: http.StatusNotFound, Code: "STU4040001", Message: "Student Not Found"}

	// 500
	ErrorInternalServerError = &HttpError{StatusCode: http.StatusInternalServerError, Code: "STU5000001", Message: "Internal Server Error"}
)
