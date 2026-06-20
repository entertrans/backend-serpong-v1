// model/to_assessment.go
package model

import (
	"time"
)

type ToAssessment struct {
	AssessmentID           uint       `gorm:"primaryKey;autoIncrement" json:"assessment_id"`
	Title                  string     `gorm:"type:varchar(255);not null" json:"title"`
	Description            *string    `gorm:"type:text" json:"description,omitempty"`
	Type                   string     `gorm:"type:enum('UB','QUIZ','UTS','UAS','TASK');not null" json:"type"`
	AutoEnrollClassID      *int       `gorm:"column:auto_enroll_class_id" json:"auto_enroll_class_id,omitempty"`
	DurationMinutes        *int       `gorm:"column:duration_minutes" json:"duration_minutes,omitempty"`
	IsRandomQuestion       bool       `gorm:"default:false" json:"is_random_question"`
	IsRandomOption         bool       `gorm:"default:false" json:"is_random_option"`
	TotalQuestionDisplayed *int       `gorm:"column:total_question_displayed" json:"total_question_displayed,omitempty"`
	PassingScore           float64    `gorm:"type:decimal(5,2);default:70.00" json:"passing_score"`
	ActiveFrom             *time.Time `gorm:"index"`                            // Ubah jadi pointer (bisa null)
	ActiveUntil            *time.Time `gorm:"index"`                            // Ubah jadi pointer (bisa null)
	Status                 string     `gorm:"type:varchar(20);default:'draft'"` // draft, active, closed
	Instruction            *string    `gorm:"type:text" json:"instruction,omitempty"`
	CreatedBy              int        `gorm:"not null" json:"created_by"`
	CreatedAt              time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relations

	Questions []ToAssessmentQuestion `gorm:"foreignKey:AssessmentID" json:"questions,omitempty"`
	Creator   *User                  `gorm:"foreignKey:CreatedBy;references:ID" json:"creator,omitempty"`
}

func (ToAssessment) TableName() string { return "to_assessment" }

// ============================================
// ASSESSMENT STUDENT (Peserta)
// ============================================
type ToAssessmentStudent struct {
	AssessmentStudentID uint      `gorm:"primaryKey;autoIncrement" json:"assessment_student_id"`
	AssessmentID        uint      `gorm:"not null;index" json:"assessment_id"`
	StudentID           int       `gorm:"not null;index" json:"student_id"`
	EnrolledAt          time.Time `gorm:"autoCreateTime" json:"enrolled_at"`
	IsActive            bool      `gorm:"default:true" json:"is_active"`

	// Relasi ke tabel siswa
	Siswa *Siswa `gorm:"foreignKey:StudentID;references:SiswaID" json:"siswa,omitempty"`
}

func (ToAssessmentStudent) TableName() string { return "to_assessment_student" }
