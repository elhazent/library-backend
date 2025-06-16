package main

import (
	"library-backend/config"
	"library-backend/models"
	"library-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Book{}, &models.Transaction{}, &models.MasterData{}, &models.TransactionDetail{}, &models.Visit{})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Perpustakaan API berjalan!")
	})
	routes.SetUpRoutes(app)

	app.Listen(":8080")

}
