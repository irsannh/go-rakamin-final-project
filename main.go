package main

import (
	"go_evermos_rakamin_irsan/config"
	"go_evermos_rakamin_irsan/migration"
	"go_evermos_rakamin_irsan/routes"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	app := config.NewFiber(cfg)

	db := config.InitDB(cfg)

	if err := migration.Migrate(db); err != nil {
		log.Fatal("Migration Failed:", err)
	}

	routes.SetupRoutes(app, db, cfg.JwtSecret)

	app.Listen("localhost:" + cfg.AppPort)
}