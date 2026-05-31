// internal/model/user.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;type:varchar(255)" json:"email"`
	Password  string         `gorm:"not null;type:varchar(255)" json:"-"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Role      string         `gorm:"type:varchar(50);default:user" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
