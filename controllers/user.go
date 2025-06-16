package controllers

import (
	"library-backend/config"
	"library-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	config.DB.Find(&users)
	return c.JSON(users)
}

type PasswordUpdateInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UbahPassword(c *fiber.Ctx) error {
	var input PasswordUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format permintaan salah"})
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	if !user.CheckPassword(input.OldPassword) {
		return c.Status(401).JSON(fiber.Map{"error": "Password lama salah"})
	}

	err := user.SetPassword(input.NewPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengubah password"})
	}

	config.DB.Save(&user)

	return c.JSON(fiber.Map{"message": "Password berhasil diubah"})
}
