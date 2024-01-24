package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/DennieDan/movie-backend/database"
	"github.com/DennieDan/movie-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// func GetAllComments(c *fiber.Ctx) error {

// }

func GetTotalCommentsByPostId(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(401).SendString("Invalid id")
	}
	var comments []models.Comment

	err = database.DB.Where("post_id = ?", postID).Find(&comments).Error
	if err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message":        "Get total comments successfully",
		"total_comments": len(comments),
	})

}

type CommentRetrieve struct {
	Id         uint64 `gorm:"primary_key;auto_increment" json:"id"`
	UserID     uint64 `json:"user_id"`
	PostID     uint64 `json:"post_id"`
	ResponseID *uint64
	Body       string            `gorm:"text;not null;" json:"content"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `gorm:"index"`
	Responses  []CommentRetrieve `gorm:"foreignKey:ResponseID;constraint:OnDelete:CASCADE;"`
	User       string            `json:"user"`
	Voters     []*models.User    `gorm:"many2many:comment_votes;" json:"voters"`
}

func preload(d *gorm.DB) *gorm.DB {
	return d.Preload("Responses", preload)

}

func GetCommentsByPostId(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(401).SendString("Invalid id")
	}
	// replies := []CommentRetrieve{}
	var comments []models.Comment

	err = database.DB.
		// Where("comments.post_id = ? AND comments.response_id IS NULL", postID).
		Where("comments.post_id = ?", postID).
		Preload(clause.Associations).
		Find(&comments).Error
	// err = database.DB.Table("comments").
	// 	Joins("INNER JOIN users ON users.id = comments.user_id").
	// 	Where("post_id = ? AND response_id IS NULL", postID).
	// 	Select("comments.*, users.username AS user").
	// 	Preload(clause.Associations, preload).
	// 	Find(&replies).Error
	if err != nil {
		c.Status(400)
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(comments)
}

func CreateComment(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	// Make a Comment object
	comment := models.Comment{
		UserID: uint64(data["user_id"].(float64)),
		PostID: uint64(data["post_id"].(float64)),
		// ResponseID: &topComment.Id,
		Body: data["body"].(string),
	}

	// Find the comment whose id is entered for response_id
	if data["response_id"] != nil {
		var response_to uint64 = uint64(data["response_id"].(float64))
		var topComment models.Comment
		if err := database.DB.First(&topComment, "id = ?", response_to); err != nil {
			log.Println(err)
		}

		if topComment.Id == 0 {
			c.Status(400)
			return c.JSON(fiber.Map{
				"message": "Invalid Comment ID",
			})
		}

		comment.ResponseID = &topComment.Id
	}

	// Insert value
	if err := database.DB.Create(&comment); err != nil {
		log.Println(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"comment": comment,
		"message": "Comment created successfully",
	})
}

func EditComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).SendString("Not Found ID")
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	var comment models.Comment
	if err := database.DB.First(&comment, uint64(id)).Error; err != nil {
		log.Println(err)
		return c.Status(401).SendString(err.Error())
	}

	comment.Body = data["body"].(string)

	err = database.DB.Save(&comment).Error
	if err != nil {
		log.Println(err)
	}

	c.Status(200) // OK
	return c.JSON(fiber.Map{
		"comment": comment,
		"message": "Comment edited successfully",
	})
}

func DeleteComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).SendString("Not found ID")
	}

	if err := database.DB.Where("id = ?", id).Delete(&models.Comment{}).Error; err != nil {
		log.Println(err)
	}

	c.Status(204) // No Content
	return c.JSON(fiber.Map{
		"message": "Comment deleted successfully",
	})
}

func addCommentVoter(user_id int, comment_id int) error {
	var user models.User
	if err := database.DB.First(&user, uint64(user_id)).Error; err != nil {
		return err
	}

	var comment models.Comment
	if err := database.DB.First(&comment, uint64(comment_id)).Error; err != nil {
		return err
	}

	// comment.Voters = append(comment.Voters, &user)
	// if err := database.DB.Save(&comment).Error; err != nil {
	// 	return err
	// }

	comment_vote := models.CommentVotes{
		UserID:    uint64(user_id),
		CommentID: uint64(comment_id),
		Score:     0,
	}

	err := database.DB.Create(&comment_vote)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func UpvoteComment(c *fiber.Ctx) error {
	// Get user_id and comment_id from parameters
	user_id, err := c.ParamsInt("user")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	comment_id, err := c.ParamsInt("comment")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	// query the join table to check if the user has already voted
	var comment_vote models.CommentVotes
	database.DB.Where("comment_id = ? AND user_id = ?", comment_id, user_id).First(&comment_vote)

	// If the record is not found, create the reaction into the join table
	if comment_vote.UserID == 0 {
		err := addCommentVoter(user_id, comment_id)
		if err != nil {
			log.Println(err)
			c.Status(400)
			return c.JSON(err)
		}
	}

	// Get the related record from the join table and update
	database.DB.Where("comment_id = ? AND user_id = ?", comment_id, user_id).First(&comment_vote)
	if comment_vote.Score != 1 {
		comment_vote.Score = 1
	} else {
		comment_vote.Score = 0
	}
	if err := database.DB.Save(&comment_vote).Error; err != nil {
		log.Println(err)
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Comment Upvoted successfully",
	})
}

func DownvoteComment(c *fiber.Ctx) error {
	// Get user_id and comment_id from parameters
	user_id, err := c.ParamsInt("user")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	comment_id, err := c.ParamsInt("comment")
	if err != nil {
		return c.Status(401).SendString("Invalid ID")
	}

	// query the join table to check if the user has already voted
	var comment_vote models.CommentVotes
	database.DB.Where("comment_id = ? AND user_id = ?", comment_id, user_id).First(&comment_vote)

	// If the record is not found, create the reaction into the join table
	if comment_vote.UserID == 0 {
		err := addCommentVoter(user_id, comment_id)
		if err != nil {
			log.Println(err)
			c.Status(400)
			return c.JSON(err)
		}
	}

	// Get the related record from the join table and update
	database.DB.Where("comment_id = ? AND user_id = ?", comment_id, user_id).First(&comment_vote)
	if comment_vote.Score != -1 {
		comment_vote.Score = -1
	} else {
		comment_vote.Score = 0
	}
	if err := database.DB.Save(&comment_vote).Error; err != nil {
		log.Println(err)
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"message": "Comment Downvoted successfully",
	})
}
