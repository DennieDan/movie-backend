package routes

import (
	// "github.com/DennieDan/movie-backend/controllers"
	"github.com/DennieDan/movie-backend/controllers"
	"github.com/DennieDan/movie-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("api/register", controllers.Register)
	app.Post("api/login", controllers.Login)
	app.Get("api/validate", middleware.RequireAuth, controllers.Validate)
	app.Post("api/logout", controllers.Logout)
	app.Put("api/update_account/:id", controllers.UpdateAccount)
	app.Put("api/change_password/:id", controllers.ChangePassword)

	app.Post("api/create_post", controllers.CreatePost)
	app.Get("api/posts", controllers.GetPosts)
	app.Put("api/edit_post/:id", controllers.EditPost)
	app.Delete("api/delete_post/:id", controllers.DeletePost)
	// app.Get("api/posts/:id", controllers.GetPostById)
	app.Get("api/movies", controllers.GetMovies)
	app.Get("api/topics", controllers.GetTopics)
	app.Get("api/users", controllers.GetUsers)
	app.Post("api/create_comment", controllers.CreateComment)
	app.Get("api/comments/:id", controllers.GetCommentsByPostId)
	app.Get("api/total_comments/:id", controllers.GetTotalCommentsByPostId)
	app.Put("api/edit_comment/:id", controllers.EditComment)
	app.Delete("api/delete_comment/:id", controllers.DeleteComment)

	app.Put("api/save_post/:user/:post", controllers.SavePost)
	app.Post("api/upvote_post/:user/:post", controllers.UpvotePost)
	app.Post("api/downvote_post/:user/:post", controllers.DownvotePost)
	app.Patch("api/unvote_post/:user/:post", controllers.UnvotePost)
}
