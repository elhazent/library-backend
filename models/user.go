package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"` // 🔁 TAG PENTING!
	Role     string `json:"role"`
}

func (u *User) SetPassword(password string) error {
	fmt.Println("[DEBUG] Plain Password before hashing:", password) // ⬅️ log
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	fmt.Println("[DEBUG] Hashed Password:", u.Password) // ⬅️ log hasil hash
	return nil
}

func (u *User) CheckPassword(password string) bool {
	fmt.Println("[DEBUG] Comparing input password:", password)
	fmt.Println("[DEBUG] Against hash stored:", u.Password)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		fmt.Println("[DEBUG] Bcrypt error:", err)
	}
	return err == nil
}
