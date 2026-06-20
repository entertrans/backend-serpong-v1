// controller/cbt/bank_helper.go
package cbt

import (
	"fmt"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

// ============================================
// HELPER FUNCTIONS
// ============================================

func getUserName(user *model.User) string {
	if user != nil {
		return user.Name
	}
	return ""
}

func getMapelName(mapel *model.Mapel) string {
	if mapel != nil {
		return mapel.NmMapel
	}
	return ""
}

func getKelasName(kelas *model.Kelas) string {
	if kelas != nil {
		return kelas.KelasNama
	}
	return ""
}

func createOptions(tx *gorm.DB, questionID uint, options []dto.CreateOptionRequest) error {
	if len(options) == 0 {
		return nil
	}

	opts := make([]model.ToQuestionOption, len(options))
	for i, opt := range options {
		opts[i] = model.ToQuestionOption{
			QuestionID:  questionID,
			OptionLabel: opt.Label,
			OptionText:  opt.Text,
			IsCorrect:   opt.IsCorrect,
			SortOrder:   opt.SortOrder,
		}
	}
	return tx.Create(&opts).Error
}

func createShortAnswers(tx *gorm.DB, questionID uint, validAnswers []string) error {
	if len(validAnswers) == 0 {
		return nil
	}

	answers := make([]model.ToQuestionShortAnswer, len(validAnswers))
	for i, ans := range validAnswers {
		answers[i] = model.ToQuestionShortAnswer{
			QuestionID:  questionID,
			ValidAnswer: ans,
		}
	}
	return tx.Create(&answers).Error
}

func createMatchingPairs(tx *gorm.DB, questionID uint, pairs []dto.CreateMatchingRequest) error {
	if len(pairs) == 0 {
		return nil
	}

	matchingPairs := make([]model.ToQuestionMatching, len(pairs))
	for i, pair := range pairs {
		matchingPairs[i] = model.ToQuestionMatching{
			QuestionID: questionID,
			LeftText:   pair.LeftText,
			RightText:  pair.RightText,
			PairOrder:  pair.PairOrder,
		}
	}
	return tx.Create(&matchingPairs).Error
}

// Convert userID dari interface ke int
func convertToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case uint:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("invalid user id type: %T", value)
	}
}

// Build question detail response
func buildQuestionDetailResponse(db *gorm.DB, question *model.ToQuestion) dto.QuestionDetailResponse {
	response := dto.QuestionDetailResponse{
		QuestionID:   question.QuestionID,
		BankID:       question.BankID,
		QuestionType: question.QuestionType,
		QuestionText: question.QuestionText,
		Explanation:  question.Explanation,
		CreatedBy:    question.CreatedBy,
		CreatedAt:    question.CreatedAt,
		UpdatedAt:    question.UpdatedAt,
	}

	switch question.QuestionType {
	case "MCQ", "MULTI_MCQ", "TRUE_FALSE":
		var options []model.ToQuestionOption
		db.Where("question_id = ?", question.QuestionID).
			Order("sort_order ASC").
			Find(&options)

		opts := make([]dto.OptionResponse, len(options))
		for i, opt := range options {
			opts[i] = dto.OptionResponse{
				OptionID:    opt.OptionID,
				OptionLabel: opt.OptionLabel,
				OptionText:  opt.OptionText,
				IsCorrect:   opt.IsCorrect,
				SortOrder:   opt.SortOrder,
			}
		}
		response.Options = opts

	case "SHORT_ANSWER":
		var answers []model.ToQuestionShortAnswer
		db.Where("question_id = ?", question.QuestionID).Find(&answers)

		validAnswers := make([]string, len(answers))
		for i, ans := range answers {
			validAnswers[i] = ans.ValidAnswer
		}
		response.ValidAnswers = validAnswers

	case "MATCHING":
		var pairs []model.ToQuestionMatching
		db.Where("question_id = ?", question.QuestionID).
			Order("pair_order ASC").
			Find(&pairs)

		matchingPairs := make([]dto.MatchingPairResponse, len(pairs))
		for i, pair := range pairs {
			matchingPairs[i] = dto.MatchingPairResponse{
				MatchingID: pair.MatchingID,
				LeftText:   pair.LeftText,
				RightText:  pair.RightText,
				PairOrder:  pair.PairOrder,
			}
		}
		response.MatchingPairs = matchingPairs
	}

	return response
}
