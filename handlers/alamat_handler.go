package handlers

import (
	"fmt"
	"go_evermos_rakamin_irsan/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PostAlamatHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		var input models.AlamatRequest

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid Input",
			})
		}

		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
				"message": "Unauthorized",
			})
		}
		userID, _ := strconv.Atoi(fmt.Sprintf("%v", userIDInt))

		alamat := models.Alamat{
			UserID: uint(userID),
			JudulAlamat: input.JudulAlamat,
			NamaPenerima: input.NamaPenerima,
			NoTelp: input.NoTelp,
			DetailAlamat: input.DetailAlamat,
		}

		if err := db.Create(&alamat).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to create alamat",
			})
		}
		
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "success",
			"category": input,
		})
	}
}


func GetAllMyAlamatHandler(db *gorm.DB) fiber.Handler {
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

		var alamat []models.Alamat
		if err := db.Where("id_user", userID).Find(&alamat).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Cannot get alamat",
			})
		}

		var response []models.AlamatResponse
		for _, p := range alamat {
			res := models.AlamatResponse{
				JudulAlamat: p.JudulAlamat,
				NamaPenerima: p.NamaPenerima,
				NoTelp: p.NoTelp,
				DetailAlamat: p.DetailAlamat,
			}
			response = append(response, res)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "success",
				"alamat": response,
			})
	}
}

func GetMyAlamatByIdHandler(db *gorm.DB) fiber.Handler {
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

		idAlamat := c.Params("id")

		var alamat models.Alamat
		if err := db.Where("id_user = ? AND id = ?", userID,idAlamat).First(&alamat).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Cannot get alamat or not belongs to you",
			})
		}

		response := models.AlamatResponse{
				JudulAlamat: alamat.JudulAlamat,
				NamaPenerima: alamat.NamaPenerima,
				NoTelp: alamat.NoTelp,
				DetailAlamat: alamat.DetailAlamat,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "success",
				"alamat": response,
			})
	}
}


func UpdateMyAlamatHandler(db *gorm.DB) fiber.Handler {
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

		idAlamat := c.Params("id")

		var alamat models.Alamat
		if err := db.First(&alamat, idAlamat).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Cannot update alamat" ,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		if alamat.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"message": "You don't have any access to update this alamat",
			})
		}

		var request models.AlamatRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid request body",
			})
		}

		if request.NamaPenerima != "" {
			alamat.NamaPenerima = request.NamaPenerima
		}

		if request.JudulAlamat != "" {
			alamat.JudulAlamat = request.JudulAlamat
		}

		if request.DetailAlamat != "" {
			alamat.DetailAlamat = request.DetailAlamat
		}

		if request.NoTelp != "" {
			alamat.NoTelp = request.NoTelp
		}

		if err := db.Save(&alamat).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": alamat,
		})
	}
}

func DeleteMyAlamatHandler(db *gorm.DB) fiber.Handler {
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

		idAlamat := c.Params("id")

		var alamat models.Alamat
		if err := db.First(&alamat, idAlamat).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Cannot delete alamat" ,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		if alamat.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"message": "You don't have any access to update this alamat",
			})
		}

		if err := db.Delete(&alamat).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to delete alamat",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"message": "Alamat deleted successfully",
		})
	}
}


