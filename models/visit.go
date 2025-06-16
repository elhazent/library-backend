package models

import (
	"time"

	"gorm.io/gorm"
)

type Visit struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Method    string         `gorm:"type:enum('manual','qr');not null" json:"method"`
	CreatedBy string         `gorm:"default:null" json:"created_by"` // admin ID jika manual
	VisitTime time.Time      `gorm:"not null" json:"visit_time"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
