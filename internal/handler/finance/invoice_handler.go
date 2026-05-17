package finance

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/controller/finance"
	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/utils"
	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	invoiceController finance.InvoiceController
}

func NewInvoiceHandler(invoiceController finance.InvoiceController) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceController: invoiceController,
	}
}

func (h *InvoiceHandler) GetAllInvoices(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, total, err := h.invoiceController.GetAllInvoices(status, page, limit)
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

func (h *InvoiceHandler) GetInvoiceByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.invoiceController.GetInvoiceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest

	// Baca body untuk logging TANPA mengkonsumsi
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
	}
	// KEMBALIKAN body ke request
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Log raw body
	fmt.Printf("Raw request body: %s\n", string(body))

	// Sekarang binding akan berhasil karena body sudah dikembalikan
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Debug: print parsed request
	fmt.Printf("Parsed: IssueDate=%s, DueDate=%s, Description=%s\n",
		req.Invoice.IssueDate, req.Invoice.DueDate, req.Invoice.Description)
	fmt.Printf("Items count: %d\n", len(req.Items))

	createdBy, exists := utils.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Convert to InvoiceRequest for controller
	var dueDatePtr *string
	if req.Invoice.DueDate != "" {
		dueDatePtr = &req.Invoice.DueDate
	}

	var descriptionPtr *string
	if req.Invoice.Description != "" {
		descriptionPtr = &req.Invoice.Description
	}

	invoiceReq := &dto.InvoiceRequest{
		IssueDate:   req.Invoice.IssueDate,
		DueDate:     dueDatePtr,
		Description: descriptionPtr,
	}

	result, err := h.invoiceController.CreateInvoice(invoiceReq, req.Items, req.Adjustments, createdBy)
	if err != nil {
		fmt.Printf("CreateInvoice error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.InvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.invoiceController.UpdateInvoice(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.invoiceController.DeleteInvoice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invoice deleted successfully"})
}

func (h *InvoiceHandler) PublishInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.invoiceController.PublishInvoice(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invoice published successfully"})
}

func (h *InvoiceHandler) CancelInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.invoiceController.CancelInvoice(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invoice cancelled successfully"})
}

func (h *InvoiceHandler) AddInvoiceItem(c *gin.Context) {
	invoiceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice id"})
		return
	}

	var req dto.InvoiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.invoiceController.AddInvoiceItem(invoiceID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *InvoiceHandler) UpdateInvoiceItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req dto.InvoiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.invoiceController.UpdateInvoiceItem(itemID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *InvoiceHandler) RemoveInvoiceItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	err = h.invoiceController.RemoveInvoiceItem(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invoice item removed successfully"})
}

func (h *InvoiceHandler) AddInvoiceAdjustment(c *gin.Context) {
	invoiceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice id"})
		return
	}

	var req dto.InvoiceAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.invoiceController.AddInvoiceAdjustment(invoiceID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *InvoiceHandler) UpdateInvoiceAdjustment(c *gin.Context) {
	adjustmentID, err := strconv.ParseUint(c.Param("adjustmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adjustment id"})
		return
	}

	var req dto.InvoiceAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.invoiceController.UpdateInvoiceAdjustment(adjustmentID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *InvoiceHandler) RemoveInvoiceAdjustment(c *gin.Context) {
	adjustmentID, err := strconv.ParseUint(c.Param("adjustmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adjustment id"})
		return
	}

	err = h.invoiceController.RemoveInvoiceAdjustment(adjustmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "invoice adjustment removed successfully"})
}
