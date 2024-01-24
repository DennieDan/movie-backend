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
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	dsn := db_user + ":" + db_password + "@tcp(" + db_host + ":" + db_port + ")/gomovieforumtest?charset=utf8mb4&parseTime=True&loc=Local"
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

	// Initialize data in the local database
	initializeData()

}

func initializeData() {

	movies := []*models.Movie{
		{Title: "Once Upon a Time in the West", Year: 1968},
		{Title: "American History X", Year: 1998},
		{Title: "Interstellar", Year: 2014},
		{Title: "Casablanca", Year: 1942},
		{Title: "City Lights", Year: 1931},
		{Title: "Psycho", Year: 1960},
		{Title: "The Green Mile", Year: 1999},
		{Title: "The Intouchables", Year: 2011},
		{Title: "Modern Times", Year: 1936},
		{Title: "Raiders of the Lost Ark", Year: 1981},
		{Title: "Rear Window", Year: 1954},
		{Title: "The Pianist", Year: 2002},
		{Title: "The Departed", Year: 2006},
		{Title: "Cinema Paradiso", Year: 1988},
		{Title: "The Lives of Others", Year: 2006},
		{Title: "Grave of the Fireflies", Year: 1988},
		{Title: "Jaws", Year: 1975},
		{Title: "Home Alone", Year: 1990},
		{Title: "Soul", Year: 2020},
	}

	topics := []*models.Topic{
		{
			Name:  "review",
			Color: "default",
		},
		{
			Name:  "spoiler",
			Color: "primary",
		},
		{
			Name:  "Bad movie",
			Color: "secondary",
		},
		{
			Name:  "ask",
			Color: "info",
		},
		{
			Name:  "Netflix",
			Color: "success",
		},
		{
			Name:  "TV Show",
			Color: "warning",
		},
		{
			Name:  "Report",
			Color: "error",
		},
	}

	DB.Create(movies)
	DB.Create(topics)

	user1 := models.User{Username: "User1", Email: "dedeui@gmail.com"}
	user1.SetPassword("password111")
	DB.Create(&user1)

	user2 := models.User{Username: "User2", Email: "hurrayyygoing@hotmail.com"}
	user2.SetPassword("huhuhucry321")
	DB.Create(&user2)

	user3 := models.User{Username: "User3", Email: "funand4587@outlook.com"}
	user3.SetPassword("thisisapassword")
	DB.Create(&user3)

	post1 := models.Post{Title: "Is Jaws worth watching 20 times?", Content: "I was really shocked as I realised today was the twentieth time I have watched Jaws!", TopicID: 4, AuthorID: 1}
	post2 := models.Post{Title: "Souls itself is not the best movie for everyone, but it is a gift for anyone who is super duper alone.", Content: "I watched Souls today with my girl friend and the title is what she told me after the film ended", TopicID: 1, AuthorID: 3}
	post3 := models.Post{Title: "Which series is suitable for family gathering?", Content: "My big family intend to sit together every month, so I am finding activities for our bonding session. Do you think Home Alone is a good option for our family this month?", TopicID: 4, AuthorID: 2}
	post4 := models.Post{Title: "This movie brings me to a great world the 1990s haha", Content: "I really suggest watch The Green Mile, especially for those people at my age and want to stay away from the hustle and bustle of the busy city. \nThe way people treated one another will make you think again about life.", TopicID: 1, AuthorID: 2}

	var movie1 models.Movie
	if err := DB.First(&movie1, uint64(17)).Error; err != nil {
		log.Println(err)
	}
	post1.MovieID = &movie1.Id
	var movie2 models.Movie
	if err := DB.First(&movie2, uint64(19)).Error; err != nil {
		log.Println(err)
	}
	post2.MovieID = &movie2.Id
	var movie4 models.Movie
	if err := DB.First(&movie4, uint64(7)).Error; err != nil {
		log.Println(err)
	}
	post4.MovieID = &movie4.Id

	if err := DB.Create(&post1); err != nil {
		log.Println(err)
	}
	DB.Create(&post2)
	DB.Create(&post3)
	DB.Create(&post4)

	comments := []*models.Comment{
		{UserID: 3, PostID: 1, Body: "This is a comment for post 1, please reply below"},
		{UserID: 2, PostID: 1, Body: "This is the second comment for post 1, please reply below"},
		{UserID: 2, PostID: 2, Body: "This is another comment for post 2, please reply below"},
		{UserID: 1, PostID: 3, Body: "This is a comment for post 3, please reply below"},
		{UserID: 2, PostID: 3, Body: "This is another comment for post 3, please reply below"},
	}

	DB.Create(comments)
}
