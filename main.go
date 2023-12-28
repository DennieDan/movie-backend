package main

import (
	"log"
	"os"

	"github.com/DennieDan/movie-backend/database"
	"github.com/gofiber/fiber/v2"
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
	// routes.Setup(app)
	app.Listen(":" + port)
}
