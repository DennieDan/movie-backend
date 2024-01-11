package main

import (
	"log"
	"os"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Utilize Connect() function in database package
	database.Connect()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	// Launch in port
	port := os.Getenv("PORT")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
		// AllowOrigins: []string{"http://localhost:4000"},
		// AllowHeaders: []string{"Origin, Content-Type, Accept"},
	}))

	app.Get("/posts", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	routes.Setup(app)
	log.Fatal(app.Listen(":" + port))
}
