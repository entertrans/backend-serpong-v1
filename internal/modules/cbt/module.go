// internal/modules/cbt/cbtmodule.go
package cbtmodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/entertrans/backend-bogor.git/internal/controller/cbt"
	"github.com/entertrans/backend-bogor.git/internal/handler"
	"github.com/entertrans/backend-bogor.git/internal/middleware"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	cbtController := cbt.NewCBTController(db)
	cbtHandler := handler.NewCBTHandler(cbtController)

	cbtGroup := rg.Group("/cbt")
	cbtGroup.Use(middleware.AuthMiddleware(cfg))
	{
		// Bank CRUD
		bankGroup := cbtGroup.Group("/bank")
		{
			bankGroup.POST("", cbtHandler.CreateBank)
			bankGroup.GET("", cbtHandler.GetBanks)
			bankGroup.GET("/:id", cbtHandler.GetBankByID)
			bankGroup.GET("/:id/questions", cbtHandler.GetQuestionsByBank)
			bankGroup.PUT("/:id", cbtHandler.UpdateBank)
			bankGroup.DELETE("/:id", cbtHandler.DeleteBank)
			bankGroup.GET("/kelas/:kelas_id", cbtHandler.GetBanksByKelas)
		}

		// Question CRUD
		questionGroup := cbtGroup.Group("/question")
		{
			questionGroup.POST("", cbtHandler.CreateQuestion)
			questionGroup.GET("", cbtHandler.GetQuestions)
			questionGroup.GET("/:id", cbtHandler.GetQuestionByID)
			questionGroup.PUT("/:id", cbtHandler.UpdateQuestion)
			questionGroup.DELETE("/:id", cbtHandler.DeleteQuestion)
		}

		// Assessment CRUD
		assessmentGroup := cbtGroup.Group("/assessment")
		{
			assessmentGroup.POST("", cbtHandler.CreateAssessment)
			assessmentGroup.GET("", cbtHandler.GetAssessments)
			assessmentGroup.GET("/:id", cbtHandler.GetAssessmentByID)
			assessmentGroup.PUT("/:id", cbtHandler.UpdateAssessment)
			assessmentGroup.DELETE("/:id", cbtHandler.DeleteAssessment)

			// Status & Date
			assessmentGroup.PATCH("/:id/status", cbtHandler.UpdateAssessmentStatus)
			assessmentGroup.PATCH("/:id/ub-date", cbtHandler.UpdateUBDate)

			// ✅ ENROLL STUDENT (TAMBAHKAN INI)
			assessmentGroup.POST("/:id/enroll", cbtHandler.EnrollStudents)
			assessmentGroup.GET("/:id/students", cbtHandler.GetAssessmentStudents)
			assessmentGroup.DELETE("/:id/students/:student_id", cbtHandler.RemoveStudentFromAssessment)

			// Manage questions in assessment
			assessmentGroup.POST("/:id/questions", cbtHandler.AddQuestionsToAssessment)
			assessmentGroup.GET("/:id/questions", cbtHandler.GetAssessmentQuestions)
			assessmentGroup.PUT("/:id/questions/:question_id/weight", cbtHandler.UpdateQuestionWeight)
			assessmentGroup.PUT("/:id/questions/:question_id/order", cbtHandler.UpdateQuestionOrder)
			assessmentGroup.DELETE("/:id/questions/:question_id", cbtHandler.RemoveQuestionFromAssessment)
		}
	}
}
