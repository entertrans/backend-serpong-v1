// internal/repository/guru_repository.go
package repository

import (
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

type GuruRepository interface {
	FindByUserID(userID uint) (*model.Guru, error)
	FindByNIP(nip string) (*model.Guru, error)
	Create(guru *model.Guru) error
	Update(guru *model.Guru) error
}

type guruRepository struct {
	db *gorm.DB
}

func NewGuruRepository(db *gorm.DB) GuruRepository {
	return &guruRepository{db: db}
}

func (r *guruRepository) FindByUserID(userID uint) (*model.Guru, error) {
	var guru model.Guru
	err := r.db.Where("user_id = ?", userID).First(&guru).Error
	if err != nil {
		return nil, err
	}
	return &guru, nil
}

func (r *guruRepository) FindByNIP(nip string) (*model.Guru, error) {
	var guru model.Guru
	err := r.db.Where("guru_nip = ?", nip).First(&guru).Error
	if err != nil {
		return nil, err
	}
	return &guru, nil
}

func (r *guruRepository) Create(guru *model.Guru) error {
	return r.db.Create(guru).Error
}

func (r *guruRepository) Update(guru *model.Guru) error {
	return r.db.Save(guru).Error
}
