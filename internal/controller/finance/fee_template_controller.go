package finance

import (
	"errors"
	"fmt"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

type FeeTemplateController interface {
	GetAllFeeTemplates(isActive *int8) ([]dto.FeeTemplateResponse, error)
	GetFeeTemplateByID(id uint64) (*dto.FeeTemplateResponse, error)
	CreateFeeTemplate(req *dto.FeeTemplateRequest) (*dto.FeeTemplateResponse, error)
	UpdateFeeTemplate(id uint64, req *dto.FeeTemplateRequest) (*dto.FeeTemplateResponse, error)
	DeleteFeeTemplate(id uint64) error
}

type feeTemplateController struct {
	db *gorm.DB
}

func NewFeeTemplateController(db *gorm.DB) FeeTemplateController {
	return &feeTemplateController{db: db}
}

func (c *feeTemplateController) GetAllFeeTemplates(isActive *int8) ([]dto.FeeTemplateResponse, error) {
	var templates []model.FinanceFeeTemplate

	query := c.db.Model(&model.FinanceFeeTemplate{})
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("fee_template_id DESC").Find(&templates).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.FeeTemplateResponse, len(templates))
	for i, t := range templates {
		result[i] = dto.FeeTemplateResponse{
			FeeTemplateID: t.FeeTemplateID,
			FeeCode:       t.FeeCode,
			FeeName:       t.FeeName,
			DefaultAmount: t.DefaultAmount,
			Description:   t.Description,
			IsActive:      t.IsActive,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		}
	}

	return result, nil
}

func (c *feeTemplateController) GetFeeTemplateByID(id uint64) (*dto.FeeTemplateResponse, error) {
	var template model.FinanceFeeTemplate
	err := c.db.First(&template, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("fee template not found")
		}
		return nil, err
	}

	return &dto.FeeTemplateResponse{
		FeeTemplateID: template.FeeTemplateID,
		FeeCode:       template.FeeCode,
		FeeName:       template.FeeName,
		DefaultAmount: template.DefaultAmount,
		Description:   template.Description,
		IsActive:      template.IsActive,
		CreatedAt:     template.CreatedAt,
		UpdatedAt:     template.UpdatedAt,
	}, nil
}

func (c *feeTemplateController) CreateFeeTemplate(req *dto.FeeTemplateRequest) (*dto.FeeTemplateResponse, error) {
	template := model.FinanceFeeTemplate{
		FeeCode:       req.FeeCode,
		FeeName:       req.FeeName,
		DefaultAmount: req.DefaultAmount,
		Description:   req.Description,
	}

	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	} else {
		template.IsActive = 1
	}

	err := c.db.Create(&template).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create fee template: %w", err)
	}

	return &dto.FeeTemplateResponse{
		FeeTemplateID: template.FeeTemplateID,
		FeeCode:       template.FeeCode,
		FeeName:       template.FeeName,
		DefaultAmount: template.DefaultAmount,
		Description:   template.Description,
		IsActive:      template.IsActive,
		CreatedAt:     template.CreatedAt,
		UpdatedAt:     template.UpdatedAt,
	}, nil
}

func (c *feeTemplateController) UpdateFeeTemplate(id uint64, req *dto.FeeTemplateRequest) (*dto.FeeTemplateResponse, error) {
	var template model.FinanceFeeTemplate
	err := c.db.First(&template, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("fee template not found")
		}
		return nil, err
	}

	// Update fields
	if req.FeeCode != nil {
		template.FeeCode = req.FeeCode
	}
	if req.FeeName != "" {
		template.FeeName = req.FeeName
	}
	template.DefaultAmount = req.DefaultAmount
	if req.Description != nil {
		template.Description = req.Description
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	err = c.db.Save(&template).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update fee template: %w", err)
	}

	return &dto.FeeTemplateResponse{
		FeeTemplateID: template.FeeTemplateID,
		FeeCode:       template.FeeCode,
		FeeName:       template.FeeName,
		DefaultAmount: template.DefaultAmount,
		Description:   template.Description,
		IsActive:      template.IsActive,
		CreatedAt:     template.CreatedAt,
		UpdatedAt:     template.UpdatedAt,
	}, nil
}

func (c *feeTemplateController) DeleteFeeTemplate(id uint64) error {
	result := c.db.Delete(&model.FinanceFeeTemplate{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("fee template not found")
	}
	return nil
}
