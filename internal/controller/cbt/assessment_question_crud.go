// controller/cbt/assessment_question_crud.go
package cbt

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

// ============================================
// ASSESSMENT QUESTION CRUD
// ============================================

// AddQuestionsToAssessment - Menambahkan soal ke assessment dengan bobot
func (ctrl *cbtController) AddQuestionsToAssessment(c *gin.Context) {
	// Ambil assessment ID dari URL
	assessmentIDStr := c.Param("id")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	// Parse request body
	var req dto.AddQuestionToAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah assessment exists
	var assessment model.ToAssessment
	if err := ctrl.db.Where("assessment_id = ? AND deleted_at IS NULL", uint(assessmentID)).First(&assessment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assessment tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Cek duplikasi soal
	var existingQuestions []model.ToAssessmentQuestion
	ctrl.db.Where("assessment_id = ?", uint(assessmentID)).Find(&existingQuestions)
	existingQuestionMap := make(map[uint]bool)
	for _, eq := range existingQuestions {
		existingQuestionMap[eq.QuestionID] = true
	}

	// Filter soal yang belum ada
	var newQuestions []dto.AssessmentQuestion
	var duplicateQuestions []uint
	for _, q := range req.Questions {
		if existingQuestionMap[uint(q.QuestionID)] {
			duplicateQuestions = append(duplicateQuestions, uint(q.QuestionID))
		} else {
			newQuestions = append(newQuestions, q)
		}
	}

	if len(duplicateQuestions) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":      "Beberapa soal sudah ada dalam assessment ini",
			"duplicates": duplicateQuestions,
		})
		return
	}

	if len(newQuestions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada soal baru yang ditambahkan"})
		return
	}

	// Mulai transaction
	tx := ctrl.db.Begin()

	// Simpan ke to_assessment_question
	assessmentQuestions := make([]model.ToAssessmentQuestion, len(newQuestions))
	for i, q := range newQuestions {
		assessmentQuestions[i] = model.ToAssessmentQuestion{
			AssessmentID: uint(assessmentID),
			QuestionID:   uint(q.QuestionID),
			Weight:       q.Weight,
			FixedOrder:   q.FixedOrder,
		}
	}

	if err := tx.Create(&assessmentQuestions).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan soal: " + err.Error()})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal commit transaction"})
		return
	}

	// Load data untuk response
	var addedQuestions []model.ToAssessmentQuestion
	ctrl.db.Where("assessment_id = ?", uint(assessmentID)).
		Preload("Question").
		Find(&addedQuestions)

	response := buildAssessmentQuestionResponse(ctrl.db, addedQuestions)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Soal berhasil ditambahkan ke assessment",
		"data":    response,
	})
}

// GetAssessmentQuestions - Mendapatkan semua soal dalam assessment dengan detail lengkap
func (ctrl *cbtController) GetAssessmentQuestions(c *gin.Context) {
	assessmentIDStr := c.Param("id")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	// Cek apakah assessment exists
	var assessment model.ToAssessment
	if err := ctrl.db.Where("assessment_id = ? AND deleted_at IS NULL", uint(assessmentID)).First(&assessment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assessment tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// AMBIL SEMUA ASSESSMENT QUESTIONS dengan JOIN ke to_question
	var assessmentQuestions []model.ToAssessmentQuestion

	// Gunakan Joins untuk memastikan data diambil dengan benar
	err = ctrl.db.
		Table("to_assessment_question").
		Select("to_assessment_question.*").
		Where("to_assessment_question.assessment_id = ?", uint(assessmentID)).
		Order("to_assessment_question.fixed_order ASC, to_assessment_question.assessment_question_id ASC").
		Find(&assessmentQuestions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response dengan manual query untuk setiap question
	result := make([]dto.AssessmentQuestionDetailResponse, len(assessmentQuestions))
	for i, aq := range assessmentQuestions {
		// Ambil question secara manual dengan query terpisah
		var question model.ToQuestion
		if err := ctrl.db.Where("question_id = ? AND deleted_at IS NULL", aq.QuestionID).First(&question).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Log warning tapi tetap lanjutkan
				fmt.Printf("Warning: Question with ID %d not found for assessment question %d\n", aq.QuestionID, aq.AssessmentQuestionID)
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Gunakan buildQuestionDetailResponse yang sudah ada
		questionDetail := buildQuestionDetailResponse(ctrl.db, &question)

		result[i] = dto.AssessmentQuestionDetailResponse{
			AssessmentQuestionID: aq.AssessmentQuestionID,
			AssessmentID:         aq.AssessmentID,
			QuestionID:           aq.QuestionID,
			Weight:               aq.Weight,
			FixedOrder:           aq.FixedOrder,
			Question:             convertToQuestionFullDetail(questionDetail),
		}
	}

	// Filter out any nil entries
	filteredResult := make([]dto.AssessmentQuestionDetailResponse, 0)
	for _, r := range result {
		if r.QuestionID != 0 {
			filteredResult = append(filteredResult, r)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  filteredResult,
		"total": len(filteredResult),
		"assessment": gin.H{
			"id":    assessment.AssessmentID,
			"title": assessment.Title,
			"type":  assessment.Type,
		},
	})
}

// UpdateQuestionWeight - Update bobot soal dalam assessment
func (ctrl *cbtController) UpdateQuestionWeight(c *gin.Context) {
	// Ambil parameter
	assessmentIDStr := c.Param("id")
	questionIDStr := c.Param("question_id")
	// Tambahkan log
	// println("=== DEBUG UpdateQuestionWeight ===")
	// println("assessmentIDStr:", assessmentIDStr)
	// println("questionIDStr:", questionIDStr)

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}

	// println("Parsed assessmentID:", assessmentID)
	// println("Parsed questionID:", questionID)

	// Parse request
	var req dto.UpdateQuestionWeightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari record assessment_question
	var assessmentQuestion model.ToAssessmentQuestion
	err = ctrl.db.Where("assessment_id = ? AND question_id = ?", uint(assessmentID), uint(questionID)).
		First(&assessmentQuestion).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan dalam assessment ini"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// println("Record found! AssessmentQuestionID:", assessmentQuestion.AssessmentQuestionID)

	// Update weight
	assessmentQuestion.Weight = req.Weight
	if err := ctrl.db.Save(&assessmentQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update bobot: " + err.Error()})
		return
	}

	// Load ulang dengan relasi
	ctrl.db.Preload("Question").First(&assessmentQuestion, assessmentQuestion.AssessmentQuestionID)

	// Gunakan buildQuestionDetailResponse yang sudah ada
	questionDetail := buildQuestionDetailResponse(ctrl.db, assessmentQuestion.Question)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bobot soal berhasil diupdate",
		"data": dto.AssessmentQuestionResponse{
			AssessmentQuestionID: assessmentQuestion.AssessmentQuestionID,
			AssessmentID:         assessmentQuestion.AssessmentID,
			QuestionID:           assessmentQuestion.QuestionID,
			Weight:               assessmentQuestion.Weight,
			FixedOrder:           assessmentQuestion.FixedOrder,
			Question:             questionDetail,
		},
	})
}

// UpdateQuestionOrder - Update urutan tetap soal dalam assessment
func (ctrl *cbtController) UpdateQuestionOrder(c *gin.Context) {
	// Ambil parameter
	assessmentIDStr := c.Param("id")
	questionIDStr := c.Param("question_id")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}

	// Parse request
	var req dto.UpdateQuestionOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari record assessment_question
	var assessmentQuestion model.ToAssessmentQuestion
	err = ctrl.db.Where("assessment_id = ? AND question_id = ?", uint(assessmentID), uint(questionID)).
		First(&assessmentQuestion).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan dalam assessment ini"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Update fixed_order
	assessmentQuestion.FixedOrder = req.FixedOrder
	if err := ctrl.db.Save(&assessmentQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update urutan: " + err.Error()})
		return
	}

	message := "Urutan soal direset ke random"
	if req.FixedOrder != nil {
		message = "Urutan soal ditetapkan menjadi urutan ke-" + strconv.Itoa(*req.FixedOrder)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data": gin.H{
			"question_id": questionID,
			"fixed_order": req.FixedOrder,
		},
	})
}

// RemoveQuestionFromAssessment - Hapus soal dari assessment
func (ctrl *cbtController) RemoveQuestionFromAssessment(c *gin.Context) {
	// Ambil parameter
	assessmentIDStr := c.Param("id")
	questionIDStr := c.Param("question_id")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment id"})
		return
	}

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}

	// Hapus record
	result := ctrl.db.Where("assessment_id = ? AND question_id = ?", uint(assessmentID), uint(questionID)).
		Delete(&model.ToAssessmentQuestion{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan dalam assessment ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Soal berhasil dihapus dari assessment",
	})
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// buildAssessmentQuestionResponse - untuk response sederhana (tanpa detail lengkap)
func buildAssessmentQuestionResponse(db *gorm.DB, assessmentQuestions []model.ToAssessmentQuestion) []dto.AssessmentQuestionResponse {
	response := make([]dto.AssessmentQuestionResponse, len(assessmentQuestions))
	for i, aq := range assessmentQuestions {
		var questionDetail dto.QuestionDetailResponse
		if aq.Question != nil {
			questionDetail = buildQuestionDetailResponse(db, aq.Question)
		}

		response[i] = dto.AssessmentQuestionResponse{
			AssessmentQuestionID: aq.AssessmentQuestionID,
			AssessmentID:         aq.AssessmentID,
			QuestionID:           aq.QuestionID,
			Weight:               aq.Weight,
			FixedOrder:           aq.FixedOrder,
			Question:             questionDetail,
		}
	}
	return response
}

// convertToQuestionFullDetail - konversi dari QuestionDetailResponse ke QuestionFullDetailResponse
func convertToQuestionFullDetail(q dto.QuestionDetailResponse) dto.QuestionFullDetailResponse {
	// Handle pointer BankID
	var bankID uint
	if q.BankID != nil {
		bankID = *q.BankID
	}

	// Handle pointer Explanation
	explanation := ""
	if q.Explanation != nil {
		explanation = *q.Explanation
	}

	return dto.QuestionFullDetailResponse{
		QuestionID:    q.QuestionID,
		BankID:        bankID,
		QuestionType:  q.QuestionType,
		QuestionText:  q.QuestionText,
		Explanation:   explanation,
		CreatedBy:     uint64(q.CreatedBy),
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
		Options:       convertOptions(q.Options),
		ValidAnswers:  q.ValidAnswers,
		MatchingPairs: q.MatchingPairs,
	}
}

// convertOptions - konversi dari OptionResponse ke QuestionOptionResponse
func convertOptions(options []dto.OptionResponse) []dto.QuestionOptionResponse {
	if options == nil {
		return nil
	}

	result := make([]dto.QuestionOptionResponse, len(options))
	for i, opt := range options {
		result[i] = dto.QuestionOptionResponse{
			OptionID:    opt.OptionID,
			OptionText:  opt.OptionText,
			IsCorrect:   opt.IsCorrect,
			OptionOrder: opt.SortOrder,
		}
	}
	return result
}
