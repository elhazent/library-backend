package models

type TransactionDetail struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	TransactionID uint `json:"transaction_id"`
	// Transaction   Transaction `json:"transaction"`
	BookID uint `json:"book_id"`
	Book   Book `json:"book"`
	Qty    int  `json:"qty"`
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}
