package handlers

import (
	"go_evermos_rakamin_irsan/models"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllTokos(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		nama := c.Query("nama", "")
		pageStr := c.Query("page")
		limitStr := c.Query("limit")

		query := db.Model(&models.Toko{})

		if nama != "" {
			query = query.Where("nama_toko LIKE ?", "%" + nama + "%")
		}

		var tokos []models.Toko
		var total int64
		query.Count(&total)

		if pageStr != "" && limitStr != "" {
			page, _ := strconv.Atoi(pageStr)
			limit, _ := strconv.Atoi(limitStr)
			if page < 1 {
				page = 1
			}
			offset := (page - 1) * limit
			query = query.Limit(limit).Offset(offset)
		}

		if err := query.Find(&tokos).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		var result []models.TokoResponse
		for _, t := range tokos {
			result = append(result, models.TokoResponse{
				NamaToko: t.NamaToko,
				URLFoto: t.URLFoto,
				UserID: t.UserID,
				CreatedAt: t.CreatedAt,
				UpdatedAt: t.UpdatedAt,
			})
		}

		response := fiber.Map{"data": result}

		if pageStr != "" && limitStr != "" {
			page, _ := strconv.Atoi(pageStr)
			limit, _ := strconv.Atoi(limitStr)
			response["page"] = page
			response["limit"] = limit
			response["total"] = total
		}

		return c.Status(fiber.StatusOK).JSON(response)
	} 
}

func GetTokoByIdHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		id := c.Params("id")
		var toko models.Toko

		if err := db.First(&toko, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Cannot get toko",
				})
			}
		}

		response := models.TokoResponse{
			NamaToko: toko.NamaToko,
			URLFoto: toko.URLFoto,
			UserID: toko.UserID,
			CreatedAt: toko.CreatedAt,
			UpdatedAt: toko.UpdatedAt,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": response,
		})
	} 
}

func GetMyTokoHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		var toko models.Toko
		if err := db.Where("user_id = ?", userID).First(&toko).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Toko not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		response := models.TokoResponse{
			NamaToko: toko.NamaToko,
			URLFoto: toko.URLFoto,
			CreatedAt: toko.CreatedAt,
			UpdatedAt: toko.UpdatedAt,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": response,
		})
	}
}

func UpdateMyTokoHandlers(db *gorm.DB) fiber.Handler {
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

		idToko := c.Params("id_toko")

		var toko models.Toko
		if err := db.First(&toko, idToko).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Toko not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		if toko.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"message": "You don't have any access to update this toko",
			})
		}

		namaToko := c.FormValue("nama_toko")
		if namaToko != "" {
			toko.NamaToko = namaToko
		}

		file, err := c.FormFile("photo")
		if err == nil {
			path := "./uploads/toko_profile/"
			if _, err := os.Stat(path); os.IsNotExist(err) {
				os.MkdirAll(path, os.ModePerm)
			}

			filePath := path + file.Filename
			if err := c.SaveFile(file, filePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error",
					"message": "Fail to save the photo",
				})
			}

			toko.URLFoto = filePath
		}

		if err := db.Save(&toko).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		response := models.TokoResponse{
			NamaToko: toko.NamaToko,
			URLFoto: toko.URLFoto,
			UserID: toko.UserID,
			CreatedAt: toko.CreatedAt,
			UpdatedAt: toko.UpdatedAt,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": response,
		})
		
	}
}