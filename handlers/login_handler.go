package handlers

import (
	"go_evermos_rakamin_irsan/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginHandler(db *gorm.DB, jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.LoginRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid request body",
			})
		}

		if req.NoTelp == "" || req.KataSandi == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Phone Number and Password are required",
			})
		}

		var user models.User
		if err := db.Where("notelp = ?", req.NoTelp).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"message": "Phone number not found",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(req.KataSandi)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid password",
			})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"is_admin": user.IsAdmin,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "failed to generate token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"token": tokenString,
		})
	}
}