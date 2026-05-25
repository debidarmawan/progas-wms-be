package server

import (
	"os"
	"progas-wms-be/constant"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RequestLogger() fiber.Handler {
	if os.Getenv(constant.GoEnv) == "development" {
		return logger.New(logger.Config{
			Format: "[${time}] ${method} ${path} - ${status} | ${latency} | ReqId ${reqHeader:X-RequestId}\n[Query] ${queryParams}\n",
		})
	}

	return logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} - ${status} | ${latency} | ReqId ${reqHeader:X-RequestId}\n",
	})
}
