package database

import (
	"log"
	"os"

	"github.com/DennieDan/movie-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DSN")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	} else {
		log.Println("Connect successfully")
	}
	DB = database

	_ = DB.SetupJoinTable(&models.Post{}, "Voters", &models.PostVotes{})

	// Auto chuyen sang SQL code
	database.AutoMigrate(
		&models.Movie{},
		&models.Topic{},
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.CommentVotes{},
		&models.PostVotes{},
	)
}
