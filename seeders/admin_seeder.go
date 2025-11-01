package seeders

import (
	"go_evermos_rakamin_irsan/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@example.com").Or("notelp = ?", "08987654321").Count(&count)

	if count > 0 {
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Nama: "Admin",
		KataSandi: string(hashed),
		NoTelp: "08987654321",
		TanggalLahir: time.Date(1990,1,1,0,0,0,0,time.UTC),
		Email: "admin@example.com",
		Pekerjaan: "Admin E-Commerce",
		Tentang: "Akun Admin",
		JenisKelamin: "L",
		IsAdmin: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Println("Admin Created")
	return nil
}