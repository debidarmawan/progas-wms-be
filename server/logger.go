package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RequestLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} - ${status} | ${latency} | ReqId ${reqHeader:X-RequestId} \n[Request Query Params] ${queryParams}\n[Request Body] ${body}\n[Response Body] ${resBody}\n\n",
	})
}
