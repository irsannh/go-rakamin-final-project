package handlers

import (
	"errors"
	"go_evermos_rakamin_irsan/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetMyProfileHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"message": "Unauthorized",
			})
		}

		var userID uint
		switch v := userIDInt.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case uint:
			userID = v
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid User ID",
			})
		}

		var user models.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": user,
		})
	}
}

func UpdateMyProfileHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"message": "Unauthorized",
			})
		}

		var userID uint
		switch v := userIDInt.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case uint:
			userID = v
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid User ID",
			})
		}

		var user models.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		var request models.RegisterAndUserRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid request body",
			})
		}

		var exists models.User
		errCheck := db.Where("email = ?", request.Email).Or("notelp = ?", request.NoTelp).Take(&exists)

		if errCheck.Error == nil && errCheck.RowsAffected > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Email or phone number already registered by another user",
			})
		}
		if errCheck.Error != nil && !errors.Is(errCheck.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Database error",
			})
		}

		if request.Nama != "" {
			user.Nama = request.Nama
		}

		if request.KataSandi != "" {
			hashed, err := bcrypt.GenerateFromPassword([]byte(request.KataSandi), bcrypt.DefaultCost)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": "error",
					"message": "Failed to hash password",
				})
			}

			user.KataSandi = string(hashed)
		}

		if request.NoTelp != "" {
			user.NoTelp = request.NoTelp
		}

		if request.TanggalLahir != "" {
			tglLahir := request.TanggalLahir
			pola := "02/01/2006"
			dateOfBirth, err := time.Parse(pola, tglLahir)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": "error",
					"message": "failed to parse date of birth",
				})
			}

			user.TanggalLahir = dateOfBirth
		}

		if request.Pekerjaan != "" {
			user.Pekerjaan = request.Pekerjaan
		}

		if request.Tentang != "" {
			user.Tentang = request.Tentang
		}

		if request.Email != "" {
			user.Email = request.Email
		}

		if request.IDProvinsi != "" {
			user.IDProvinsi = request.IDProvinsi
		}

		if request.IDKota != "" {
			user.IDKota = request.IDKota
		}

		if err := db.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		response := models.UserResponse{
			Nama: user.Nama,
			NoTelp: user.NoTelp,
			Email: user.NoTelp,
			TanggalLahir: user.TanggalLahir.String(),
			Pekerjaan: user.Pekerjaan,
			Tentang: user.Tentang,
			IDProvinsi: user.IDProvinsi,
			IDKota: user.IDKota,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": response,
		})
	}
}

