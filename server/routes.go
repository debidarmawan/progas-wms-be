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
	"gorm.io/gorm"
)

func Routes(f *fiber.App, db *gorm.DB) {
	f.Use(RequestLogger())

	if os.Getenv("GO_ENV") == "development" || os.Getenv("GO_ENV") == "localhost" {
		docs.SwaggerInfo.BasePath = "/api/v1"
		f.Get("/swagger/*", swaggo.New(swaggo.Config{
			InstanceName: "swagger",
			URL:          "/swagger/doc.json",
		}))
	}

	routerGroup := f.Group("/api/v1")
	routerGroup.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})
	routerGroup.Use(PanicHandler)
	routerGroup.Use(corsMiddleware())

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
	driverRepo := repository.NewDriverRepository(db)
	deliveryOrderRepo := repository.NewDeliveryOrderRepository(db)
	customerPORepo := repository.NewCustomerPORepository(db)
	salesOrderRepo := repository.NewSalesOrderRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	cylinderLedgerRepo := repository.NewCylinderLedgerRepository(db)
	workOrderRepo := repository.NewWorkOrderRepository(db)
	sparepartMovementRepo := repository.NewSparepartMovementRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	authUsecase := usecase.NewAuthUseCase(userRepo, auditLogRepo, rbacRepo)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	userUsecase := usecase.NewUserUsecase(txManager, userRepo, roleRepo, driverRepo, auditLogRepo)
	masterItemUsecase := usecase.NewMasterItemUsecase(txManager, masterItemRepo, sparepartStockRepo, auditLogRepo)
	cylinderUsecase := usecase.NewCylinderUsecase(txManager, cylinderRepo, masterItemRepo, customerRepo, vendorRepo, cylinderLedgerRepo, userRepo, auditLogRepo)
	customerUsecase := usecase.NewCustomerUsecase(txManager, customerRepo, auditLogRepo)
	customerItemPriceRepo := repository.NewCustomerItemPriceRepository(db)
	customerItemPriceUsecase := usecase.NewCustomerItemPriceUsecase(txManager, customerItemPriceRepo, customerRepo, masterItemRepo, auditLogRepo)
	vendorUsecase := usecase.NewVendorUsecase(txManager, vendorRepo, cylinderRepo, auditLogRepo)
	inboundUsecase := usecase.NewInboundUsecase(txManager, cylinderRepo, cylinderLedgerRepo, auditLogRepo)
	fillingBatchUsecase := usecase.NewFillingBatchUsecase(txManager, fillingBatchRepo, cylinderRepo, cylinderLedgerRepo, masterItemRepo, auditLogRepo)
	fleetUsecase := usecase.NewFleetUsecase(txManager, fleetRepo, auditLogRepo)
	driverUsecase := usecase.NewDriverUsecase(txManager, driverRepo, userRepo, auditLogRepo)
	invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepo, customerRepo)
	paymentUsecase := usecase.NewPaymentUsecase(txManager, paymentRepo, invoiceRepo, customerRepo, auditLogRepo)
	deliveryOrderUsecase := usecase.NewDeliveryOrderUsecase(txManager, deliveryOrderRepo, cylinderRepo, cylinderLedgerRepo, customerRepo, fleetRepo, auditLogRepo, customerItemPriceUsecase, salesOrderRepo, invoiceUsecase)
	customerPOUsecase := usecase.NewCustomerPOUsecase(txManager, customerPORepo, customerRepo, masterItemRepo, auditLogRepo)
	salesOrderUsecase := usecase.NewSalesOrderUsecase(txManager, salesOrderRepo, customerPORepo, customerRepo, masterItemRepo, auditLogRepo)
	exchangeUsecase := usecase.NewExchangeUsecase(txManager, cylinderRepo, cylinderLedgerRepo, customerRepo, auditLogRepo)
	workOrderUsecase := usecase.NewWorkOrderUsecase(txManager, workOrderRepo, masterItemRepo, sparepartStockRepo, sparepartMovementRepo, auditLogRepo)
	stockOpnameUsecase := usecase.NewStockOpnameUsecase(txManager, masterItemRepo, sparepartStockRepo, sparepartMovementRepo, auditLogRepo)
	hydrotestUsecase := usecase.NewHydrotestUsecase(txManager, cylinderRepo, cylinderLedgerRepo, auditLogRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo, sparepartStockRepo, cylinderRepo, customerRepo)
	reportUsecase := usecase.NewReportUsecase(cylinderLedgerRepo, cylinderRepo)
	inventoryUsecase := usecase.NewInventoryUsecase(customerRepo, cylinderRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	roleHandler := handler.NewRoleHandler(roleUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	masterItemHandler := handler.NewMasterItemHandler(masterItemUsecase)
	cylinderHandler := handler.NewCylinderHandler(cylinderUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)
	customerItemPriceHandler := handler.NewCustomerItemPriceHandler(customerItemPriceUsecase)
	vendorHandler := handler.NewVendorHandler(vendorUsecase)
	inboundHandler := handler.NewInboundHandler(inboundUsecase)
	fillingBatchHandler := handler.NewFillingBatchHandler(fillingBatchUsecase)
	fleetHandler := handler.NewFleetHandler(fleetUsecase)
	driverHandler := handler.NewDriverHandler(driverUsecase)
	deliveryOrderHandler := handler.NewDeliveryOrderHandler(deliveryOrderUsecase)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUsecase, paymentUsecase)
	customerPOHandler := handler.NewCustomerPOHandler(customerPOUsecase)
	salesOrderHandler := handler.NewSalesOrderHandler(salesOrderUsecase)
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
	protected.Get("/profile", authHandler.Profile)
	protected.Get("/roles/", middleware.Authorize(rbacRepo, constant.PermRoleRead), roleHandler.FindAll)
	protected.Get("/roles/:id", middleware.Authorize(rbacRepo, constant.PermRoleRead), roleHandler.FindById)
	userRead := middleware.Authorize(rbacRepo, constant.PermUserRead)
	protected.Get("/users", userRead, userHandler.FindAll)
	protected.Get("/users/:id", userRead, userHandler.FindById)

	userWrite := middleware.Authorize(rbacRepo, constant.PermUserWrite)
	protected.Post("/users", userWrite, userHandler.CreateUser)
	protected.Put("/users/:id", userWrite, userHandler.UpdateUser)
	protected.Delete("/users/:id", userWrite, userHandler.DeleteUser)

	masterItemRead := middleware.Authorize(rbacRepo, constant.PermMasterItemRead)
	protected.Get("/master-items", masterItemRead, masterItemHandler.FindAll)
	protected.Get("/master-items/:id", masterItemRead, masterItemHandler.FindById)

	masterItemWrite := middleware.Authorize(rbacRepo, constant.PermMasterItemWrite)
	protected.Post("/master-items", masterItemWrite, masterItemHandler.Create)
	protected.Post("/master-items/bulk", masterItemWrite, masterItemHandler.CreateBulk)
	protected.Put("/master-items/:id", masterItemWrite, masterItemHandler.Update)

	cylinderRead := middleware.Authorize(rbacRepo, constant.PermCylinderRead)
	protected.Get("/cylinders", cylinderRead, cylinderHandler.FindAll)
	protected.Get("/cylinders/by-barcode/:sn", cylinderRead, cylinderHandler.FindByBarcode)
	protected.Get("/cylinders/:id", cylinderRead, cylinderHandler.FindById)
	protected.Get("/cylinders/:id/history", cylinderRead, cylinderHandler.History)

	cylinderWrite := middleware.Authorize(rbacRepo, constant.PermCylinderWrite)
	protected.Post("/cylinders", cylinderWrite, cylinderHandler.Create)
	protected.Put("/cylinders/:id", cylinderWrite, cylinderHandler.Update)

	customerRead := middleware.Authorize(rbacRepo, constant.PermCustomerRead)
	protected.Get("/customers", customerRead, customerHandler.FindAll)
	protected.Get("/customers/:id", customerRead, customerHandler.FindById)

	customerWrite := middleware.Authorize(rbacRepo, constant.PermCustomerWrite)
	protected.Post("/customers", customerWrite, customerHandler.Create)
	protected.Put("/customers/:id", customerWrite, customerHandler.Update)

	customerPricingRead := middleware.Authorize(rbacRepo, constant.PermCustomerPricingRead)
	protected.Get("/customers/:customerId/pricing", customerPricingRead, customerItemPriceHandler.FindAllByCustomer)
	protected.Get("/customers/:customerId/items/:masterItemId/pricing", customerPricingRead, customerItemPriceHandler.FindActive)

	customerPricingWrite := middleware.Authorize(rbacRepo, constant.PermCustomerPricingWrite)
	protected.Post("/customers/:customerId/pricing", customerPricingWrite, customerItemPriceHandler.Create)
	protected.Put("/customers/:customerId/pricing/:id", customerPricingWrite, customerItemPriceHandler.Update)
	protected.Delete("/customers/:customerId/pricing/:id", customerPricingWrite, customerItemPriceHandler.Delete)

	customerPORead := middleware.Authorize(rbacRepo, constant.PermCustomerPORead)
	protected.Get("/sales/customer-pos", customerPORead, customerPOHandler.FindAll)
	protected.Get("/sales/customer-pos/:id", customerPORead, customerPOHandler.FindById)

	customerPOWrite := middleware.Authorize(rbacRepo, constant.PermCustomerPOWrite)
	protected.Post("/sales/customer-pos", customerPOWrite, customerPOHandler.Create)
	protected.Put("/sales/customer-pos/:id/confirm", customerPOWrite, customerPOHandler.Confirm)

	salesOrderRead := middleware.Authorize(rbacRepo, constant.PermSalesOrderRead)
	protected.Get("/sales/sales-orders", salesOrderRead, salesOrderHandler.FindAll)
	protected.Get("/sales/sales-orders/:id", salesOrderRead, salesOrderHandler.FindById)

	salesOrderWrite := middleware.Authorize(rbacRepo, constant.PermSalesOrderWrite)
	protected.Post("/sales/sales-orders", salesOrderWrite, salesOrderHandler.Create)
	protected.Put("/sales/sales-orders/:id/confirm", salesOrderWrite, salesOrderHandler.Confirm)
	protected.Put("/sales/sales-orders/:id/cancel", salesOrderWrite, salesOrderHandler.Cancel)

	vendorRead := middleware.Authorize(rbacRepo, constant.PermVendorRead)
	protected.Get("/vendors", vendorRead, vendorHandler.FindAll)
	protected.Get("/vendors/:id", vendorRead, vendorHandler.FindById)

	vendorWrite := middleware.Authorize(rbacRepo, constant.PermVendorWrite)
	protected.Post("/vendors", vendorWrite, vendorHandler.Create)
	protected.Put("/vendors/:id", vendorWrite, vendorHandler.Update)
	protected.Delete("/vendors/:id", vendorWrite, vendorHandler.Delete)

	inboundWrite := middleware.Authorize(rbacRepo, constant.PermInboundEmptyReceive)
	protected.Post("/inbound/empty-receive", inboundWrite, inboundHandler.EmptyReceive)

	productionQCPreFill := middleware.Authorize(rbacRepo, constant.PermProductionQCPreFill)
	protected.Post("/production/qc/pre-fill", productionQCPreFill, inboundHandler.PreFillQC)

	productionQCPostFill := middleware.Authorize(rbacRepo, constant.PermProductionQCPostFill)
	protected.Post("/production/qc/post-fill", productionQCPostFill, inboundHandler.PostFillQC)

	fillingBatchRead := middleware.Authorize(rbacRepo, constant.PermFillingBatchRead)
	protected.Get("/production/filling-batches", fillingBatchRead, fillingBatchHandler.FindAll)
	protected.Get("/production/filling-batches/:id", fillingBatchRead, fillingBatchHandler.FindById)

	fillingBatchWrite := middleware.Authorize(rbacRepo, constant.PermFillingBatchWrite)
	protected.Post("/production/filling-batches", fillingBatchWrite, fillingBatchHandler.Submit)

	fleetRead := middleware.Authorize(rbacRepo, constant.PermFleetRead)
	protected.Get("/logistics/fleet", fleetRead, fleetHandler.FindAll)
	protected.Get("/logistics/fleet/:id", fleetRead, fleetHandler.FindById)

	fleetWrite := middleware.Authorize(rbacRepo, constant.PermFleetWrite)
	protected.Post("/logistics/fleet", fleetWrite, fleetHandler.Create)
	protected.Put("/logistics/fleet/:id", fleetWrite, fleetHandler.Update)

	driverRead := middleware.Authorize(rbacRepo, constant.PermDriverRead)
	protected.Get("/logistics/drivers", driverRead, driverHandler.FindAll)
	protected.Get("/logistics/drivers/:id", driverRead, driverHandler.FindById)

	driverWrite := middleware.Authorize(rbacRepo, constant.PermDriverWrite)
	protected.Post("/logistics/drivers", driverWrite, driverHandler.Create)
	protected.Put("/logistics/drivers/:id", driverWrite, driverHandler.Update)
	protected.Delete("/logistics/drivers/:id", driverWrite, driverHandler.Delete)

	doRead := middleware.Authorize(rbacRepo, constant.PermDORead)
	protected.Get("/outbound/delivery-orders", doRead, deliveryOrderHandler.FindAll)
	protected.Get("/outbound/delivery-orders/:id", doRead, deliveryOrderHandler.FindById)

	doWrite := middleware.Authorize(rbacRepo, constant.PermDOCreate)
	protected.Post("/outbound/delivery-orders", doWrite, deliveryOrderHandler.Issue)

	exchangeWrite := middleware.Authorize(rbacRepo, constant.PermExchangeProcess)
	protected.Post("/outbound/exchange", exchangeWrite, exchangeHandler.Process)

	invoiceRead := middleware.Authorize(rbacRepo, constant.PermInvoiceRead)
	protected.Get("/finance/invoices", invoiceRead, invoiceHandler.FindAll)
	protected.Get("/finance/invoices/:id", invoiceRead, invoiceHandler.FindById)

	paymentWrite := middleware.Authorize(rbacRepo, constant.PermPaymentWrite)
	protected.Post("/finance/invoices/:id/payments", paymentWrite, invoiceHandler.RecordPayment)

	workOrderRead := middleware.Authorize(rbacRepo, constant.PermWorkOrderRead)
	protected.Get("/maintenance/work-orders", workOrderRead, workOrderHandler.FindAll)
	protected.Get("/maintenance/work-orders/:id", workOrderRead, workOrderHandler.FindById)

	workOrderWrite := middleware.Authorize(rbacRepo, constant.PermWorkOrderWrite)
	protected.Post("/maintenance/work-orders", workOrderWrite, workOrderHandler.Create)
	protected.Post("/maintenance/work-orders/:id/complete", workOrderWrite, workOrderHandler.Complete)

	stockOpnameWrite := middleware.Authorize(rbacRepo, constant.PermInventoryStockOpname)
	protected.Post("/inventory/spareparts/stock-opname", stockOpnameWrite, stockOpnameHandler.Submit)

	hydrotestRead := middleware.Authorize(rbacRepo, constant.PermCylinderHydrotest)
	protected.Get("/maintenance/hydrotest/due", hydrotestRead, hydrotestHandler.FindDue)
	protected.Post("/maintenance/cylinders/:id/hydrotest", hydrotestRead, hydrotestHandler.Record)

	dashboardRead := middleware.Authorize(rbacRepo, constant.PermDashboardRead)
	protected.Get("/dashboard/summary", dashboardRead, dashboardHandler.GetSummary)

	reportLedger := middleware.Authorize(rbacRepo, constant.PermReportLedger)
	protected.Get("/reports/stock-ledger", reportLedger, reportHandler.StockLedger)

	reportTurnaround := middleware.Authorize(rbacRepo, constant.PermReportTurnaround)
	protected.Get("/reports/turnaround", reportTurnaround, reportHandler.Turnaround)

	virtualWarehouse := middleware.Authorize(rbacRepo, constant.PermInventoryVirtual)
	protected.Get("/inventory/virtual-warehouse", virtualWarehouse, inventoryHandler.VirtualWarehouse)
}
