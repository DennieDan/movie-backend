package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *fiber.Ctx) error {
	// Get the cookie off request
	tokenString := c.Cookies("Authorization")

	// Decode/validate it

	// Parse takes the token string and a function for looking up the key.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "unauthenticated",
			})
		}

		// Find the user with the token sub (subject)
		var user models.User
		database.DB.First(&user, claims["sub"])

		if user.Id == 0 {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "unauthenticated",
			})
		}

		// Attach to request
		c.Locals("user", user)

		// Continue
		return c.Next()
	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
}
