package finance

import (
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/controller/finance"
	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/utils"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentController finance.PaymentController
}

func NewPaymentHandler(paymentController finance.PaymentController) *PaymentHandler {
	return &PaymentHandler{
		paymentController: paymentController,
	}
}

func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	siswaNIS := c.Query("siswa_nis")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, total, err := h.paymentController.GetAllPayments(siswaNIS, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.paymentController.GetPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	siswaNIS := c.Param("siswaNis")
	if siswaNIS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "siswa_nis is required"})
		return
	}

	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy, exists := utils.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	result, err := h.paymentController.CreatePayment(siswaNIS, &req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.paymentController.DeletePayment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment deleted successfully"})
}
