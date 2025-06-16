package utils

import (
	"library-backend/config"
	"library-backend/models"
	"strconv"
)

func GetHargaDendaPerHari() int {
	var md models.MasterData
	if err := config.DB.Where("key_name = ?", "harga_denda_per_hari").First(&md).Error; err != nil {
		return 0 // atau fallback default
	}

	val, err := strconv.Atoi(md.Value)
	if err != nil {
		return 0
	}
	return val
}
