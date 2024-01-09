package controllers

import (
	"fmt"
	"log"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
)

func CreatePost(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	// Check ve null la do front end
	post := models.Post{
		Title:   data["title"].(string),
		Content: data["content"].(string),
		// MovieID:  data["movie_id"].(float64), // lay MovieID va TopicID la do frontend
		// TopicID:  data["topic_id"].(float64),
		// AuthorID: data["topic_id"].(float64),
		MovieID:  uint64(data["movie_id"].(float64)), // convert from float64 to uint32 as specified in the struct
		TopicID:  uint64(data["topic_id"].(float64)),
		AuthorID: uint64(data["author_id"].(float64)), // do authentication
		// created_at and updated_at will be handled by mysql
	}

	err := database.DB.Create(&post)
	if err != nil {
		log.Println(err)
	}

	c.Status(201) // created
	return c.JSON(fiber.Map{
		"post":    post,
		"message": "Post created successfully",
	})

}

func GetPosts(c *fiber.Ctx) error {
	posts := []models.Post{}
	err := database.DB.Joins("Movie").Joins("Topic").Joins("Author").Preload("Comments").Preload("Voters").Preload("Savers").Find(&posts).Error

	// query := database.DB.Table("posts").
	// 	Joins("INNER JOIN movies ON movies.id = posts.movie_id").
	// 	Joins("INNER JOIN topics ON topics.id = posts.topic_id").
	// 	Joins("INNER JOIN users ON users.id = posts.author_id").
	// 	Select("posts.*, movies.title AS movie, topics.name AS topic, users.id AS author")

	// err := query.Find(&posts)

	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(posts)
}

// func GetPostById(c *fiber.Ctx) error {
// 	id, err := c.ParamsInt("id")
// 	log.Println(id)

// 	if err != nil {
// 		return c.Status(401).SendString("Invalid id")
// 	}

// 	post := &models.Post{}

// 	database.DB.First(&post, id)
// 	return c.JSON(post)
// }

func EditPost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(401).SendString("Invalid id")
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	var post models.Post
	if err := database.DB.First(&post, uint64(id)).Error; err != nil {
		log.Println(err)
		return c.Status(401).SendString(err.Error())
	}
	post.Title = data["title"].(string)
	post.Content = data["content"].(string)
	post.MovieID = uint64(data["movie_id"].(float64))
	post.TopicID = uint64(data["topic_id"].(float64))

	err = database.DB.Save(&post).Error

	if err != nil {
		log.Println(err)
	}

	c.Status(200) // OK
	return c.JSON(fiber.Map{
		"post":    post,
		"message": "Post edited successfully",
	})

}
