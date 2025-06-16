package controllers

import (
	"library-backend/config"
	"library-backend/models"

	"github.com/gofiber/fiber/v2"
)

// ✅ GET /api/kategori
func GetAllKategori(c *fiber.Ctx) error {
	var kategori []models.Category

	if err := config.DB.Find(&kategori).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data kategori",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil mengambil data kategori",
		"data":    kategori,
	})
}

// ✅ POST /api/kategori
func CreateKategori(c *fiber.Ctx) error {
	var input models.Category

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama kategori tidak boleh kosong",
		})
	}

	// Cek jika nama sudah ada
	var existing models.Category
	if err := config.DB.Where("name = ?", input.Name).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Kategori dengan nama ini sudah ada",
		})
	}

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan kategori",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Kategori berhasil ditambahkan",
		"data":    input,
	})
}
