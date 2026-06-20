package model

import "time"

type ToQuestionBank struct {
	BankID      uint       `gorm:"primaryKey;autoIncrement" json:"bank_id"`
	BankName    string     `gorm:"type:varchar(255);not null" json:"bank_name"`
	KdMapel     uint64     `gorm:"not null" json:"kd_mapel"`
	KelasID     uint64     `gorm:"not null" json:"kelas_id"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedBy   uint64     `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	User  *User  `gorm:"foreignKey:CreatedBy;references:ID" json:"user,omitempty"`
	Mapel *Mapel `gorm:"foreignKey:KdMapel;references:KdMapel" json:"mapel,omitempty"`
	Kelas *Kelas `gorm:"foreignKey:KelasID;references:KelasId" json:"kelas,omitempty"`
}

func (ToQuestionBank) TableName() string { return "to_question_bank" }

// ToQuestion - tabel to_question
type ToQuestion struct {
	QuestionID   uint       `gorm:"primaryKey;autoIncrement" json:"question_id"`
	BankID       *uint      `gorm:"index" json:"bank_id,omitempty"`
	QuestionType string     `gorm:"type:enum('MCQ','MULTI_MCQ','TRUE_FALSE','SHORT_ANSWER','ESSAY','MATCHING');not null" json:"question_type"`
	QuestionText string     `gorm:"type:text;not null" json:"question_text"`
	Explanation  *string    `gorm:"type:text" json:"explanation,omitempty"`
	CreatedBy    int        `gorm:"not null" json:"created_by"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	Bank *ToQuestionBank `gorm:"foreignKey:BankID;references:BankID" json:"bank,omitempty"`
}

func (ToQuestion) TableName() string { return "to_question" }

// ToQuestionOption - tabel option untuk MCQ/MULTI_MCQ/TRUE_FALSE
type ToQuestionOption struct {
	OptionID       uint    `gorm:"primaryKey;autoIncrement" json:"option_id"`
	QuestionID     uint    `gorm:"not null;index" json:"question_id"`
	OptionLabel    string  `gorm:"type:varchar(10);not null" json:"option_label"`
	OptionText     string  `gorm:"type:text;not null" json:"option_text"`
	OptionImageURL *string `gorm:"type:varchar(500)" json:"option_image_url,omitempty"`
	IsCorrect      bool    `gorm:"default:false" json:"is_correct"`
	SortOrder      int     `gorm:"default:0" json:"sort_order"`
}

func (ToQuestionOption) TableName() string { return "to_question_option" }

// ToQuestionShortAnswer - tabel valid answer untuk SHORT_ANSWER
type ToQuestionShortAnswer struct {
	ShortAnswerID uint   `gorm:"primaryKey;autoIncrement" json:"short_answer_id"`
	QuestionID    uint   `gorm:"not null;index" json:"question_id"`
	ValidAnswer   string `gorm:"type:text;not null" json:"valid_answer"`
}

func (ToQuestionShortAnswer) TableName() string { return "to_question_short_answer" }

// ToQuestionMatching - tabel matching pairs untuk MATCHING
type ToQuestionMatching struct {
	MatchingID uint   `gorm:"primaryKey;autoIncrement" json:"matching_id"`
	QuestionID uint   `gorm:"not null;index" json:"question_id"`
	LeftText   string `gorm:"type:varchar(255);not null" json:"left_text"`
	RightText  string `gorm:"type:varchar(255);not null" json:"right_text"`
	PairOrder  int    `gorm:"default:0" json:"pair_order"`
}

func (ToQuestionMatching) TableName() string { return "to_question_matching" }
