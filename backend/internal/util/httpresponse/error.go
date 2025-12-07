package httpresponse

import (
	"devhub-backend/internal/domain/errs"
	"devhub-backend/internal/infra/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code" example:"TR-XXXXXX"`
	Message string `json:"message" example:"Error message"`
	Data    any    `json:"data,omitempty"`
}

// Error sends an error response with the appropriate status code.
func Error(c *gin.Context, err error) {
	// Unwrap the error and extract the error response.
	httpCode, errResp := unwrapError(err)
	// Log the error with the appropriate context.
	appLogger := logger.FromContext(c.Request.Context())
	appLogger.Error(c.Request.Context(), errResp.Message, err, nil)
	// Send the error response.
	c.AbortWithStatusJSON(httpCode, errResp)
}

// unwrapError processes the error and extracts information for the response.
func unwrapError(err error) (httpCode int, errResp ErrorResponse) {
	// Default error response for non-domain errors.
	// This will be used if the error is not a errs.DomainError.
	httpCode = http.StatusInternalServerError
	errResp = ErrorResponse{
		Code:    errs.StatusCodeGenericInternalServerError,
		Message: "An unexpected error occurred. Please try again later.",
	}

	// Try to unwrap the error and find a first valid errs.DomainError in the chain.
	if domainErr := UnwrapDomainError(err); domainErr != nil {
		httpCode = domainErr.GetHTTPCode()
		errResp.Code = domainErr.Code()
		errResp.Message = domainErr.GetMessage()
		errResp.Data = domainErr.GetData()
	}

	return httpCode, errResp
}

func UnwrapDomainError(err error) errs.DomainError {
	unwrapErr := err
	for unwrapErr != nil {
		// Check if the error explicitly implements DomainError and has a BaseError.
		if domainErr, ok := unwrapErr.(errs.DomainError); ok && errs.ExtractBaseError(domainErr) != nil {
			return domainErr
		}

		// Try to unwrap the next error in the chain.
		type unwrapper interface {
			Unwrap() error
		}
		// If the error does not implement an unwrapper, stop unwrapping.
		if unwrappableErr, ok := unwrapErr.(unwrapper); ok {
			unwrapErr = unwrappableErr.Unwrap()
		} else {
			break
		}
	}
	return nil
}
