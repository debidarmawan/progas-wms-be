package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RequestLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${reqHeader:X-RequestId}] ${method} ${path} - ${status} | ${latency} | ${time}\n[${reqHeader:X-RequestId} | Req Header]\n${reqHeaders}\n[${reqHeader:X-RequestId} | Request Query Params] ${queryParams}\n[${reqHeader:X-RequestId} | Request Body] ${body}\n[${reqHeader:X-RequestId} | Response Body] ${resBody}\n\n",
	})
}
