package controllers

import (
	"log"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetTopics(c *fiber.Ctx) error {
	topics := []models.Topic{}
	err := database.DB.Model(&models.Topic{}).Preload("Posts").Find(&topics).Error
	if err != nil {
		log.Println(err)
	}

	c.Status(200)

	return c.JSON(topics)
}
