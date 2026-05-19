package server

import (
	"os"
	"progas-wms-be/constant"
	docs "progas-wms-be/docs"
	"progas-wms-be/handler"
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

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	rbacRepo := repository.NewRbacRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	masterItemRepo := repository.NewMasterItemRepository(db)
	sparepartStockRepo := repository.NewSparepartStockRepository(db)
	cylinderRepo := repository.NewCylinderRepository(db)
	customerRepo := repository.NewCustomerRepository(db)

	authUsecase := usecase.NewAuthUseCase(userRepo, auditLogRepo)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	userUsecase := usecase.NewUserUsecase(txManager, userRepo, roleRepo, auditLogRepo)
	masterItemUsecase := usecase.NewMasterItemUsecase(txManager, masterItemRepo, sparepartStockRepo, auditLogRepo)
	cylinderUsecase := usecase.NewCylinderUsecase(txManager, cylinderRepo, masterItemRepo, customerRepo, auditLogRepo)
	customerUsecase := usecase.NewCustomerUsecase(txManager, customerRepo, auditLogRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	roleHandler := handler.NewRoleHandler(roleUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	masterItemHandler := handler.NewMasterItemHandler(masterItemUsecase)
	cylinderHandler := handler.NewCylinderHandler(cylinderUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	authHandler.Routes(routerGroup)

	protected := routerGroup.Group("", VerifyAuthToken)
	protected.Post("/logout", middleware.Authorize(rbacRepo, constant.PermAuthLogout), authHandler.Logout)
	roleHandler.Routes(protected.Group("", middleware.Authorize(rbacRepo, constant.PermRoleRead)))
	userHandler.Routes(protected.Group("", middleware.Authorize(rbacRepo, constant.PermUserCreate)))

	masterItemRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermMasterItemRead))
	masterItemRead.Get("/master-items", masterItemHandler.FindAll)
	masterItemRead.Get("/master-items/:id", masterItemHandler.FindById)

	masterItemWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermMasterItemWrite))
	masterItemWrite.Post("/master-items", masterItemHandler.Create)
	masterItemWrite.Put("/master-items/:id", masterItemHandler.Update)

	cylinderRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermCylinderRead))
	cylinderRead.Get("/cylinders", cylinderHandler.FindAll)
	cylinderRead.Get("/cylinders/:id", cylinderHandler.FindById)

	cylinderWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermCylinderWrite))
	cylinderWrite.Post("/cylinders", cylinderHandler.Create)

	customerRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermCustomerRead))
	customerRead.Get("/customers", customerHandler.FindAll)
	customerRead.Get("/customers/:id", customerHandler.FindById)

	customerWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermCustomerWrite))
	customerWrite.Post("/customers", customerHandler.Create)
	customerWrite.Put("/customers/:id", customerHandler.Update)
}
