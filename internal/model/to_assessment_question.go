// model/to_assessment_question.go
package model

type ToAssessmentQuestion struct {
	AssessmentQuestionID uint `gorm:"primaryKey;autoIncrement" json:"assessment_question_id"`
	AssessmentID         uint `gorm:"not null;index" json:"assessment_id"`
	QuestionID           uint `gorm:"not null" json:"question_id"`
	Weight               int  `gorm:"default:1" json:"weight"`
	FixedOrder           *int `gorm:"column:fixed_order" json:"fixed_order,omitempty"`

	// Relations
	Assessment *ToAssessment `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
	Question   *ToQuestion   `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (ToAssessmentQuestion) TableName() string { return "to_assessment_question" }
