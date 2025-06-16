package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID             uint                `gorm:"primaryKey" json:"id"`
	UserID         uint                `gorm:"not null" json:"user_id"`
	User           User                `json:"user"`
	TanggalPinjam  time.Time           `json:"tanggal_pinjam"`
	TanggalKembali time.Time           `json:"tanggal_kembali"`
	TanggalSelesai *time.Time          `json:"tanggal_selesai"`
	Status         string              `json:"status"`
	TotalDenda     int                 `json:"total_denda"`
	Details        []TransactionDetail `json:"details"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"deleted_at"`
	CreatedBy      string              `gorm:"not null" json:"created_by"` // admin ID yang memproses
	UpdatedBy      string              `gorm:"not null" json:"updated_by"`
}

// TableName sets the insert table name for this struct type
func (Transaction) TableName() string {
	return "transactions"
}
