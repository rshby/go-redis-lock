package dto

type ApiResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	ErrorCode  string `json:"error_code,omitempty"`
}
