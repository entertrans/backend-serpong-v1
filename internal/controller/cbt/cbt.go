// controller/cbt/cbt.go
package cbt

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CBTController interface {
	// Bank methods
	CreateBank(c *gin.Context)
	GetBanks(c *gin.Context)
	GetBankByID(c *gin.Context)
	UpdateBank(c *gin.Context)
	DeleteBank(c *gin.Context)
	GetQuestionsByBank(c *gin.Context)
	GetBanksByKelas(c *gin.Context)

	// Question methods
	CreateQuestion(c *gin.Context)
	GetQuestions(c *gin.Context)
	GetQuestionByID(c *gin.Context)
	UpdateQuestion(c *gin.Context)
	DeleteQuestion(c *gin.Context)

	// Assessment methods
	CreateAssessment(c *gin.Context)
	GetAssessments(c *gin.Context)
	GetAssessmentByID(c *gin.Context)
	UpdateAssessment(c *gin.Context)
	DeleteAssessment(c *gin.Context)

	// UB specific methods
	UpdateAssessmentStatus(c *gin.Context)
	UpdateUBDate(c *gin.Context)

	// Assessment Question methods (TAMBAHKAN INI)
	AddQuestionsToAssessment(c *gin.Context)
	GetAssessmentQuestions(c *gin.Context)
	UpdateQuestionWeight(c *gin.Context)
	UpdateQuestionOrder(c *gin.Context)
	RemoveQuestionFromAssessment(c *gin.Context)

	// Enroll Student methods (TAMBAHKAN INI)
	EnrollStudents(c *gin.Context)
	GetAssessmentStudents(c *gin.Context)
	RemoveStudentFromAssessment(c *gin.Context)
}

type cbtController struct {
	db *gorm.DB
}

func NewCBTController(db *gorm.DB) CBTController {
	return &cbtController{
		db: db,
	}
}
