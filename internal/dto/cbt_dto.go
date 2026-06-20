package dto

import "time"

// ========== BANK DTO ==========
type CreateBankRequest struct {
	BankName    string `json:"bank_name" binding:"required"`
	KdMapel     uint64 `json:"kd_mapel" binding:"required"`
	KelasID     uint64 `json:"kelas_id" binding:"required"`
	Description string `json:"description"`
}

type UpdateBankRequest struct {
	BankName    *string `json:"bank_name,omitempty"`
	KdMapel     *uint64 `json:"kd_mapel,omitempty"`
	KelasID     *uint64 `json:"kelas_id,omitempty"`
	Description *string `json:"description,omitempty"`
}

type BankResponse struct {
	BankID         uint      `json:"bank_id"`
	BankName       string    `json:"bank_name"`
	KdMapel        uint64    `json:"kd_mapel"`
	MapelName      string    `json:"mapel_name"`
	KelasID        uint64    `json:"kelas_id"`
	KelasName      string    `json:"kelas_name"`
	Description    *string   `json:"description,omitempty"`
	TotalQuestions int64     `json:"total_questions"`
	CreatedBy      uint64    `json:"created_by"`
	CreatedByName  string    `json:"created_by_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BankDetailResponse - untuk halaman detail bank (dengan statistik tipe soal)
type BankDetailResponse struct {
	BankID         uint             `json:"bank_id"`
	BankName       string           `json:"bank_name"`
	KdMapel        uint64           `json:"kd_mapel"`
	MapelName      string           `json:"mapel_name"`
	KelasID        uint64           `json:"kelas_id"`
	KelasName      string           `json:"kelas_name"`
	Description    *string          `json:"description,omitempty"`
	TotalQuestions int64            `json:"total_questions"`
	QuestionTypes  map[string]int64 `json:"question_types"` // Statistik tipe soal
	CreatedBy      uint64           `json:"created_by"`
	CreatedByName  string           `json:"created_by_name"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// dto/bank.go
type BankByKelasResponse struct {
	BankID         uint      `json:"bank_id"`
	BankName       string    `json:"bank_name"`
	KdMapel        uint64    `json:"kd_mapel"`
	MapelName      string    `json:"mapel_name"`
	Description    *string   `json:"description,omitempty"`
	TotalQuestions int64     `json:"total_questions"`
	CreatedAt      time.Time `json:"created_at"`
}

// ========== QUESTION DTO ==========
type CreateQuestionRequest struct {
	BankID        *uint                   `json:"bank_id,omitempty"`
	QuestionType  string                  `json:"question_type" binding:"required,oneof=MCQ MULTI_MCQ TRUE_FALSE SHORT_ANSWER ESSAY MATCHING"`
	QuestionText  string                  `json:"question_text" binding:"required"`
	Explanation   *string                 `json:"explanation,omitempty"`
	Options       []CreateOptionRequest   `json:"options,omitempty"`
	ValidAnswers  []string                `json:"valid_answers,omitempty"`
	MatchingPairs []CreateMatchingRequest `json:"matching_pairs,omitempty"`
}

type CreateOptionRequest struct {
	Label     string `json:"label" binding:"required,max=10"`
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
	SortOrder int    `json:"sort_order"`
}

type CreateMatchingRequest struct {
	LeftText  string `json:"left_text" binding:"required"`
	RightText string `json:"right_text" binding:"required"`
	PairOrder int    `json:"pair_order"`
}

type UpdateQuestionRequest struct {
	BankID        *uint                   `json:"bank_id,omitempty"`
	QuestionType  *string                 `json:"question_type,omitempty" binding:"omitempty,oneof=MCQ MULTI_MCQ TRUE_FALSE SHORT_ANSWER ESSAY MATCHING"`
	QuestionText  string                  `json:"question_text,omitempty"`
	Explanation   *string                 `json:"explanation,omitempty"`
	Options       []CreateOptionRequest   `json:"options,omitempty"`
	ValidAnswers  []string                `json:"valid_answers,omitempty"`
	MatchingPairs []CreateMatchingRequest `json:"matching_pairs,omitempty"`
}
type QuestionResponse struct {
	QuestionID   uint      `json:"question_id"`
	BankID       *uint     `json:"bank_id,omitempty"`
	QuestionType string    `json:"question_type"`
	QuestionText string    `json:"question_text"`
	Explanation  *string   `json:"explanation,omitempty"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// QuestionDetailResponse - untuk list question di detail bank
type QuestionDetailResponse struct {
	QuestionID    uint                   `json:"question_id"`
	BankID        *uint                  `json:"bank_id,omitempty"`
	QuestionType  string                 `json:"question_type"`
	QuestionText  string                 `json:"question_text"`
	Explanation   *string                `json:"explanation,omitempty"`
	CreatedBy     int                    `json:"created_by"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Options       []OptionResponse       `json:"options,omitempty"`
	ValidAnswers  []string               `json:"valid_answers,omitempty"`
	MatchingPairs []MatchingPairResponse `json:"matching_pairs,omitempty"`
}

type OptionResponse struct {
	OptionID    uint   `json:"option_id"`
	OptionLabel string `json:"option_label"`
	OptionText  string `json:"option_text"`
	IsCorrect   bool   `json:"is_correct"`
	SortOrder   int    `json:"sort_order"`
}

type MatchingPairResponse struct {
	MatchingID uint   `json:"matching_id"`
	LeftText   string `json:"left_text"`
	RightText  string `json:"right_text"`
	PairOrder  int    `json:"pair_order"`
}
