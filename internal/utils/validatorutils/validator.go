package validatorutils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rshby/go-redis-lock/http/httpresponse"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

// GetHttpErrorByTag is function to get httpError by tag
func GetHttpErrorByTag(err error) *httpresponse.HttpError {
	validationError, ok := err.(validator.ValidationErrors)
	if !ok {
		return httpresponse.ErrorBadRequest
	}

	var httpError *httpresponse.HttpError
	for _, fieldError := range validationError {
		switch fieldError.Tag() {
		case "required":
			httpError = httpresponse.ErrorBadRequest.WithMessage(fmt.Sprintf("%s is required", fieldError.Field()))
			return httpError
		case "email":
			httpError = httpresponse.ErrorBadRequest.WithMessage(fmt.Sprintf("%s should be an email format", fieldError.Field()))
			return httpError
		case "gt":
			httpError = httpresponse.ErrorBadRequest.WithMessage(fmt.Sprintf("count must be more than %s", fieldError.Param()))
			return httpError
		default:
			httpError = httpresponse.ErrorBadRequest
			return httpError
		}
	}

	return httpError
}
