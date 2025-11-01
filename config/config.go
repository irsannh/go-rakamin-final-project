package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName string
	AppPort string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	JwtSecret string
}

func LoadConfig() *AppConfig {
	_ = godotenv.Load()

	return &AppConfig{
		AppName: os.Getenv("APP_NAME"),
		AppPort: os.Getenv("APP_PORT"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}