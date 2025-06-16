package routes

import (
	"library-backend/controllers"
	"library-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")
	admin := api.Group("/admin", middleware.Protected(), middleware.IsAdmin())

	//anggota
	api.Get("/users", controllers.GetUsers)

	auth := api.Group("/auth")
	auth.Put("/change-password", middleware.Protected(), controllers.UbahPassword)
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	book := api.Group("/book", middleware.Protected())
	book.Get("/", controllers.Getbook)
	book.Get("/:id", controllers.GetbookByID)

	transaction := api.Group("/transaction", middleware.Protected())
	transaction.Get("/", controllers.GetTransaksiUser)
	transaction.Get("/:id", controllers.GetTransaksiByID)

	visit := api.Group("/visits", middleware.Protected())
	visit.Post("/", controllers.GetVisitHistory)
	visit.Get("/:id", controllers.GetVisitByID)

	//admin
	visitAdmin := admin.Group("/visits")
	visitAdmin.Post("/input", controllers.ProcessVisitInput)
	visitAdmin.Get("/report", controllers.VisitStatistic)  // admin
	visitAdmin.Get("/user/:id", controllers.VisitByUserID) // admin
	visitAdmin.Get("/:id", controllers.GetVisitByID)

	transactionAdmin := admin.Group("/transaction")
	transactionAdmin.Post("/create", controllers.CreateTransaction)
	transactionAdmin.Put("/return/:id", controllers.ReturnTransaction)
	transactionAdmin.Put("/flag-lunas/:id", controllers.LunasDenda)
	transactionAdmin.Get("/:id", controllers.GetTransaksiByID)
	transactionAdmin.Get("/all", controllers.GetAllTransaksiAdmin)      // all user
	transactionAdmin.Get("/user/:id", controllers.GetTransaksiByUserID) // admin by ID

	bookAdmin := admin.Group("/book")
	bookAdmin.Post("/add", controllers.Createbook)
	bookAdmin.Put("/:id", controllers.Updatebook)
	bookAdmin.Delete("/:id", controllers.Deletebook)

	categoryAdmin := admin.Group("/category")
	categoryAdmin.Get("/", controllers.GetAllKategori)
	categoryAdmin.Post("/", controllers.CreateKategori)

}
