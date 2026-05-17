package finance

import (
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/controller/finance"
	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/gin-gonic/gin"
)

type FeeTemplateHandler struct {
	feeTemplateController finance.FeeTemplateController
}

func NewFeeTemplateHandler(feeTemplateController finance.FeeTemplateController) *FeeTemplateHandler {
	return &FeeTemplateHandler{
		feeTemplateController: feeTemplateController,
	}
}

func (h *FeeTemplateHandler) GetAllFeeTemplates(c *gin.Context) {
	isActiveStr := c.Query("is_active")
	var isActive *int8
	if isActiveStr != "" {
		val, err := strconv.ParseInt(isActiveStr, 10, 8)
		if err == nil {
			active := int8(val)
			isActive = &active
		}
	}

	result, err := h.feeTemplateController.GetAllFeeTemplates(isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FeeTemplateHandler) GetFeeTemplateByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.feeTemplateController.GetFeeTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FeeTemplateHandler) CreateFeeTemplate(c *gin.Context) {
	var req dto.FeeTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.feeTemplateController.CreateFeeTemplate(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result})
}

func (h *FeeTemplateHandler) UpdateFeeTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.FeeTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.feeTemplateController.UpdateFeeTemplate(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *FeeTemplateHandler) DeleteFeeTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.feeTemplateController.DeleteFeeTemplate(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "fee template deleted successfully"})
}
