package jobs

import (
	"fmt"

	"os"
	"time"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

// func ReminderJatuhTempo() {
// 	today := time.Now()
// 	deadline := today.Add(48 * time.Hour).Format("2006-01-02")

// 	var peminjaman []models.Borrow
// 	config.DB.Preload("User").Where("tanggal_kembali BETWEEN ? AND ? AND status = ?", today.Format("2006-01-02"), deadline, "dipinjam").Find(&peminjaman)

// 	for _, p := range peminjaman {
// 		fmt.Println("[DEBUG] Email Reminder:", p.User.Email) // ⬅️ log hasil hash
// 		sendEmail(p.User.Email, p.User.Nama, p.ReturnDate)
// 	}
// }

func sendEmail(to, name string, kembali time.Time) {
	mj := mailjet.NewMailjetClient(os.Getenv("MJ_API_KEY"), os.Getenv("MJ_SECRET_KEY"))

	messageInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "noreply.manassemalo@istekwiduri.ac.id",
				Name:  "Perpustakaan ISTEK Widuri",
			},
			To: &mailjet.RecipientsV31{
				{
					Email: to,
					Name:  name,
				},
			},
			Subject:  "📚 Pengingat Pengembalian Buku",
			TextPart: fmt.Sprintf("Halo %s, buku yang Anda pinjam akan jatuh tempo pada %s. Harap dikembalikan tepat waktu.", name, kembali.Format("02 Jan 2006")),
			HTMLPart: fmt.Sprintf("<h3>Halo %s</h3><p>Buku yang Anda pinjam akan jatuh tempo pada <strong>%s</strong>.<br>Harap dikembalikan tepat waktu agar tidak kena denda.</p>", name, kembali.Format("02 Jan 2006")),
		},
	}

	messages := mailjet.MessagesV31{Info: messageInfo}

	res, err := mj.SendMailV31(&messages)

	if err != nil {
		fmt.Println("❌ Gagal kirim email ke", to, ":", err)
	} else {
		fmt.Println("✅ Email dikirim ke", to, ":", res)
	}
}
