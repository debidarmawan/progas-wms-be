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
	vendorRepo := repository.NewVendorRepository(db)
	fillingBatchRepo := repository.NewFillingBatchRepository(db)
	fleetRepo := repository.NewFleetRepository(db)
	deliveryOrderRepo := repository.NewDeliveryOrderRepository(db)
	cylinderLedgerRepo := repository.NewCylinderLedgerRepository(db)
	workOrderRepo := repository.NewWorkOrderRepository(db)
	sparepartMovementRepo := repository.NewSparepartMovementRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	authUsecase := usecase.NewAuthUseCase(userRepo, auditLogRepo)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	userUsecase := usecase.NewUserUsecase(txManager, userRepo, roleRepo, auditLogRepo)
	masterItemUsecase := usecase.NewMasterItemUsecase(txManager, masterItemRepo, sparepartStockRepo, auditLogRepo)
	cylinderUsecase := usecase.NewCylinderUsecase(txManager, cylinderRepo, masterItemRepo, customerRepo, vendorRepo, auditLogRepo)
	customerUsecase := usecase.NewCustomerUsecase(txManager, customerRepo, auditLogRepo)
	vendorUsecase := usecase.NewVendorUsecase(txManager, vendorRepo, cylinderRepo, auditLogRepo)
	inboundUsecase := usecase.NewInboundUsecase(txManager, cylinderRepo, cylinderLedgerRepo, auditLogRepo)
	fillingBatchUsecase := usecase.NewFillingBatchUsecase(txManager, fillingBatchRepo, cylinderRepo, cylinderLedgerRepo, masterItemRepo, auditLogRepo)
	fleetUsecase := usecase.NewFleetUsecase(txManager, fleetRepo, auditLogRepo)
	deliveryOrderUsecase := usecase.NewDeliveryOrderUsecase(txManager, deliveryOrderRepo, cylinderRepo, cylinderLedgerRepo, customerRepo, fleetRepo, auditLogRepo)
	exchangeUsecase := usecase.NewExchangeUsecase(txManager, cylinderRepo, cylinderLedgerRepo, customerRepo, auditLogRepo)
	workOrderUsecase := usecase.NewWorkOrderUsecase(txManager, workOrderRepo, masterItemRepo, sparepartStockRepo, sparepartMovementRepo, auditLogRepo)
	stockOpnameUsecase := usecase.NewStockOpnameUsecase(txManager, masterItemRepo, sparepartStockRepo, sparepartMovementRepo, auditLogRepo)
	hydrotestUsecase := usecase.NewHydrotestUsecase(txManager, cylinderRepo, cylinderLedgerRepo, auditLogRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo, sparepartStockRepo)
	reportUsecase := usecase.NewReportUsecase(cylinderLedgerRepo, cylinderRepo)
	inventoryUsecase := usecase.NewInventoryUsecase(customerRepo, cylinderRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	roleHandler := handler.NewRoleHandler(roleUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	masterItemHandler := handler.NewMasterItemHandler(masterItemUsecase)
	cylinderHandler := handler.NewCylinderHandler(cylinderUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)
	vendorHandler := handler.NewVendorHandler(vendorUsecase)
	inboundHandler := handler.NewInboundHandler(inboundUsecase)
	fillingBatchHandler := handler.NewFillingBatchHandler(fillingBatchUsecase)
	fleetHandler := handler.NewFleetHandler(fleetUsecase)
	deliveryOrderHandler := handler.NewDeliveryOrderHandler(deliveryOrderUsecase)
	exchangeHandler := handler.NewExchangeHandler(exchangeUsecase, rbacRepo)
	workOrderHandler := handler.NewWorkOrderHandler(workOrderUsecase)
	stockOpnameHandler := handler.NewStockOpnameHandler(stockOpnameUsecase)
	hydrotestHandler := handler.NewHydrotestHandler(hydrotestUsecase)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)
	reportHandler := handler.NewReportHandler(reportUsecase)
	inventoryHandler := handler.NewInventoryHandler(inventoryUsecase)

	authHandler.Routes(routerGroup)

	protected := routerGroup.Group("", VerifyAuthToken)
	protected.Post("/logout", middleware.Authorize(rbacRepo, constant.PermAuthLogout), authHandler.Logout)
	roleHandler.Routes(protected.Group("", middleware.Authorize(rbacRepo, constant.PermRoleRead)))
	userRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermUserRead))
	userRead.Get("/users", userHandler.FindAll)
	userRead.Get("/users/:id", userHandler.FindById)

	userWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermUserWrite))
	userWrite.Post("/users", userHandler.CreateUser)
	userWrite.Put("/users/:id", userHandler.UpdateUser)
	userWrite.Delete("/users/:id", userHandler.DeleteUser)

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

	vendorRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermVendorRead))
	vendorRead.Get("/vendors", vendorHandler.FindAll)
	vendorRead.Get("/vendors/:id", vendorHandler.FindById)

	vendorWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermVendorWrite))
	vendorWrite.Post("/vendors", vendorHandler.Create)
	vendorWrite.Put("/vendors/:id", vendorHandler.Update)
	vendorWrite.Delete("/vendors/:id", vendorHandler.Delete)

	inboundWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermInboundEmptyReceive))
	inboundWrite.Post("/inbound/empty-receive", inboundHandler.EmptyReceive)

	productionQCPreFill := protected.Group("", middleware.Authorize(rbacRepo, constant.PermProductionQCPreFill))
	productionQCPreFill.Post("/production/qc/pre-fill", inboundHandler.PreFillQC)

	productionQCPostFill := protected.Group("", middleware.Authorize(rbacRepo, constant.PermProductionQCPostFill))
	productionQCPostFill.Post("/production/qc/post-fill", inboundHandler.PostFillQC)

	fillingBatchRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermFillingBatchRead))
	fillingBatchRead.Get("/production/filling-batches", fillingBatchHandler.FindAll)
	fillingBatchRead.Get("/production/filling-batches/:id", fillingBatchHandler.FindById)

	fillingBatchWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermFillingBatchWrite))
	fillingBatchWrite.Post("/production/filling-batches", fillingBatchHandler.Submit)

	fleetRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermFleetRead))
	fleetRead.Get("/logistics/fleet", fleetHandler.FindAll)
	fleetRead.Get("/logistics/fleet/:id", fleetHandler.FindById)

	fleetWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermFleetWrite))
	fleetWrite.Post("/logistics/fleet", fleetHandler.Create)
	fleetWrite.Put("/logistics/fleet/:id", fleetHandler.Update)

	doRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermDORead))
	doRead.Get("/outbound/delivery-orders", deliveryOrderHandler.FindAll)
	doRead.Get("/outbound/delivery-orders/:id", deliveryOrderHandler.FindById)

	doWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermDOCreate))
	doWrite.Post("/outbound/delivery-orders", deliveryOrderHandler.Issue)

	exchangeWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermExchangeProcess))
	exchangeWrite.Post("/outbound/exchange", exchangeHandler.Process)

	workOrderRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermWorkOrderRead))
	workOrderRead.Get("/maintenance/work-orders", workOrderHandler.FindAll)
	workOrderRead.Get("/maintenance/work-orders/:id", workOrderHandler.FindById)

	workOrderWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermWorkOrderWrite))
	workOrderWrite.Post("/maintenance/work-orders", workOrderHandler.Create)
	workOrderWrite.Post("/maintenance/work-orders/:id/complete", workOrderHandler.Complete)

	stockOpnameWrite := protected.Group("", middleware.Authorize(rbacRepo, constant.PermInventoryStockOpname))
	stockOpnameWrite.Post("/inventory/spareparts/stock-opname", stockOpnameHandler.Submit)

	hydrotestRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermCylinderHydrotest))
	hydrotestRead.Get("/maintenance/hydrotest/due", hydrotestHandler.FindDue)
	hydrotestRead.Post("/maintenance/cylinders/:id/hydrotest", hydrotestHandler.Record)

	dashboardRead := protected.Group("", middleware.Authorize(rbacRepo, constant.PermDashboardRead))
	dashboardRead.Get("/dashboard/summary", dashboardHandler.GetSummary)

	reportLedger := protected.Group("", middleware.Authorize(rbacRepo, constant.PermReportLedger))
	reportLedger.Get("/reports/stock-ledger", reportHandler.StockLedger)

	reportTurnaround := protected.Group("", middleware.Authorize(rbacRepo, constant.PermReportTurnaround))
	reportTurnaround.Get("/reports/turnaround", reportHandler.Turnaround)

	virtualWarehouse := protected.Group("", middleware.Authorize(rbacRepo, constant.PermInventoryVirtual))
	virtualWarehouse.Get("/inventory/virtual-warehouse", inventoryHandler.VirtualWarehouse)
}
