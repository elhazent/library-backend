package controllers

import (
	"errors"
	"library-backend/config"
	"library-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Getbook(c *fiber.Ctx) error {
	var books []models.Book

	// Ambil query parameter page dan pageSize, default ke 1 dan 10 jika tidak ada
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var total int64
	config.DB.Model(&models.Book{}).Count(&total)

	offset := (page - 1) * pageSize

	if err := config.DB.Preload("Category").Limit(pageSize).Offset(offset).Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data buku",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil mengambil data buku",
		"data": fiber.Map{
			"data":       books,
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func GetbookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	if err := config.DB.Preload("Category").First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Buku tidak ditemukan",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil detail buku",
		})
	}

	return c.JSON(book)
}

func Createbook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data salah",
			"data":    nil,
		})
	}

	if err := config.DB.Create(&book).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat buku",
			"data":    nil,
		})
	}

	// Ambil data buku beserta category setelah berhasil dibuat
	if err := config.DB.Preload("Category").First(&book, book.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data buku beserta kategori",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Buku berhasil dibuat",
		"data":    book,
	})
}

func Updatebook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	// Cari buku berdasarkan ID
	if err := config.DB.First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Buku tidak ditemukan"})
	}

	// Parse body ke struct book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data salah"})
	}

	// Simpan perubahan
	if err := config.DB.Save(&book).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memperbarui buku"})
	}

	// Ambil data buku beserta category setelah update
	if err := config.DB.Preload("Category").First(&book, book.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data buku beserta kategori"})
	}

	return c.JSON(book)
}

func Deletebook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	// Cek apakah book dengan ID tersebut ada (termasuk yang sudah soft delete)
	if err := config.DB.Unscoped().First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "book tidak ditemukan"})
	}

	// Soft delete
	if err := config.DB.Delete(&book).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus book"})
	}

	return c.JSON(fiber.Map{"message": "book berhasil dihapus (soft delete)"})
}
