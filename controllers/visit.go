package controllers

import (
	"library-backend/config"
	"library-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type VisitInput struct {
	UserID      uint   `json:"user_id"`
	VisitMethod string `json:"visit_method"`
}

func ProcessVisitInput(c *fiber.Ctx) error {
	var input VisitInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format salah"})
	}

	// Validasi method hanya boleh "qr" atau "manual"
	if input.VisitMethod != "qr" && input.VisitMethod != "manual" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Metode kunjungan tidak valid (hanya 'qr' atau 'manual')",
		})
	}

	admin := c.Locals("user").(*jwt.Token)
	claims := admin.Claims.(jwt.MapClaims)
	adminID := uint(claims["id"].(float64))

	timestamp := time.Now()
	var user models.User
	if err := config.DB.First(&user, adminID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	visit := models.Visit{
		UserID:    input.UserID,
		Method:    input.VisitMethod,
		CreatedBy: user.Nama,
		VisitTime: timestamp,
	}

	if err := config.DB.Create(&visit).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan kunjungan"})
	}

	return c.JSON(fiber.Map{"message": "Kunjungan berhasil dicatat"})
}

func VisitStatistic(c *fiber.Ctx) error {
	type Result struct {
		UserID uint   `json:"user_id"`
		Name   string `json:"name"`
		Count  int    `json:"total_kunjungan"`
	}

	var results []Result

	config.DB.Table("visits").
		Select("users.id as user_id, users.nama as name, COUNT(visits.id) as count").
		Joins("JOIN users ON users.id = visits.user_id").
		Group("users.id, users.nama").
		Order("count DESC").
		Scan(&results)

	return c.JSON(results)
}

func VisitByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")

	var visits []models.Visit
	err := config.DB.Preload("User").
		Where("user_id = ?", userID).
		Order("visit_time desc").
		Find(&visits).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil riwayat kunjungan"})
	}

	return c.JSON(visits)
}

func GetVisitHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var visits []models.Visit
	err := config.DB.Preload("User").
		Where("user_id = ?", userID).
		Order("visit_time desc").
		Find(&visits).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil riwayat kunjungan"})
	}

	return c.JSON(visits)
}

func GetVisitByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var visit models.Visit
	if err := config.DB.Preload("User").First(&visit, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Kunjungan tidak ditemukan",
		})
	}

	// Ambil user dari token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))
	role := claims["role"].(string)

	// Hanya pemilik kunjungan atau admin yang boleh akses
	if visit.UserID != userID && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Anda tidak memiliki izin untuk melihat kunjungan ini",
		})
	}

	return c.JSON(visit)
}
