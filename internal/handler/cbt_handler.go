// handler/cbt_handler.go
package handler

import (
	"github.com/entertrans/backend-bogor.git/internal/controller/cbt"
	"github.com/gin-gonic/gin"
)

type CBTHandler struct {
	controller cbt.CBTController
}

func NewCBTHandler(controller cbt.CBTController) *CBTHandler {
	return &CBTHandler{controller: controller}
}

// Bank handlers
func (h *CBTHandler) CreateBank(c *gin.Context) {
	h.controller.CreateBank(c)
}

func (h *CBTHandler) GetBanks(c *gin.Context) {
	h.controller.GetBanks(c)
}

func (h *CBTHandler) GetBankByID(c *gin.Context) {
	h.controller.GetBankByID(c)
}

func (h *CBTHandler) GetBanksByKelas(c *gin.Context) {
	h.controller.GetBanksByKelas(c)
}

func (h *CBTHandler) UpdateBank(c *gin.Context) {
	h.controller.UpdateBank(c)
}

func (h *CBTHandler) DeleteBank(c *gin.Context) {
	h.controller.DeleteBank(c)
}

func (h *CBTHandler) GetQuestionsByBank(c *gin.Context) {
	h.controller.GetQuestionsByBank(c)
}

// Question handlers
func (h *CBTHandler) CreateQuestion(c *gin.Context) {
	h.controller.CreateQuestion(c)
}

func (h *CBTHandler) GetQuestions(c *gin.Context) {
	h.controller.GetQuestions(c)
}

func (h *CBTHandler) GetQuestionByID(c *gin.Context) {
	h.controller.GetQuestionByID(c)
}

func (h *CBTHandler) UpdateQuestion(c *gin.Context) {
	h.controller.UpdateQuestion(c)
}

func (h *CBTHandler) DeleteQuestion(c *gin.Context) {
	h.controller.DeleteQuestion(c)
}

// Assessment handlers
func (h *CBTHandler) CreateAssessment(c *gin.Context) {
	h.controller.CreateAssessment(c)
}

func (h *CBTHandler) GetAssessments(c *gin.Context) {
	h.controller.GetAssessments(c)
}

func (h *CBTHandler) GetAssessmentByID(c *gin.Context) {
	h.controller.GetAssessmentByID(c)
}

func (h *CBTHandler) UpdateAssessment(c *gin.Context) {
	h.controller.UpdateAssessment(c)
}

func (h *CBTHandler) DeleteAssessment(c *gin.Context) {
	h.controller.DeleteAssessment(c)
}

// UB (Ujian Bulanan) specific handlers
func (h *CBTHandler) UpdateAssessmentStatus(c *gin.Context) {
	h.controller.UpdateAssessmentStatus(c)
}

func (h *CBTHandler) UpdateUBDate(c *gin.Context) {
	h.controller.UpdateUBDate(c)
}

// Assessment Question handlers
func (h *CBTHandler) AddQuestionsToAssessment(c *gin.Context) {
	h.controller.AddQuestionsToAssessment(c)
}

func (h *CBTHandler) GetAssessmentQuestions(c *gin.Context) {
	h.controller.GetAssessmentQuestions(c)
}

func (h *CBTHandler) UpdateQuestionWeight(c *gin.Context) {
	h.controller.UpdateQuestionWeight(c)
}

func (h *CBTHandler) UpdateQuestionOrder(c *gin.Context) {
	h.controller.UpdateQuestionOrder(c)
}

func (h *CBTHandler) RemoveQuestionFromAssessment(c *gin.Context) {
	h.controller.RemoveQuestionFromAssessment(c)
}

func (h *CBTHandler) EnrollStudents(c *gin.Context) {
	h.controller.EnrollStudents(c)
}

func (h *CBTHandler) GetAssessmentStudents(c *gin.Context) {
	h.controller.GetAssessmentStudents(c)
}

func (h *CBTHandler) RemoveStudentFromAssessment(c *gin.Context) {
	h.controller.RemoveStudentFromAssessment(c)
}
