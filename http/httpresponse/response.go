package httpresponse

import (
	"github.com/gin-gonic/gin"
	"github.com/rshby/go-redis-lock/http/httpresponse/dto"
	"net/http"
)

func ResponseError(ctx *gin.Context, httpError *HttpError) {
	ctx.Status(httpError.StatusCode)
	ctx.JSON(httpError.StatusCode, WrapApiResponse(httpError, "", nil))
}

func ResponseOK(ctx *gin.Context, message string, data any) {
	ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, WrapApiResponse(nil, message, data))
}

func WrapApiResponse(httpError *HttpError, message string, data any) *dto.ApiResponse {
	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	}

	if httpError != nil {
		response.StatusCode = httpError.StatusCode
		response.Message = httpError.Message
		response.ErrorCode = httpError.Code
	}

	return &response
}
