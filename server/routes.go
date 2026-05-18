package server

import (
	"os"
	_ "progas-wms-be/docs"
	"progas-wms-be/handler"
	"progas-wms-be/repository"
	"progas-wms-be/usecase"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/gorm"
)

func Routes(f *fiber.App, db *gorm.DB) {
	if os.Getenv("GO_ENV") == "development" {
		f.Get("/swagger/*", swaggo.HandlerDefault)
	}

	routerGroup := f.Group("/api/v1")
	routerGroup.Use(PanicHandler)

	f.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} - ${status} | ${latency} | ReqId ${reqHeader:X-RequestId} \n[Request Headers] ${reqHeaders}\n[Request Query Params] ${queryParams}\n[Request Body] ${body}\n[Response Body] ${resBody}\n\n",
	}))

	routerGroup.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
	}))

	// txManager := helper.NewTxManager(db)

	// INIT REPOSITORY
	userRepo := repository.NewUserRepository(db)

	// INIT USECASE
	authUsecase := usecase.NewAuthUseCase(userRepo)

	// INIT HANDLER
	authHandler := handler.NewAuthHandler(authUsecase)

	// ROUTING HANDLER
	authHandler.Routes(routerGroup)
}
