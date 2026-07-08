// controller/cbt/bank_crud.go
package cbt

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

// ============================================
// BANK CRUD
// ============================================

func (ctrl *cbtController) CreateBank(c *gin.Context) {
	var req dto.CreateBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gunakan "userID" (sesuai middleware)
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - user not found"})
		return
	}

	// Convert userID ke uint64
	var createdBy uint64
	switch v := userIDValue.(type) {
	case float64:
		createdBy = uint64(v)
	case int:
		createdBy = uint64(v)
	case uint:
		createdBy = uint64(v)
	case int64:
		createdBy = uint64(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	bank := &model.ToQuestionBank{
		BankName:    req.BankName,
		KdMapel:     req.KdMapel,
		KelasID:     req.KelasID,
		Description: &req.Description,
		CreatedBy:   createdBy,
	}

	if err := ctrl.db.Create(bank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load relasi untuk response
	ctrl.db.Preload("User").Preload("Mapel").Preload("Kelas").First(bank, bank.BankID)

	// Hitung total soal
	var totalQuestions int64
	ctrl.db.Model(&model.ToQuestion{}).Where("bank_id = ? AND deleted_at IS NULL", bank.BankID).Count(&totalQuestions)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bank soal berhasil dibuat",
		"data": dto.BankResponse{
			BankID:         bank.BankID,
			BankName:       bank.BankName,
			KdMapel:        bank.KdMapel,
			MapelName:      getMapelName(bank.Mapel),
			KelasID:        bank.KelasID,
			KelasName:      getKelasName(bank.Kelas),
			Description:    bank.Description,
			TotalQuestions: totalQuestions,
			CreatedBy:      bank.CreatedBy,
			CreatedByName:  getUserName(bank.User),
			CreatedAt:      bank.CreatedAt,
			UpdatedAt:      bank.UpdatedAt,
		},
	})
}

// controller/cbt/cbt_controller.go

func (ctrl *cbtController) GetBanks(c *gin.Context) {
	// ✅ Ambil user info dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Parse query parameters
	search := c.Query("search")
	mapelID := c.Query("mapel_id")
	kelasID := c.Query("kelas_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var banks []model.ToQuestionBank
	var total int64

	// ✅ Build query dengan filter berdasarkan role
	query := ctrl.db.Model(&model.ToQuestionBank{}).
		Where("deleted_at IS NULL").
		Preload("User").
		Preload("Mapel").
		Preload("Kelas")

	// ✅ LOGIC FILTER BERDASARKAN ROLE
	role := userRole.(string)
	userIDUint := uint(userID.(float64))

	if role == "admin" {
		// Admin: lihat SEMUA bank soal
		// Tidak perlu filter tambahan
	} else if role == "guru" {
		// Guru: lihat bank soal yang dibuat oleh:
		// 1. Guru itu sendiri (created_by = userID)
		// 2. Admin (created_by = user dengan role 'admin')
		query = query.Where(
			"created_by = ? OR created_by IN (SELECT id FROM users WHERE role = 'admin' AND deleted_at IS NULL)",
			userIDUint,
		)
	} else {
		// Role lain: hanya lihat milik sendiri
		query = query.Where("created_by = ?", userIDUint)
	}

	// Filter search
	if search != "" {
		query = query.Where("bank_name LIKE ?", "%"+search+"%")
	}

	// Filter mapel
	if mapelID != "" {
		query = query.Where("kd_mapel = ?", mapelID)
	}

	// Filter kelas
	if kelasID != "" {
		query = query.Where("kelas_id = ?", kelasID)
	}

	// Count total
	query.Count(&total)

	// Get paginated data
	offset := (page - 1) * limit
	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&banks).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response
	result := make([]dto.BankResponse, len(banks))
	for i, bank := range banks {
		// Hitung jumlah soal di bank ini
		var totalQuestions int64
		ctrl.db.Model(&model.ToQuestion{}).
			Where("bank_id = ? AND deleted_at IS NULL", bank.BankID).
			Count(&totalQuestions)

		result[i] = dto.BankResponse{
			BankID:         bank.BankID,
			BankName:       bank.BankName,
			KdMapel:        bank.KdMapel,
			MapelName:      getMapelName(bank.Mapel),
			KelasID:        bank.KelasID,
			KelasName:      getKelasName(bank.Kelas),
			Description:    bank.Description,
			TotalQuestions: totalQuestions,
			CreatedBy:      bank.CreatedBy,
			CreatedByName:  getUserName(bank.User),
			CreatedAt:      bank.CreatedAt,
			UpdatedAt:      bank.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// controller/cbt/cbt_controller.go

func (ctrl *cbtController) GetBankByID(c *gin.Context) {
	// ✅ Ambil user info dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank id"})
		return
	}

	var bank model.ToQuestionBank

	// ✅ Build query dengan filter akses
	query := ctrl.db.
		Preload("User").
		Preload("Mapel").
		Preload("Kelas").
		Where("bank_id = ? AND deleted_at IS NULL", uint(id))

	// ✅ Cek akses berdasarkan role
	role := userRole.(string)
	userIDUint := uint(userID.(float64))

	if role != "admin" {
		if role == "guru" {
			// Guru: cek apakah bank milik sendiri atau milik admin
			query = query.Where(
				"created_by = ? OR created_by IN (SELECT id FROM users WHERE role = 'admin' AND deleted_at IS NULL)",
				userIDUint,
			)
		} else {
			// Role lain: hanya milik sendiri
			query = query.Where("created_by = ?", userIDUint)
		}
	}

	err = query.First(&bank).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "bank soal tidak ditemukan atau tidak memiliki akses"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Hitung jumlah soal di bank ini
	var totalQuestions int64
	ctrl.db.Model(&model.ToQuestion{}).
		Where("bank_id = ? AND deleted_at IS NULL", bank.BankID).
		Count(&totalQuestions)

	// Hitung statistik tipe soal
	type QuestionTypeCount struct {
		QuestionType string
		Count        int64
	}
	var typeStats []QuestionTypeCount
	ctrl.db.Model(&model.ToQuestion{}).
		Select("question_type, count(*) as count").
		Where("bank_id = ? AND deleted_at IS NULL", bank.BankID).
		Group("question_type").
		Scan(&typeStats)

	questionTypes := make(map[string]int64)
	for _, stat := range typeStats {
		questionTypes[stat.QuestionType] = stat.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.BankDetailResponse{
			BankID:         bank.BankID,
			BankName:       bank.BankName,
			KdMapel:        bank.KdMapel,
			MapelName:      getMapelName(bank.Mapel),
			KelasID:        bank.KelasID,
			KelasName:      getKelasName(bank.Kelas),
			Description:    bank.Description,
			TotalQuestions: totalQuestions,
			QuestionTypes:  questionTypes,
			CreatedBy:      bank.CreatedBy,
			CreatedByName:  getUserName(bank.User),
			CreatedAt:      bank.CreatedAt,
			UpdatedAt:      bank.UpdatedAt,
		},
	})
}

func (ctrl *cbtController) UpdateBank(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank id"})
		return
	}

	var req dto.UpdateBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bank model.ToQuestionBank
	err = ctrl.db.Where("bank_id = ? AND deleted_at IS NULL", uint(id)).
		First(&bank).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "bank soal tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if req.BankName != nil {
		bank.BankName = *req.BankName
	}
	if req.KdMapel != nil {
		bank.KdMapel = *req.KdMapel
	}
	if req.KelasID != nil {
		bank.KelasID = *req.KelasID
	}
	if req.Description != nil {
		bank.Description = req.Description
	}

	if err := ctrl.db.Save(&bank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.Preload("User").Preload("Mapel").Preload("Kelas").First(&bank, bank.BankID)

	var totalQuestions int64
	ctrl.db.Model(&model.ToQuestion{}).
		Where("bank_id = ? AND deleted_at IS NULL", bank.BankID).
		Count(&totalQuestions)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bank soal berhasil diupdate",
		"data": dto.BankResponse{
			BankID:         bank.BankID,
			BankName:       bank.BankName,
			KdMapel:        bank.KdMapel,
			MapelName:      getMapelName(bank.Mapel),
			KelasID:        bank.KelasID,
			KelasName:      getKelasName(bank.Kelas),
			Description:    bank.Description,
			TotalQuestions: totalQuestions,
			CreatedBy:      bank.CreatedBy,
			CreatedByName:  getUserName(bank.User),
			CreatedAt:      bank.CreatedAt,
			UpdatedAt:      bank.UpdatedAt,
		},
	})
}

func (ctrl *cbtController) DeleteBank(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank id"})
		return
	}

	var totalQuestions int64
	ctrl.db.Model(&model.ToQuestion{}).
		Where("bank_id = ? AND deleted_at IS NULL", uint(id)).
		Count(&totalQuestions)

	if totalQuestions > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tidak dapat menghapus bank yang masih memiliki soal. Hapus semua soal terlebih dahulu.",
		})
		return
	}

	result := ctrl.db.Model(&model.ToQuestionBank{}).
		Where("bank_id = ?", uint(id)).
		Update("deleted_at", gorm.Expr("NOW()"))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank soal tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank soal berhasil dihapus"})
}

func (ctrl *cbtController) GetBanksByKelas(c *gin.Context) {
	// ✅ Ambil user info dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Ambil parameter kelas_id dari URL
	kelasIDStr := c.Param("kelas_id")
	kelasID, err := strconv.ParseUint(kelasIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelas id"})
		return
	}

	// Optional: filter search
	search := c.Query("search")

	// Optional: pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var banks []model.ToQuestionBank
	var total int64

	// ✅ Build query dengan filter role
	query := ctrl.db.Model(&model.ToQuestionBank{}).
		Where("kelas_id = ? AND deleted_at IS NULL", uint(kelasID)).
		Preload("Mapel").
		Preload("User") // Preload user untuk mendapatkan created_by_name

	// ✅ LOGIC FILTER BERDASARKAN ROLE
	role := userRole.(string)
	userIDUint := uint(userID.(float64))

	if role == "admin" {
		// Admin: lihat SEMUA bank soal di kelas ini
		// Tidak perlu filter tambahan
	} else if role == "guru" {
		// Guru: lihat bank soal yang dibuat oleh:
		// 1. Guru itu sendiri (created_by = userID)
		// 2. Admin (created_by = user dengan role 'admin')
		query = query.Where(
			"created_by = ? OR created_by IN (SELECT id FROM users WHERE role = 'admin' AND deleted_at IS NULL)",
			userIDUint,
		)
	} else {
		// Role lain: hanya lihat milik sendiri
		query = query.Where("created_by = ?", userIDUint)
	}

	// Filter search
	if search != "" {
		query = query.Where("bank_name LIKE ?", "%"+search+"%")
	}

	// Hitung total
	query.Count(&total)

	// Execute query dengan pagination
	offset := (page - 1) * limit
	err = query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&banks).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format response
	result := make([]dto.BankByKelasResponse, len(banks))
	for i, bank := range banks {
		// Hitung total soal di bank ini
		var totalQuestions int64
		ctrl.db.Model(&model.ToQuestion{}).
			Where("bank_id = ? AND deleted_at IS NULL", bank.BankID).
			Count(&totalQuestions)

		result[i] = dto.BankByKelasResponse{
			BankID:         bank.BankID,
			BankName:       bank.BankName,
			KdMapel:        bank.KdMapel,
			MapelName:      getMapelName(bank.Mapel),
			Description:    bank.Description,
			TotalQuestions: totalQuestions,
			// CreatedBy:      bank.CreatedBy,
			// CreatedByName:  getUserName(bank.User),
			CreatedAt: bank.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     result,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"kelas_id": kelasID,
		"message":  "Berhasil mengambil data bank soal berdasarkan kelas",
	})
}
