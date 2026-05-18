package server

import (
	"log"
	"progas-wms-be/config"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequestId(c fiber.Ctx) error {
	requestid := uuid.NewString()
	c.Request().Header.Add("X-RequestId", requestid)
	c.Set("X-RequestId", requestid)

	return c.Next()
}

func PanicHandler(c fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())

			log.Println("[PANIC]")

			if e, ok := r.(error); ok {
				log.Println(e.Error())
			} else {
				log.Println(r)
			}

			log.Println(stack)

			err = c.Status(fiber.StatusInternalServerError).JSON(global.Response[any]{
				Code:    fiber.StatusInternalServerError,
				Data:    nil,
				Status:  global.ResponseStatus.FailedResponse,
				Message: "Internal Server Error",
			})
		}
	}()

	return c.Next()
}

func VerifyAuthToken(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(global.Response[any]{
			Code:    fiber.StatusUnauthorized,
			Status:  global.ResponseStatus.FailedResponse,
			Message: "Unauthorized",
		})
	}

	tokenString := authHeader[7:]

	token, err := jwt.ParseWithClaims(tokenString, &dto.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv(constant.AuthTokenSecretKey)), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(global.Response[any]{
			Code:    fiber.StatusUnauthorized,
			Status:  global.ResponseStatus.FailedResponse,
			Message: "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(*dto.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(global.Response[any]{
			Code:    fiber.StatusUnauthorized,
			Status:  global.ResponseStatus.FailedResponse,
			Message: "Invalid token claims",
		})
	}

	c.Locals("user_id", claims.UserId)
	c.Locals("role_id", claims.RoleId)

	return c.Next()
}
