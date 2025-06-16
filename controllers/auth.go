package controllers

import (
	"os"
	"time"

	"library-backend/config"
	"library-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var input models.User

	// Parsing request body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format permintaan tidak valid",
			"status":  false,
		})
	}

	// Cek jika email sudah terdaftar
	var existing models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email sudah terdaftar",
			"status":  false,
		})
	}

	// Hash password
	if err := input.SetPassword(input.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengenkripsi password",
			"status":  false,
		})
	}

	input.Role = "anggota" // default role

	// Simpan ke database
	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mendaftarkan user",
			"status":  false,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"status":  true,
		"data":    input,
	})
}

func Login(c *fiber.Ctx) error {

	var user models.User

	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	config.DB.Where("email = ?", input.Email).First(&user)

	if user.ID == 0 || !user.CheckPassword(input.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	return c.JSON(fiber.Map{
		"message":    "Berhasil login",
		"success":    true,
		"token":      t,
		"token_type": "bearer",
	})
}
