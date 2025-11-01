package handlers

import "github.com/gofiber/fiber/v2"

func WelcomeHandlers(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to Rakamin Evermos Back-End API E-Commerce Project, made by Irsan Nur Hidayat",
	})
}