package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ListProvinciesHandler(c *fiber.Ctx) error {
	response, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	defer response.Body.Close()

	var data interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	return c.JSON(data)
}

func ListCitiesByProvinciesHandler(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	link := "https://emsifa.github.io/api-wilayah-indonesia/api/regencies/" + provID + ".json"
	response, err := http.Get(link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	defer response.Body.Close()

	var data interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	return c.JSON(data)
}

func ProvinceByIdHandlers(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	link := "https://emsifa.github.io/api-wilayah-indonesia/api/province/" + provID + ".json"
	response, err := http.Get(link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	defer response.Body.Close()

	var data interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	return c.JSON(data)
}

func CityByIdHandler(c *fiber.Ctx) error {
	cityID := c.Params("city_id")
	link := "https://emsifa.github.io/api-wilayah-indonesia/api/regency/" + cityID + ".json"
	response, err := http.Get(link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	defer response.Body.Close()

	var data interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": err.Error(),
		})
	}

	return c.JSON(data)
}