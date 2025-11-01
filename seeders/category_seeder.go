package seeders

import (
	"go_evermos_rakamin_irsan/models"
	"log"

	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) error {
	categories := []models.Category{
		{NamaCategory: "Baju Muslim"},
		{NamaCategory: "Baju Muslimah"},
		{NamaCategory: "Peci"},
		{NamaCategory: "Jilbab"},
		{NamaCategory: "Sarung"},
	}

	for _, cat := range categories {
		if err := db.FirstOrCreate(&cat, models.Category{NamaCategory: cat.NamaCategory}).Error; err != nil {
			return err
		}
	}

	log.Println("Category Created")
	return nil
}