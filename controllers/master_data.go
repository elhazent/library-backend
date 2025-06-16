package controllers

import (
	"library-backend/config"
	"library-backend/models"

	"github.com/gofiber/fiber/v2"
)

// POST /masterdata
type MasterDataInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func CreateOrUpdateMasterData(c *fiber.Ctx) error {
	var input MasterDataInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format salah"})
	}

	var md models.MasterData
	err := config.DB.Where("key = ?", input.Key).First(&md).Error
	if err == nil {
		md.Value = input.Value
		config.DB.Save(&md)
	} else {
		md = models.MasterData{KeyName: input.Key, Value: input.Value}
		config.DB.Create(&md)
	}

	return c.JSON(fiber.Map{"message": "Master data berhasil disimpan", "data": md})
}
