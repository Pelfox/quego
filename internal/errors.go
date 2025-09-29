package internal

import (
	"github.com/labstack/echo/v4"
)

// ErrorCode defines a symbolic identifier for a category of errors that can
// occur in the responses. Error codes provide a stable way to classify errors.
type ErrorCode string

const (
	// ErrorCodeFunctionNotFound indicates that the requested function does not
	// exist (has not been registered in the system).
	ErrorCodeFunctionNotFound ErrorCode = "FUNCTION_NOT_FOUND"
	// ErrorCodeDatabase indicates that an error occurred during a database
	// operation, such as a failed query or connection issue.
	ErrorCodeDatabase ErrorCode = "DATABASE_ERROR"
	// ErrorCodeInvalidBody indicates that the request body could not be parsed.
	ErrorCodeInvalidBody ErrorCode = "INVALID_BODY"
)

// GenericError represents an application error that can be safely serialized
// and returned in an API response.
type GenericError struct {
	// Code identifies the type of error that occurred.
	Code ErrorCode `json:"code"`
	// Message provides additional context about the error.
	Message string `json:"message"`
}

// RespondError sends a JSON response containing a structured error to the
// client. It wraps the given error code and message in a `GenericError` and
// writes it with the specified HTTP status code.
func RespondError(ctx echo.Context, status int, code ErrorCode, message string) error {
	return ctx.JSON(status, &GenericError{Code: code, Message: message})
}
