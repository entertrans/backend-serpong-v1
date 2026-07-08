// controller/cbt/assessment_crud.go
package cbt

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

// ============================================
// ASSESSMENT CRUD
// ============================================

func (ctrl *cbtController) CreateAssessment(c *gin.Context) {
	var req dto.CreateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract userID dari context
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - user not found"})
		return
	}

	createdBy, err := convertToInt(userIDValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var activeFrom *time.Time
	var activeUntil *time.Time

	// Untuk UB, waktu tidak wajib
	if req.Type == "UB" {
		// Jika UB, active_from dan active_until bisa null
		if req.ActiveFrom != nil && *req.ActiveFrom != "" {
			parsed, err := time.Parse("2006-01-02T15:04", *req.ActiveFrom)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid active_from format"})
				return
			}
			activeFrom = &parsed
		}

		if req.ActiveUntil != nil && *req.ActiveUntil != "" {
			parsed, err := time.Parse("2006-01-02T15:04", *req.ActiveUntil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid active_until format"})
				return
			}
			activeUntil = &parsed
		}
	} else {
		// Untuk tipe lain, waktu wajib diisi
		if req.ActiveFrom == nil || *req.ActiveFrom == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "active_from is required for this assessment type"})
			return
		}
		if req.ActiveUntil == nil || *req.ActiveUntil == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "active_until is required for this assessment type"})
			return
		}

		parsedFrom, err := time.Parse("2006-01-02T15:04", *req.ActiveFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid active_from format"})
			return
		}
		activeFrom = &parsedFrom

		parsedUntil, err := time.Parse("2006-01-02T15:04", *req.ActiveUntil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid active_until format"})
			return
		}
		activeUntil = &parsedUntil

		// Validasi tanggal untuk tipe non-UB
		if activeFrom.After(*activeUntil) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "active_from must be before active_until"})
			return
		}

		if activeFrom.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "active_from cannot be in the past"})
			return
		}
	}

	// Mulai transaction
	tx := ctrl.db.Begin()

	// Set status awal
	status := "draft"
	if req.Type != "UB" {
		status = "active" // Untuk tipe lain langsung active
	}

	// 1. Create assessment
	assessment := &model.ToAssessment{
		Title:                  req.Title,
		Description:            req.Description,
		Type:                   req.Type,
		AutoEnrollClassID:      req.AutoEnrollClassID,
		DurationMinutes:        req.DurationMinutes,
		IsRandomQuestion:       req.IsRandomQuestion,
		IsRandomOption:         req.IsRandomOption,
		TotalQuestionDisplayed: req.TotalQuestionDisplayed,
		PassingScore:           req.PassingScore,
		ActiveFrom:             activeFrom,
		ActiveUntil:            activeUntil,
		Status:                 status,
		Instruction:            req.Instruction,
		CreatedBy:              createdBy,
	}

	if err := tx.Create(assessment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Create assessment questions (opsional)
	if len(req.Questions) > 0 {
		if err := createAssessmentQuestions(tx, assessment.AssessmentID, req.Questions); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// Load assessment with questions for response
	response := buildAssessmentResponse(ctrl.db, assessment)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Assessment berhasil dibuat",
		"data":    response,
	})
}

// Endpoint untuk membuka/menutup UB dengan update tanggal otomatis
func (ctrl *cbtController) UpdateAssessmentStatus(c *gin.Context) {
	// Ambil ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	var req dto.UpdateAssessmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah assessment ada
	var assessment model.ToAssessment
	if err := ctrl.db.First(&assessment, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}

	// Hanya untuk tipe UB
	if assessment.Type != "UB" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status update only allowed for UB type"})
		return
	}

	now := time.Now()
	updates := make(map[string]interface{})

	switch req.Status {
	case "active":
		// MEMBUKA/MENGAKTIFKAN UB (bisa dilakukan berkali-kali)
		// - Set active_from = sekarang
		// - Set active_until = null (tidak terbatas sampai ditutup)
		// - Ubah status ke "active"
		updates["status"] = "active"
		updates["active_from"] = now
		updates["active_until"] = nil

		// Update database
		if err := ctrl.db.Model(&assessment).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "UB berhasil diaktifkan. Siswa sekarang dapat mengerjakan.",
			"data": gin.H{
				"assessment_id": assessment.AssessmentID,
				"status":        "active",
				"active_from":   now,
				"active_until":  nil,
				"note":          "UB akan aktif sampai ditutup manual oleh admin",
			},
		})
		return

	case "closed":
		// MENUTUP/MENONAKTIFKAN UB (bisa dilakukan berkali-kali)
		// - Set active_until = sekarang
		// - Ubah status ke "closed"
		updates["status"] = "closed"
		updates["active_until"] = now

		// Jika active_from null (belum pernah dibuka), set dengan sekarang juga
		if assessment.ActiveFrom == nil {
			updates["active_from"] = now
		}

		// Update database
		if err := ctrl.db.Model(&assessment).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Refresh assessment untuk mendapatkan data terbaru
		ctrl.db.First(&assessment, uint(id))

		c.JSON(http.StatusOK, gin.H{
			"message": "UB berhasil ditutup. Siswa tidak dapat mengerjakan lagi.",
			"data": gin.H{
				"assessment_id": assessment.AssessmentID,
				"status":        "closed",
				"active_from":   assessment.ActiveFrom,
				"active_until":  now,
				"note":          "UB dapat diaktifkan kembali kapan saja untuk tahun ajaran berikutnya",
			},
		})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
}

// Endpoint untuk update tanggal UB secara manual (jika diperlukan)
func (ctrl *cbtController) UpdateUBDate(c *gin.Context) {
	// Ambil ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	var req dto.UpdateUBDateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah assessment ada dan tipenya UB
	var assessment model.ToAssessment
	if err := ctrl.db.First(&assessment, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}

	if assessment.Type != "UB" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This endpoint only for UB type"})
		return
	}

	updates := make(map[string]interface{})

	if req.ActiveFrom != nil {
		updates["active_from"] = req.ActiveFrom
	}

	if req.ActiveUntil != nil {
		updates["active_until"] = req.ActiveUntil
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Update database
	if err := ctrl.db.Model(&assessment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload assessment
	ctrl.db.First(&assessment, uint(id))

	c.JSON(http.StatusOK, gin.H{
		"message": "Tanggal UB berhasil diupdate",
		"data": gin.H{
			"assessment_id": assessment.AssessmentID,
			"active_from":   assessment.ActiveFrom,
			"active_until":  assessment.ActiveUntil,
			"status":        assessment.Status,
		},
	})
}

func (ctrl *cbtController) GetAssessments(c *gin.Context) {
	search := c.Query("search")
	assessmentType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var assessments []model.ToAssessment
	var total int64

	query := ctrl.db.Model(&model.ToAssessment{}).
		Where("deleted_at IS NULL").
		Preload("Creator")

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	if assessmentType != "" {
		query = query.Where("type = ?", assessmentType)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&assessments).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]dto.AssessmentResponse, len(assessments))
	for i, assessment := range assessments {
		result[i] = buildAssessmentResponse(ctrl.db, &assessment)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (ctrl *cbtController) GetAssessmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	var assessment model.ToAssessment
	err = ctrl.db.
		Preload("Creator").
		Where("assessment_id = ? AND deleted_at IS NULL", uint(id)).
		First(&assessment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "assessment tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	response := buildAssessmentResponse(ctrl.db, &assessment)

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (ctrl *cbtController) UpdateAssessment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	var req dto.UpdateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assessment model.ToAssessment
	err = ctrl.db.Where("assessment_id = ? AND deleted_at IS NULL", uint(id)).
		First(&assessment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "assessment tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Mulai transaction
	tx := ctrl.db.Begin()

	// Update assessment fields
	if req.Title != nil {
		assessment.Title = *req.Title
	}
	if req.Description != nil {
		assessment.Description = req.Description
	}
	if req.Type != nil {
		assessment.Type = *req.Type
	}
	if req.AutoEnrollClassID != nil {
		assessment.AutoEnrollClassID = req.AutoEnrollClassID
	}
	if req.DurationMinutes != nil {
		assessment.DurationMinutes = req.DurationMinutes
	}
	if req.IsRandomQuestion != nil {
		assessment.IsRandomQuestion = *req.IsRandomQuestion
	}
	if req.IsRandomOption != nil {
		assessment.IsRandomOption = *req.IsRandomOption
	}
	if req.TotalQuestionDisplayed != nil {
		assessment.TotalQuestionDisplayed = req.TotalQuestionDisplayed
	}
	if req.PassingScore != nil {
		assessment.PassingScore = *req.PassingScore
	}

	// Handle ActiveFrom (bisa set ke null)
	if req.ActiveFrom != nil {
		assessment.ActiveFrom = *req.ActiveFrom
	}

	// Handle ActiveUntil (bisa set ke null)
	if req.ActiveUntil != nil {
		assessment.ActiveUntil = *req.ActiveUntil
	}

	if req.Instruction != nil {
		assessment.Instruction = req.Instruction
	}

	if err := tx.Save(&assessment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update questions if provided
	if req.Questions != nil {
		// Delete existing questions
		if err := tx.Where("assessment_id = ?", assessment.AssessmentID).Delete(&model.ToAssessmentQuestion{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete old questions: " + err.Error()})
			return
		}

		// Create new questions
		if err := createAssessmentQuestions(tx, assessment.AssessmentID, req.Questions); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	response := buildAssessmentResponse(ctrl.db, &assessment)

	c.JSON(http.StatusOK, gin.H{
		"message": "Assessment berhasil diupdate",
		"data":    response,
	})
}

func (ctrl *cbtController) DeleteAssessment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	result := ctrl.db.Model(&model.ToAssessment{}).
		Where("assessment_id = ?", uint(id)).
		Update("deleted_at", gorm.Expr("NOW()"))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment berhasil dihapus"})
}

// ============================================
// ENROL STUDENT
// ============================================

// EnrollStudents - Mendaftarkan siswa ke assessment
func (cc *cbtController) EnrollStudents(c *gin.Context) {
	assessmentIDStr := c.Param("id")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID assessment tidak valid"})
		return
	}

	// Validasi assessment exists
	var assessment model.ToAssessment
	if err := cc.db.First(&assessment, assessmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assessment tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa assessment"})
		return
	}

	// Parse request body
	var req dto.EnrollStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid: " + err.Error()})
		return
	}

	if len(req.StudentIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimal pilih 1 siswa"})
		return
	}

	// Ambil daftar student_id yang sudah terdaftar
	var existingStudentIDs []int
	cc.db.Model(&model.ToAssessmentStudent{}).
		Where("assessment_id = ?", assessmentID).
		Pluck("student_id", &existingStudentIDs)

	// Buat map untuk quick lookup
	existingMap := make(map[int]bool)
	for _, id := range existingStudentIDs {
		existingMap[id] = true
	}

	// Filter siswa yang belum terdaftar
	var newStudents []model.ToAssessmentStudent
	var duplicateCount int

	for _, studentID := range req.StudentIDs {
		if existingMap[studentID] {
			duplicateCount++
			continue
		}

		newStudents = append(newStudents, model.ToAssessmentStudent{
			AssessmentID: uint(assessmentID),
			StudentID:    studentID,
			IsActive:     true,
		})
	}

	if len(newStudents) == 0 {
		message := "Semua siswa sudah terdaftar"
		if duplicateCount > 0 {
			message = "Semua siswa sudah terdaftar (ada " + strconv.Itoa(duplicateCount) + " siswa duplikat)"
		}
		c.JSON(http.StatusOK, dto.EnrollStudentResponse{
			AssessmentID:      uint(assessmentID),
			TotalEnrolled:     0,
			FailedEnrollments: duplicateCount,
			Message:           message,
		})
		return
	}

	// Bulk insert siswa baru
	if err := cc.db.Create(&newStudents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan siswa: " + err.Error()})
		return
	}

	response := dto.EnrollStudentResponse{
		AssessmentID:      uint(assessmentID),
		TotalEnrolled:     len(newStudents),
		FailedEnrollments: duplicateCount,
		Message:           "Berhasil mendaftarkan " + strconv.Itoa(len(newStudents)) + " siswa",
	}

	if duplicateCount > 0 {
		response.Message += ", " + strconv.Itoa(duplicateCount) + " siswa sudah terdaftar"
	}

	c.JSON(http.StatusOK, response)
}

// GetAssessmentStudents - Mendapatkan daftar siswa yang terdaftar di assessment
func (cc *cbtController) GetAssessmentStudents(c *gin.Context) {
	assessmentIDStr := c.Param("id")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID assessment tidak valid"})
		return
	}

	// Query students dengan relasi ke tbl_siswa dan kelas
	var enrolledStudents []model.ToAssessmentStudent
	if err := cc.db.
		Preload("Siswa").
		Preload("Siswa.Kelas").
		Where("assessment_id = ? AND is_active = ?", assessmentID, true).
		Find(&enrolledStudents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	// Format response
	var studentList []dto.AssessmentStudentResponse
	for _, es := range enrolledStudents {
		student := dto.AssessmentStudentResponse{
			AssessmentStudentID: es.AssessmentStudentID,
			SiswaID:             es.StudentID,
			IsActive:            es.IsActive,
			EnrolledAt:          es.EnrolledAt.Format("2006-01-02T15:04:05Z07:00"),
			AttemptStatus:       "not_started", // Default
		}

		if es.Siswa != nil {
			// 🔥 Handle pointer fields dengan aman
			if es.Siswa.SiswaNama != nil {
				student.SiswaNama = *es.Siswa.SiswaNama
			}
			if es.Siswa.SiswaNIS != nil {
				student.SiswaNIS = *es.Siswa.SiswaNIS
			}
			if es.Siswa.SiswaNISN != nil {
				student.SiswaNISN = *es.Siswa.SiswaNISN
			}

			// 🔥 Ambil nama kelas (langsung assign string, bukan pointer)
			if es.Siswa.Kelas.KelasId != 0 {
				student.Kelas = es.Siswa.Kelas.KelasNama // ✅ Langsung assign string
			} else {
				student.Kelas = "Tidak ada kelas"
			}
		}

		studentList = append(studentList, student)
	}

	c.JSON(http.StatusOK, gin.H{
		"assessment_id": assessmentID,
		"total":         len(studentList),
		"students":      studentList,
	})
}

// Helper function untuk pointer string
func getStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// RemoveStudentFromAssessment - Menghapus siswa dari assessment (soft delete / set is_active false)
func (cc *cbtController) RemoveStudentFromAssessment(c *gin.Context) {
	assessmentIDStr := c.Param("id")
	studentIDStr := c.Param("student_id")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID assessment tidak valid"})
		return
	}

	studentID, err := strconv.ParseInt(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID siswa tidak valid"})
		return
	}

	// Cari data enrollment
	var enrollment model.ToAssessmentStudent
	if err := cc.db.
		Where("assessment_id = ? AND student_id = ?", assessmentID, studentID).
		First(&enrollment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak terdaftar di assessment ini"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa data"})
		return
	}

	// Soft delete dengan set is_active = false
	if err := cc.db.Model(&enrollment).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus siswa dari assessment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Siswa berhasil dihapus dari assessment",
		"assessment_id": assessmentID,
		"student_id":    studentID,
	})
}

// ============================================
// HELPER FUNCTIONS FOR ASSESSMENT
// ============================================

func createAssessmentQuestions(tx *gorm.DB, assessmentID uint, questions []dto.AssessmentQuestion) error {
	if len(questions) == 0 {
		return nil
	}

	assessmentQuestions := make([]model.ToAssessmentQuestion, len(questions))
	for i, q := range questions {
		assessmentQuestions[i] = model.ToAssessmentQuestion{
			AssessmentID: assessmentID,
			QuestionID:   uint(q.QuestionID),
			Weight:       q.Weight,
			FixedOrder:   q.FixedOrder,
		}
	}
	return tx.Create(&assessmentQuestions).Error
}

func buildAssessmentResponse(db *gorm.DB, assessment *model.ToAssessment) dto.AssessmentResponse {
	// Hitung total questions
	var totalQuestions int64
	db.Model(&model.ToAssessmentQuestion{}).
		Where("assessment_id = ?", assessment.AssessmentID).
		Count(&totalQuestions)

	// Load assessment questions with question details
	var assessmentQuestions []model.ToAssessmentQuestion
	db.Where("assessment_id = ?", assessment.AssessmentID).
		Preload("Question").
		Find(&assessmentQuestions)

	questions := make([]dto.AssessmentQuestionDetail, len(assessmentQuestions))
	for i, aq := range assessmentQuestions {
		var questionDetail dto.QuestionDetailResponse
		if aq.Question != nil {
			questionDetail = buildQuestionDetailResponse(db, aq.Question)
		}

		questions[i] = dto.AssessmentQuestionDetail{
			AssessmentQuestionID: aq.AssessmentQuestionID,
			QuestionID:           aq.QuestionID,
			Weight:               aq.Weight,
			FixedOrder:           aq.FixedOrder,
			Question:             questionDetail,
		}
	}

	createdByName := ""
	if assessment.Creator != nil {
		createdByName = assessment.Creator.Name
	}

	// Get class name if auto_enroll_class_id exists
	var className string
	var classNamePtr *string
	if assessment.AutoEnrollClassID != nil && *assessment.AutoEnrollClassID > 0 {
		var kelas model.Kelas
		if err := db.Table("tbl_kelas").
			Where("kelas_id = ?", *assessment.AutoEnrollClassID).
			First(&kelas).Error; err == nil {
			className = kelas.KelasNama
			classNamePtr = &className
		}
	}

	return dto.AssessmentResponse{
		AssessmentID:           assessment.AssessmentID,
		Title:                  assessment.Title,
		Description:            assessment.Description,
		Type:                   assessment.Type,
		AutoEnrollClassID:      assessment.AutoEnrollClassID,
		AutoEnrollClassName:    classNamePtr, // Tambahkan field ini
		DurationMinutes:        assessment.DurationMinutes,
		IsRandomQuestion:       assessment.IsRandomQuestion,
		IsRandomOption:         assessment.IsRandomOption,
		TotalQuestionDisplayed: assessment.TotalQuestionDisplayed,
		PassingScore:           assessment.PassingScore,
		ActiveFrom:             assessment.ActiveFrom,
		ActiveUntil:            assessment.ActiveUntil,
		Instruction:            assessment.Instruction,
		TotalQuestions:         int(totalQuestions),
		CreatedBy:              assessment.CreatedBy,
		CreatedByName:          createdByName,
		CreatedAt:              assessment.CreatedAt,
		UpdatedAt:              assessment.UpdatedAt,
		Questions:              questions,
	}
}
