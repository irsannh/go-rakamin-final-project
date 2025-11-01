package handlers

import (
	"errors"
	"fmt"
	"go_evermos_rakamin_irsan/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.RegisterAndUserRequest
		var exists models.User

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid Input",
			})
		}

		if input.Nama == "" || input.Email == "" || input.KataSandi == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "name, email, and password are required",
			})
		}

		errCheck := db.Where("email = ?", input.Email).Or("notelp = ?", input.NoTelp).First(&exists).Error

		if errCheck == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Email or phone number already registered",
			})
		} else if !errors.Is(errCheck, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Database error",
			})
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(input.KataSandi), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to hash password",
			})
		}

		tglLahir := input.TanggalLahir
		pola := "02/01/2006"
		dateOfBirth, err := time.Parse(pola, tglLahir)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "failed to parse date of birth",
			})
		}

		user := models.User{
			Nama: input.Nama,
			KataSandi: string(hashed),
			NoTelp: input.NoTelp,
			TanggalLahir: dateOfBirth,
			JenisKelamin: input.JenisKelamin,
			Pekerjaan: input.Pekerjaan,
			Tentang: input.Tentang,
			Email: input.Email,
			IDProvinsi: input.IDProvinsi,
			IDKota: input.IDKota,
		}
		
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to create user",
			})
		}

		toko := models.Toko{
			UserID: user.ID,
			NamaToko: fmt.Sprintf("Toko-00%d", user.ID),
			URLFoto: "",
		}

		if err := db.Create(&toko).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to created store",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
			"user id": user.ID,
		})
	}
}