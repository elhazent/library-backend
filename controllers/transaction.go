package controllers

import (
	"fmt"
	"library-backend/config"
	"library-backend/models"
	"library-backend/utils"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type CreateTransactionInput struct {
	UserID    uint `json:"user_id" validate:"required"`
	BookItems []struct {
		BookID uint `json:"book_id"`
		Qty    int  `json:"qty"`
	} `json:"book_items"`
}

func CreateTransaction(c *fiber.Ctx) error {
	admin := c.Locals("user").(*jwt.Token)
	claims := admin.Claims.(jwt.MapClaims)
	adminID := uint(claims["id"].(float64))

	var input CreateTransactionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input tidak valid"})
	}

	// Validasi: semua buku tersedia?
	for _, item := range input.BookItems {
		var book models.Book
		if err := config.DB.First(&book, item.BookID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Buku dengan ID %d tidak ditemukan", item.BookID)})
		}
		if book.Stock < item.Qty {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("Stok buku '%s' tidak mencukupi (tersisa %d, diminta %d)", book.Title, book.Stock, item.Qty),
			})
		}
	}

	var user models.User
	if err := config.DB.First(&user, adminID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	tx := models.Transaction{
		UserID:         input.UserID,
		CreatedBy:      user.Nama,
		UpdatedBy:      user.Nama, // <-- di sini ditambahkan
		TanggalPinjam:  time.Now(),
		TanggalKembali: time.Now().Add(7 * 24 * time.Hour),
		Status:         "dipinjam",
	}

	if err := config.DB.Create(&tx).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat transaksi"})
	}

	// Simpan detail dan kurangi stok
	for _, item := range input.BookItems {
		detail := models.TransactionDetail{
			TransactionID: tx.ID,
			BookID:        item.BookID,
			Qty:           item.Qty,
		}
		config.DB.Create(&detail)

		config.DB.Model(&models.Book{}).
			Where("id = ?", item.BookID).
			Update("stock", gorm.Expr("stock - ?", item.Qty))
	}

	var completeTx models.Transaction
	if err := config.DB.Preload("User").Preload("Details.Book").
		First(&completeTx, tx.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memuat detail transaksi"})
	}
	return c.JSON(fiber.Map{
		"message": "Transaksi berhasil dibuat",
		"data":    completeTx,
	})
}

func ReturnTransaction(c *fiber.Ctx) error {
	admin := c.Locals("user").(*jwt.Token)
	claims := admin.Claims.(jwt.MapClaims)
	adminID := uint(claims["id"].(float64))
	id := c.Params("id")
	var tx models.Transaction

	if err := config.DB.Preload("Details").First(&tx, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaksi tidak ditemukan"})
	}
	if tx.Status != "dipinjam" {
		return c.Status(400).JSON(fiber.Map{"error": "Transaksi bukan status dipinjam"})
	}

	var user models.User
	if err := config.DB.First(&user, adminID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	now := time.Now()

	// Hitung hari keterlambatan
	delayHours := now.Sub(tx.TanggalKembali).Hours() / 24
	keterlambatan := int(math.Max(0, delayHours))

	// Ambil harga denda dari master_data
	hargaPerHari := utils.GetHargaDendaPerHari()
	totalDenda := keterlambatan * hargaPerHari

	// Update transaksi
	tx.TanggalSelesai = &now
	tx.TotalDenda = totalDenda
	tx.UpdatedAt = now
	tx.UpdatedBy = user.Nama

	if totalDenda > 0 {
		tx.Status = "denda"
	} else {
		tx.Status = "dikembalikan"
	}

	if err := config.DB.Save(&tx).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update transaksi"})
	}

	// Kembalikan stok buku
	for _, d := range tx.Details {
		config.DB.Model(&models.Book{}).
			Where("id = ?", d.BookID).
			Update("stock", gorm.Expr("stock + ?", d.Qty))
	}

	return c.JSON(fiber.Map{
		"message":         "Pengembalian berhasil",
		"denda":           totalDenda,
		"hari_terlambat":  keterlambatan,
		"tanggal_kembali": now.Format("2006-01-02 15:04:05"),
	})
}

func LunasDenda(c *fiber.Ctx) error {
	admin := c.Locals("user").(*jwt.Token)
	claims := admin.Claims.(jwt.MapClaims)
	adminID := uint(claims["id"].(float64))
	id := c.Params("id")
	var tx models.Transaction

	if err := config.DB.First(&tx, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaksi tidak ditemukan"})
	}

	if tx.Status != "denda" {
		return c.Status(400).JSON(fiber.Map{"error": "Transaksi bukan status denda"})
	}

	var user models.User
	if err := config.DB.First(&user, adminID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	now := time.Now()
	tx.UpdatedAt = now
	tx.UpdatedBy = user.Nama

	tx.Status = "dikembalikan"
	if err := config.DB.Save(&tx).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memperbarui status"})
	}

	return c.JSON(fiber.Map{"message": "Denda berhasil dilunasi"})
}

func GetTransaksiUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))
	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")

	var transaksis []models.Transaction
	db := config.DB.Preload("User").Preload("Details.Book").
		Where("user_id = ?", userID)

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if from != "" && to != "" {
		layout := "2006-01-02"
		fromDate, _ := time.Parse(layout, from)
		toDate, _ := time.Parse(layout, to)
		db = db.Where("tanggal_pinjam BETWEEN ? AND ?", fromDate, toDate)
	}

	if err := db.Order("created_at desc").Find(&transaksis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil transaksi user"})
	}
	return c.JSON(transaksis)
}

func GetTransaksiByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")
	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")

	var transaksis []models.Transaction
	db := config.DB.Preload("User").Preload("Details.Book").
		Where("user_id = ?", userID)

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if from != "" && to != "" {
		layout := "2006-01-02"
		fromDate, _ := time.Parse(layout, from)
		toDate, _ := time.Parse(layout, to)
		db = db.Where("tanggal_pinjam BETWEEN ? AND ?", fromDate, toDate)
	}

	if err := db.Order("created_at desc").Find(&transaksis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil transaksi user"})
	}
	return c.JSON(transaksis)
}

func GetAllTransaksiAdmin(c *fiber.Ctx) error {
	status := c.Query("status") // e.g., ?status=dikembalikan
	from := c.Query("from")     // e.g., ?from=2025-06-01
	to := c.Query("to")         // e.g., ?to=2025-06-30

	var transaksis []models.Transaction
	db := config.DB.Preload("User").Preload("Details.Book")

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if from != "" && to != "" {
		layout := "2006-01-02"
		fromDate, err1 := time.Parse(layout, from)
		toDate, err2 := time.Parse(layout, to)
		if err1 == nil && err2 == nil {
			db = db.Where("tanggal_pinjam BETWEEN ? AND ?", fromDate, toDate)
		}
	}

	if err := db.Order("created_at desc").Find(&transaksis).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil transaksi"})
	}
	return c.JSON(transaksis)
}

func GetTransaksiByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var tx models.Transaction
	if err := config.DB.Preload("User").Preload("Details.Book").First(&tx, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaksi tidak ditemukan"})
	}

	// Ambil data user dari JWT
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	currentUserID := uint(claims["id"].(float64))
	currentUserRole := claims["role"].(string)

	// Hanya izinkan admin atau pemilik transaksi
	if tx.UserID != currentUserID && currentUserRole != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Anda tidak memiliki akses ke transaksi ini",
		})
	}

	return c.JSON(tx)
}
