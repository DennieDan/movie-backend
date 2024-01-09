package controllers

import (
	"log"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetMovies(c *fiber.Ctx) error {
	movies := []models.Movie{}
	err := database.DB.Model(&models.Movie{}).Preload("Posts").Find(&movies).Error
	if err != nil {
		log.Println(err)
	}

	c.Status(200)

	return c.JSON(movies)
}
