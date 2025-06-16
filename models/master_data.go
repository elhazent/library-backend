package models

type MasterData struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	KeyName string `gorm:"unique;not null" json:"key_name"` // Contoh: "harga_denda_per_hari"
	Value   string `gorm:"not null" json:"value"`           // Simpan dalam bentuk string, nanti dikonversi
}

// TableName sets the insert table name for this struct type
func (MasterData) TableName() string {
	return "master_data"
}
