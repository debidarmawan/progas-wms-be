package global

import (
	"progas-wms-be/dto"

	"github.com/gofiber/fiber/v3"
)

var ResponseStatus = struct {
	FailedResponse  string
	SuccessResponse string
	RetryResponse   string
}{
	FailedResponse:  "FAILED",
	SuccessResponse: "OK",
	RetryResponse:   "RETRY",
}

type Response[T any] struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Data      T      `json:"data"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`
}

type Result[T any] struct {
	Data       *T
	Error      error
	StatusCode int
}

type ResultChan[T any] struct {
	Data  *T
	Error ErrorResponse
}

func (r *Result[T]) ToResponseError() *Response[*T] {
	return &Response[*T]{
		Code:    r.StatusCode,
		Data:    r.Data,
		Status:  ResponseStatus.FailedResponse,
		Message: r.Error.Error(),
	}
}

func CreateResponse[T any](data *T, statusCode int, c fiber.Ctx) error {
	response := Response[T]{
		Code:    statusCode,
		Data:    *data,
		Status:  ResponseStatus.SuccessResponse,
		Message: "",
	}

	return c.Status(statusCode).JSON(response)
}

func MessageResponse(message string, statusCode int, c fiber.Ctx) error {
	return CreateResponse(&dto.Message{Message: message}, statusCode, c)
}

func CreateMessageResponse(message string, statusCode int, c fiber.Ctx) error {
	response := Response[any]{
		Code:    statusCode,
		Data:    nil,
		Status:  ResponseStatus.SuccessResponse,
		Message: message,
	}

	return c.Status(statusCode).JSON(response)
}
