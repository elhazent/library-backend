package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `json:"title"`
	Author     string         `json:"author"`
	Publisher  string         `json:"publisher"`
	Year       int            `json:"year"`
	CategoryID uint           `json:"category_id"`
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Stock      int            `json:"stock"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
