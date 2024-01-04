package routes

import (
	"github.com/DennieDan/movie-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("api/register", controllers.Register)
	app.Post("api/create_post", controllers.CreatePost)
}
