package finance

import (
	"github.com/entertrans/backend-bogor.git/internal/config"
	financeController "github.com/entertrans/backend-bogor.git/internal/controller/finance"
	financeHandler "github.com/entertrans/backend-bogor.git/internal/handler/finance"
	"github.com/entertrans/backend-bogor.git/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Initialize controllers
	feeTemplateController := financeController.NewFeeTemplateController(db)
	invoiceController := financeController.NewInvoiceController(db)
	paymentController := financeController.NewPaymentController(db)

	// Initialize handlers
	feeTemplateHandler := financeHandler.NewFeeTemplateHandler(feeTemplateController)
	invoiceHandler := financeHandler.NewInvoiceHandler(invoiceController)
	paymentHandler := financeHandler.NewPaymentHandler(paymentController)

	// Finance routes group with authentication
	finance := rg.Group("/finance")
	finance.Use(middleware.AuthMiddleware(cfg))
	{
		// ==================== FEE TEMPLATES ====================
		finance.GET("/fee-templates", feeTemplateHandler.GetAllFeeTemplates)
		finance.GET("/fee-templates/:id", feeTemplateHandler.GetFeeTemplateByID)
		finance.POST("/fee-templates", feeTemplateHandler.CreateFeeTemplate)
		finance.PUT("/fee-templates/:id", feeTemplateHandler.UpdateFeeTemplate)
		finance.DELETE("/fee-templates/:id", feeTemplateHandler.DeleteFeeTemplate)

		// ==================== INVOICES ====================
		finance.GET("/invoices", invoiceHandler.GetAllInvoices)
		finance.GET("/invoices/:id", invoiceHandler.GetInvoiceByID)
		finance.POST("/invoices", invoiceHandler.CreateInvoice)
		finance.PUT("/invoices/:id", invoiceHandler.UpdateInvoice)
		finance.DELETE("/invoices/:id", invoiceHandler.DeleteInvoice)
		finance.PUT("/invoices/:id/publish", invoiceHandler.PublishInvoice)
		finance.PUT("/invoices/:id/cancel", invoiceHandler.CancelInvoice)

		// Invoice Items
		finance.POST("/invoices/:id/items", invoiceHandler.AddInvoiceItem)
		finance.PUT("/invoices/items/:itemId", invoiceHandler.UpdateInvoiceItem)
		finance.DELETE("/invoices/items/:itemId", invoiceHandler.RemoveInvoiceItem)

		// Invoice Adjustments
		finance.POST("/invoices/:id/adjustments", invoiceHandler.AddInvoiceAdjustment)
		finance.PUT("/invoices/adjustments/:adjustmentId", invoiceHandler.UpdateInvoiceAdjustment)
		finance.DELETE("/invoices/adjustments/:adjustmentId", invoiceHandler.RemoveInvoiceAdjustment)

		// ==================== PAYMENTS ====================
		finance.GET("/payments", paymentHandler.GetAllPayments)
		finance.GET("/payments/:id", paymentHandler.GetPaymentByID)
		finance.POST("/payments/siswa/:siswaNis", paymentHandler.CreatePayment)
		finance.DELETE("/payments/:id", paymentHandler.DeletePayment)
	}
}
