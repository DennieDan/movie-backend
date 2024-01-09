package routes

import (
	// "github.com/DennieDan/movie-backend/controllers"
	"github.com/DennieDan/movie-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// app.Post("api/register", controllers.Register)
	app.Post("api/create_post", controllers.CreatePost)
	app.Get("api/posts", controllers.GetPosts)
	app.Put("api/edit_post/:id", controllers.EditPost)
	// app.Get("api/posts/:id", controllers.GetPostById)
	app.Get("api/movies", controllers.GetMovies)
	app.Get("api/topics", controllers.GetTopics)
	app.Get("api/users", controllers.GetUsers)
}
