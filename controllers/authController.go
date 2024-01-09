package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}
	err := database.DB.Preload("Posts").Preload("Comments").Preload("VoteComments").Preload("VotePosts").Preload("SavedPosts").Find(&users)
	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(users)
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z0-9._%+\-]+`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	// Check email
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid Email address",
		})
	}

	// Check if email already exist in database
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Check if password is less than 6 characters
	if len(data["password"].(string)) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Password must be greater than 6 characters",
		})
	}

	user := models.User{
		Username:   data["username"].(string),
		Email:      strings.TrimSpace(data["email"].(string)),
		AvatarPath: "",
	}

	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfully",
	})
}
