package handlers

import (
	"go_evermos_rakamin_irsan/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllCategoriesHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		var categories []models.Category
		if err := db.Find(&categories).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Cannot get categories",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "success",
				"categories": categories,
			})
	} 
}

func GetCategoryByIdHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		id := c.Params("id")
		var category models.Category

		if err := db.First(&category, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Cannot get category",
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"category": category,
		})
	} 
}

func PostCategoryHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		var input models.CategoryRequest

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid Input",
			})
		}

		category := models.Category{NamaCategory: input.NamaCategory}
		if err := db.FirstOrCreate(&category, models.Category{NamaCategory: input.NamaCategory}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Cannot create category",
			})
		}
		
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "success",
			"category": category,
		})
	} 
}

func PutCategoryHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		id := c.Params("id")
		categoryID, _ := strconv.Atoi(id)
		var category models.Category

		if err := db.First(&category, categoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Category not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error",
					"message": "Cannot fetch category",
				})
		}

		var input models.CategoryRequest
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid Input",
			})
		}

		category.NamaCategory = input.NamaCategory

		if err := db.Save(&category).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Cannot update category",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "success",
			"category": category, 
		})
	} 
}

func DeleteCategoryHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		id := c.Params("id")
		categoryID, _ := strconv.Atoi(id)
		var category models.Category

		if err := db.First(&category, categoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Category not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error",
					"message": "Cannot fetch category",
				})
		}

		if err := db.Delete(&category).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Cannot delete category",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"message": "Category deleted",
		})
	} 
}