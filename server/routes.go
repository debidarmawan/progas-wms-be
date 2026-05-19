package server

import (
	"os"
	docs "progas-wms-be/docs"
	"progas-wms-be/handler"
	"progas-wms-be/constant"
	"progas-wms-be/helper"
	"progas-wms-be/middleware"
	"progas-wms-be/repository"
	"progas-wms-be/usecase"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/gorm"
)

func Routes(f *fiber.App, db *gorm.DB) {
	f.Use(RequestLogger())

	if os.Getenv("GO_ENV") == "development" {
		docs.SwaggerInfo.BasePath = "/api/v1"
		f.Get("/swagger/*", swaggo.HandlerDefault)
	}

	routerGroup := f.Group("/api/v1")
	routerGroup.Use(PanicHandler)
	routerGroup.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
	}))

	txManager := helper.NewTxManager(db)

	// INIT REPOSITORY
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	rbacRepo := repository.NewRbacRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	// INIT USECASE
	authUsecase := usecase.NewAuthUseCase(userRepo, auditLogRepo)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	userUsecase := usecase.NewUserUsecase(txManager, userRepo, roleRepo, auditLogRepo)

	// INIT HANDLER
	authHandler := handler.NewAuthHandler(authUsecase)
	roleHandler := handler.NewRoleHandler(roleUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	// Public routes
	authHandler.Routes(routerGroup)

	// Protected routes (JWT + RBAC per endpoint)
	protected := routerGroup.Group("", VerifyAuthToken)
	protected.Post("/logout", middleware.Authorize(rbacRepo, constant.PermAuthLogout), authHandler.Logout)
	roleHandler.Routes(protected.Group("", middleware.Authorize(rbacRepo, constant.PermRoleRead)))
	userHandler.Routes(protected.Group("", middleware.Authorize(rbacRepo, constant.PermUserCreate)))
}
