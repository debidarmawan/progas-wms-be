package global

import (
	"errors"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
)

type ErrorResponse interface {
	Error() string
	GetMessage() string
	GetCode() int
	GetData() any
	ToResponse(c fiber.Ctx) error
	GetStack() string
	ToError() error
}

type errorResponse struct {
	StatusCode int
	ErrorCode  string
	Message    string
	Detail     error
	Data       any
	Stack      string
}

func (e *errorResponse) Error() string {
	if e.Detail != nil {
		return e.Detail.Error()
	}

	return e.Message
}

func (e *errorResponse) GetMessage() string {
	return e.Message
}

func (e *errorResponse) GetCode() int {
	return e.StatusCode
}

func (e *errorResponse) GetData() any {
	return e.Data
}

func (e *errorResponse) GetStack() string {
	return e.Stack
}

func (e *errorResponse) ToError() error {
	if e.Detail != nil {
		return e.Detail
	}

	return errors.New(e.Message)
}

func (e *errorResponse) ToResponse(c fiber.Ctx) error {
	return c.Status(e.StatusCode).JSON(Response[any]{
		Code:      e.StatusCode,
		Status:    ResponseStatus.FailedResponse,
		Data:      e.Data,
		Message:   e.Message,
		ErrorCode: e.ErrorCode,
	})
}

func BadRequestError(message string) ErrorResponse {
	return &errorResponse{
		StatusCode: fiber.StatusBadRequest,
		ErrorCode:  "BAD_REQUEST",
		Message:    message,
		Detail:     nil,
	}
}

func BadRequestErrorWithData(message string, data any, errorCode string) ErrorResponse {
	return &errorResponse{
		StatusCode: fiber.StatusBadRequest,
		ErrorCode:  errorCode,
		Message:    message,
		Data:       data,
		Detail:     nil,
	}
}

func InternalServerError(err error) ErrorResponse {
	return &errorResponse{
		StatusCode: fiber.StatusInternalServerError,
		ErrorCode:  "INTERNAL_SERVER_ERROR",
		Message:    err.Error(),
		Detail:     err,
		Stack:      string(debug.Stack()),
	}
}

func NotFoundError(messages ...string) ErrorResponse {
	message := "Not Found"
	if len(messages) > 0 {
		message = messages[0]
	}

	return &errorResponse{
		StatusCode: fiber.StatusNotFound,
		ErrorCode:  "NOT_FOUND",
		Message:    message,
		Detail:     nil,
		Data:       nil,
	}
}

func ForbiddenError() ErrorResponse {
	return &errorResponse{
		StatusCode: fiber.StatusForbidden,
		ErrorCode:  "FORBIDDEN",
		Message:    "Forbidden",
		Detail:     nil,
		Data:       nil,
	}
}

func UnauthorizedError() ErrorResponse {
	return &errorResponse{
		StatusCode: fiber.StatusUnauthorized,
		ErrorCode:  "UNAUTHORIZED",
		Message:    "Unauthorized",
		Detail:     nil,
		Data:       nil,
	}
}
