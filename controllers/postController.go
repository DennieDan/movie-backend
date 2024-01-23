package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
		// MovieID:  uint64(data["movie_id"].(float64)), // convert from float64 to uint32 as specified in the struct
		TopicID:  uint64(data["topic_id"].(float64)),
		AuthorID: uint64(data["author_id"].(float64)), // do authentication
		// created_at and updated_at will be handled by mysql
	}

	var movie models.Movie
	if err := database.DB.First(&movie, uint64(data["movie_id"].(float64))).Error; err != nil {
		log.Println(err)
		c.Status(400)
		return c.JSON(fiber.Map{
			"error": err,
		})
	}

	if data["movie_id"] != nil {
		post.MovieID = &movie.Id
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

type PostRetrieve struct {
	Id        uint64           `gorm:"primary_key;auto_increment" json:"id"`
	Title     string           `gorm:"size:255;not null;unique" json:"title"`
	Content   string           `gorm:"text;not null;" json:"content"`
	MovieID   uint64           `gorm:"default:null" json:"movie_id"`
	TopicID   uint64           `json:"topic_id"`
	AuthorID  uint64           `gorm:"not null;" json:"author_id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index"`
	Movie     models.Movie     `json:"movie"`
	Topic     models.Topic     `json:"topic"`
	Author    models.User      `gorm:"foreignKey:AuthorID;references:Id" json:"author"`
	Comments  []models.Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;" json:"comments"`
	Votes     int              `json:"votes"`
}

func GetPosts(c *fiber.Ctx) error {
	// posts := []models.Post{}
	posts := []PostRetrieve{}
	// err := database.DB.Joins("Movie").Joins("Topic").Joins("Author").Preload("Comments").Preload("Voters").Preload("Savers").Find(&posts).Error
	// err := database.DB.Joins("Movie").Joins("Topic").Joins("Author").Preload("Comments").Joins("Voters").Preload("Savers").Find(&posts).Error

	query := database.DB.Table("posts").
		Joins("LEFT JOIN movies ON movies.id = posts.movie_id").
		Joins("LEFT JOIN topics ON topics.id = posts.topic_id").
		Joins("LEFT JOIN users ON users.id = posts.author_id").
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id").
		Joins("Movie").Joins("Topic").Joins("Author").Preload("Comments").
		Select("posts.*, movies.title AS movie, topics.name AS topic, users.id AS author, SUM(post_votes.score) as votes").
		Group("posts.id")

	err := query.Find(&posts)

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
	var movie models.Movie
	if err := database.DB.First(&movie, uint64(data["movie_id"].(float64))).Error; err != nil {
		log.Println(err)
		c.Status(400)
		return c.JSON(fiber.Map{
			"error": err,
		})
	}
	post.Title = data["title"].(string)
	post.Content = data["content"].(string)
	post.MovieID = &movie.Id
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

func DeletePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).SendString("Not found ID")
	}

	if err := database.DB.Where("id = ?", id).Delete(&models.Post{}).Error; err != nil {
		log.Println(err)
	}

	c.Status(204) // No Content
	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}

func SavePost(c *fiber.Ctx) error {
	user_id, err := c.ParamsInt("user")
	if err != nil {
		return c.Status(404).SendString("No user found")
	}

	post_id, err := c.ParamsInt("post")
	if err != nil {
		return c.Status(404).SendString("No post found")
	}

	var user models.User
	var post models.Post
	if err := database.DB.First(&user, uint64(user_id)).Error; err != nil {
		log.Println(err)
		return c.Status(401).SendString(err.Error())
	}
	if err := database.DB.First(&post, uint64(post_id)).Error; err != nil {
		log.Println(err)
		return c.Status(401).SendString(err.Error())
	}

	user.SavedPosts = append(user.SavedPosts, &post)
	post.Savers = append(post.Savers, &user)

	err = database.DB.Save(&user).Error
	if err != nil {
		log.Println(err)
	}
	err = database.DB.Save(&post).Error
	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Post saved successfully",
	})
}

func AddVoter(user_id int, post_id int) error {
	var user models.User
	if err := database.DB.First(&user, uint64(user_id)).Error; err != nil {
		return err
	}

	var post models.Post
	if err := database.DB.First(&post, uint64(post_id)).Error; err != nil {
		return err
	}

	post.Voters = append(post.Voters, &user)
	if err := database.DB.Save(&post).Error; err != nil {
		return err
	}
	return nil
}

// Don't use post without created_at and updated_at
func UpvotePost(c *fiber.Ctx) error {
	// Get user_id and post_id from parameters
	user_id, err := c.ParamsInt("user")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	post_id, err := c.ParamsInt("post")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	// query the join table to check if the user has already voted
	var post_vote models.PostVotes
	database.DB.Where("post_id = ? AND user_id = ?", post_id, user_id).First(&post_vote)

	// If the record is not found, create the reaction into the join table
	if post_vote.UserID == 0 {
		err := AddVoter(user_id, post_id)
		if err != nil {
			log.Println(err)
			c.Status(400)
			return c.JSON(err)
		}
	}

	// Get the related record from the join table and update
	database.DB.Where("post_id = ? AND user_id = ?", post_id, user_id).First(&post_vote)
	if post_vote.Score != 1 {
		post_vote.Score = 1
	} else {
		post_vote.Score = 0
	}
	if err := database.DB.Save(&post_vote).Error; err != nil {
		log.Println(err)
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	// if post_vote.PostID != 0 {
	// 	post_vote.Score = 1
	// 	if err := database.DB.Save(&post_vote).Error; err != nil {
	// 		log.Println(err)
	// 	}
	// } else {
	// 	var user models.User
	// 	if err := database.DB.First(&user, uint64(user_id)).Error; err != nil {
	// 		log.Println(err)
	// 		return c.Status(401).SendString(err.Error())
	// 	}
	// 	post.Voters = append(post.Voters, &user)
	// 	if err := database.DB.Save(&post).Error; err != nil {
	// 		log.Println(err)
	// 	}

	// 	database.DB.Where("post_id = ? AND user_id = ?", post_id, user_id).First(&post_vote)
	// 	if post_vote.PostID != 0 {
	// 		post_vote.Score = 1
	// 		if err := database.DB.Save(&post_vote).Error; err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// }

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Upvoted successfully",
	})

}

func DownvotePost(c *fiber.Ctx) error {
	// Get user_id and post_id from parameters
	user_id, err := c.ParamsInt("user")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	post_id, err := c.ParamsInt("post")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	// query the join table to check if the user has already voted
	var post_vote models.PostVotes
	database.DB.Where("post_id = ? AND user_id = ?", post_id, user_id).First(&post_vote)

	// If the record is not found, create the reaction into the join table
	if post_vote.UserID == 0 {
		err := AddVoter(user_id, post_id)
		if err != nil {
			log.Println(err)
			c.Status(400)
			return c.JSON(err)
		}
	}

	// Get the related record from the join table and update
	database.DB.Where("post_id = ? AND user_id = ?", post_id, user_id).First(&post_vote)
	if post_vote.Score != -1 {
		post_vote.Score = -1
	} else {
		post_vote.Score = 0
	}
	if err := database.DB.Save(&post_vote).Error; err != nil {
		log.Println(err)
		c.Status(400)
		return c.JSON(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Downvoted successfully",
	})
}

func UnvotePost(c *fiber.Ctx) error {
	// Get user_id and post_id from parameters
	user_id, err := c.ParamsInt("user")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	post_id, err := c.ParamsInt("post")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	// query the join table to check if the user has already voted
	var post_vote models.PostVotes
	database.DB.Where("post_id = ? AND user_id = ?", post_id, user_id).First(&post_vote)

	// If the record is not found, create the reaction into the join table
	if post_vote.UserID == 0 {
		c.Status(401)
		return c.JSON(fiber.Map{
			"message": "No record found",
		})
	}

	post_vote.Score = 0

	if err := database.DB.Save(&post_vote).Error; err != nil {
		log.Println(err)
		c.Status(400)
		return c.JSON(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Unvoted successfully",
	})

}
