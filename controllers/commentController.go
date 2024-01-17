package controllers

import (
	"fmt"
	"log"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// func GetAllComments(c *fiber.Ctx) error {

// }

// func GetCommentsByPostId(c *fiber.Ctx) error { viet lai getPosts ben postController
// 	postId, err := c.ParamsInt("id")
// 	if err != nil {
// 		return c.Status(401).SendString("Invalid id")
// 	}

// 	var comments []models.Comment

// 	if err := database.DB.

// }
func preload(d *gorm.DB) *gorm.DB {
	return d.Preload("Responses", preload)
}

func GetCommentsByPostId(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(401).SendString("Invalid id")
	}
	var comments []models.Comment

	err = database.DB.Where("post_id = ? AND response_id IS NULL", postID).Preload(clause.Associations, preload).Find(&comments).Error
	// err = database.DB.Where("post_id = ?", postID).Preload("Responses").Find(&comments).Error

	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(comments)
}

func CreateComment(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	// Make a Comment object
	comment := models.Comment{
		UserID: uint64(data["user_id"].(float64)),
		PostID: uint64(data["post_id"].(float64)),
		// ResponseID: &topComment.Id,
		Body: data["body"].(string),
	}

	// Find the comment whose id is entered for response_id
	if data["response_id"] != nil {
		var response_to uint64 = uint64(data["response_id"].(float64))
		var topComment models.Comment
		if err := database.DB.First(&topComment, "id = ?", response_to); err != nil {
			log.Println(err)
		}

		if topComment.Id == 0 {
			c.Status(400)
			return c.JSON(fiber.Map{
				"message": "Invalid Comment ID",
			})
		}

		comment.ResponseID = &topComment.Id
	}

	// Insert value
	if err := database.DB.Create(&comment); err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"comment": comment,
		"message": "Comment created successfully",
	})
}
