package controllers

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func Login(c *fiber.Ctx) error {
	// Get username and pass off the request body
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		fmt.Println("Unable to parse body")
	}

	// Look up the requested user via username
	var user models.User
	database.DB.First(&user, "username = ?", body["username"])
	if user.Id == 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Compare sent in pass with saved user pass hash
	// err := bcrypt.CompareHashAndPassword(user.Password, []byte(body["password"].(string)))
	if err := user.ComparePassword(body["password"].(string)); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"error": "Incorrect Password",
		})
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,                                    // subject
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // expire in 30 days
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	// send it back (in JSON response) or (in Cookie -> better 25:23)
	c.Status(200)

	// return c.JSON(fiber.Map{
	// 	"token": tokenString,
	// })
	cookie := fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "You are successfully logged in",
		"user":    user,
	})
}

func Validate(c *fiber.Ctx) error {
	user := c.Locals("user") // get user from context which is assigned in requireAuth
	c.Status(200)
	return c.JSON(fiber.Map{
		"message": user,
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set the expiry time an hour ago
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	}

	c.Cookie(&cookie)

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
