package migration

import (
	"go_evermos_rakamin_irsan/models"
	"go_evermos_rakamin_irsan/seeders"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Category{},
   		&models.Toko{},
    	&models.User{},
    	&models.Alamat{},
    	&models.Product{},
    	&models.Trx{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.FotoProduk{},
		&models.DetailTrx{},
	); err != nil {
		return err
	}
	log.Println("Database migrated")

	if err := FixLogFotoFK(db); err != nil {
		return err
	}
	log.Println("Foreign key between log_produk & foto_produk fixed")

	if err := seeders.SeedAdmin(db); err != nil {
		return err
	}

	if err := seeders.SeedCategories(db); err != nil {
		return err
	}

	log.Println("Admin seeded")
	return nil
	
}