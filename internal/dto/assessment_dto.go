// dto/assessment_dto.go
package dto

import "time"

// ========== ASSESSMENT DTO ==========

// Di file DTO
type CreateAssessmentRequest struct {
	Title                  string               `json:"title" binding:"required"`
	Description            *string              `json:"description,omitempty"`
	Type                   string               `json:"type" binding:"required,oneof=UB QUIZ UTS UAS TASK"`
	AutoEnrollClassID      *int                 `json:"auto_enroll_class_id,omitempty"`
	DurationMinutes        *int                 `json:"duration_minutes,omitempty"`
	IsRandomQuestion       bool                 `json:"is_random_question"`
	IsRandomOption         bool                 `json:"is_random_option"`
	TotalQuestionDisplayed *int                 `json:"total_question_displayed,omitempty"`
	PassingScore           float64              `json:"passing_score"`
	ActiveFrom             *string              `json:"active_from,omitempty"`  // Bisa null untuk UB
	ActiveUntil            *string              `json:"active_until,omitempty"` // Bisa null untuk UB
	Instruction            *string              `json:"instruction,omitempty"`
	Questions              []AssessmentQuestion `json:"questions"`
}

type UpdateAssessmentRequest struct {
	Title                  *string              `json:"title,omitempty"`
	Description            *string              `json:"description,omitempty"`
	Type                   *string              `json:"type,omitempty" binding:"omitempty,oneof=UB QUIZ UTS UAS TASK"`
	AutoEnrollClassID      *int                 `json:"auto_enroll_class_id,omitempty"`
	DurationMinutes        *int                 `json:"duration_minutes,omitempty"`
	IsRandomQuestion       *bool                `json:"is_random_question,omitempty"`
	IsRandomOption         *bool                `json:"is_random_option,omitempty"`
	TotalQuestionDisplayed *int                 `json:"total_question_displayed,omitempty"`
	PassingScore           *float64             `json:"passing_score,omitempty"`
	ActiveFrom             **time.Time          `json:"active_from,omitempty"`  // Pointer ke pointer untuk bisa set null
	ActiveUntil            **time.Time          `json:"active_until,omitempty"` // Pointer ke pointer untuk bisa set null
	Instruction            *string              `json:"instruction,omitempty"`
	Questions              []AssessmentQuestion `json:"questions,omitempty"`
}

type AssessmentQuestion struct {
	QuestionID int  `json:"question_id" binding:"required"`
	Weight     int  `json:"weight"`
	FixedOrder *int `json:"fixed_order,omitempty"`
}

type AssessmentResponse struct {
	AssessmentID           uint                       `json:"assessment_id"`
	Title                  string                     `json:"title"`
	Description            *string                    `json:"description,omitempty"`
	Type                   string                     `json:"type"`
	AutoEnrollClassID      *int                       `json:"auto_enroll_class_id,omitempty"`
	AutoEnrollClassName    *string                    `json:"auto_enroll_class_name,omitempty"`
	DurationMinutes        *int                       `json:"duration_minutes,omitempty"`
	IsRandomQuestion       bool                       `json:"is_random_question"`
	IsRandomOption         bool                       `json:"is_random_option"`
	TotalQuestionDisplayed *int                       `json:"total_question_displayed,omitempty"`
	PassingScore           float64                    `json:"passing_score"`
	ActiveFrom             *time.Time                 `json:"active_from"`  // Pointer
	ActiveUntil            *time.Time                 `json:"active_until"` // Pointer
	Instruction            *string                    `json:"instruction,omitempty"`
	TotalQuestions         int                        `json:"total_questions"`
	CreatedBy              int                        `json:"created_by"`
	CreatedByName          string                     `json:"created_by_name"`
	CreatedAt              time.Time                  `json:"created_at"`
	UpdatedAt              time.Time                  `json:"updated_at"`
	Questions              []AssessmentQuestionDetail `json:"questions,omitempty"`
}

type AssessmentQuestionDetail struct {
	AssessmentQuestionID uint                   `json:"assessment_question_id"`
	QuestionID           uint                   `json:"question_id"`
	Weight               int                    `json:"weight"`
	FixedOrder           *int                   `json:"fixed_order,omitempty"`
	Question             QuestionDetailResponse `json:"question"`
}
type UpdateAssessmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active closed"`
}

type UpdateUBDateRequest struct {
	ActiveFrom  *time.Time `json:"active_from,omitempty"`
	ActiveUntil *time.Time `json:"active_until,omitempty"`
}

type AddQuestionToAssessmentRequest struct {
	Questions []AssessmentQuestion `json:"questions" binding:"required,min=1"`
}

type UpdateQuestionWeightRequest struct {
	Weight int `json:"weight" binding:"required,min=1"`
}

type UpdateQuestionOrderRequest struct {
	FixedOrder *int `json:"fixed_order"` // null = random, angka = urutan tetap
}

type AssessmentQuestionResponse struct {
	AssessmentQuestionID uint                   `json:"assessment_question_id"`
	AssessmentID         uint                   `json:"assessment_id"`
	QuestionID           uint                   `json:"question_id"`
	Weight               int                    `json:"weight"`
	FixedOrder           *int                   `json:"fixed_order,omitempty"`
	Question             QuestionDetailResponse `json:"question"`
}

type AssessmentQuestionDetailResponse struct {
	AssessmentQuestionID uint                       `json:"assessment_question_id"`
	AssessmentID         uint                       `json:"assessment_id"`
	QuestionID           uint                       `json:"question_id"`
	Weight               int                        `json:"weight"`
	FixedOrder           *int                       `json:"fixed_order,omitempty"`
	Question             QuestionFullDetailResponse `json:"question"`
}

// QuestionFullDetailResponse - mirip dengan response di /bank/:id/questions
type QuestionFullDetailResponse struct {
	QuestionID   uint      `json:"question_id"`
	BankID       uint      `json:"bank_id"`
	QuestionType string    `json:"question_type"`
	QuestionText string    `json:"question_text"`
	Explanation  string    `json:"explanation"`
	CreatedBy    uint64    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Untuk tipe PG
	Options []QuestionOptionResponse `json:"options,omitempty"`

	// Untuk tipe SHORT_ANSWER
	ValidAnswers []string `json:"valid_answers,omitempty"`

	// Untuk tipe MATCHING
	MatchingPairs []MatchingPairResponse `json:"matching_pairs,omitempty"`
}

type QuestionOptionResponse struct {
	OptionID    uint   `json:"option_id"`
	OptionText  string `json:"option_text"`
	IsCorrect   bool   `json:"is_correct"`
	OptionOrder int    `json:"option_order"`
}

// ========== ENROLL STUDENT DTO ==========

type EnrollStudentRequest struct {
	StudentIDs []int `json:"student_ids" binding:"required,min=1"`
}

type EnrollStudentResponse struct {
	AssessmentID      uint   `json:"assessment_id"`
	TotalEnrolled     int    `json:"total_enrolled"`
	FailedEnrollments int    `json:"failed_enrollments"`
	Message           string `json:"message"`
}

// ========== GET ASSESSMENT STUDENTS DTO ==========

type AssessmentStudentResponse struct {
	AssessmentStudentID uint    `json:"assessment_student_id"`
	SiswaID             int     `json:"siswa_id"`
	SiswaNama           string  `json:"siswa_nama"`
	SiswaNIS            string  `json:"siswa_nis"`
	SiswaNISN           string  `json:"siswa_nisn"`
	Kelas               string  `json:"kelas"` // 🔥 String biasa, bukan pointer
	AttemptStatus       string  `json:"attempt_status"`
	Score               *string `json:"score"`        // Bisa null
	SubmittedAt         *string `json:"submitted_at"` // Bisa null
	IsActive            bool    `json:"is_active"`
	EnrolledAt          string  `json:"enrolled_at"`
}
