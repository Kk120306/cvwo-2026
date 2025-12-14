package database

/*
This file has been geenrated fully with chatgpt.
It is for code testing purposes and aims to populate
data to the database for development
*/

import (
	"log"
	"math/rand"
	"time"

	"github.com/Kk120306/cvwo-2026/backend/models"
	"gorm.io/gorm"
)

/*
Seed populates the database with initial data.
Safe to run multiple times.
*/
func Seed() {
	log.Println("🌱 Seeding database...")

	seedUsers(DB)
	seedTopics(DB)
	seedPosts(DB)
	seedComments(DB)
	seedVotes(DB)

	log.Println("✅ Database seeded successfully")
}

/* ===================== USERS ===================== */

func seedUsers(db *gorm.DB) {
	users := []models.User{
		{Username: "admin", IsAdmin: true},
		{Username: "alice"},
		{Username: "bob"},
		{Username: "charlie"},
	}

	for _, user := range users {
		var existing models.User
		err := db.Where("username = ?", user.Username).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			db.Create(&user)
			log.Println("👤 Created user:", user.Username)
		}
	}
}

/* ===================== TOPICS ===================== */

func seedTopics(db *gorm.DB) {
	topics := []models.Topic{
		{Name: "Technology", Slug: "technology"},
		{Name: "Programming", Slug: "programming"},
		{Name: "Artificial Intelligence", Slug: "ai"},
		{Name: "Startups", Slug: "startups"},
		{Name: "Life", Slug: "life"},
	}

	for _, topic := range topics {
		var existing models.Topic
		err := db.Where("slug = ?", topic.Slug).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			db.Create(&topic)
			log.Println("🏷 Created topic:", topic.Name)
		}
	}
}

/* ===================== POSTS ===================== */

func seedPosts(db *gorm.DB) {
	var topics []models.Topic
	var users []models.User

	db.Find(&topics)
	db.Find(&users)

	if len(topics) == 0 || len(users) == 0 {
		log.Println("⚠️ Skipping post seeding")
		return
	}

	posts := []models.Post{
		{
			Title:    "Welcome to the Forum",
			Content:  "This is the first pinned post of the platform.",
			IsPinned: true,
		},
		{
			Title:   "Best way to learn Go?",
			Content: "Share your resources, tips, and experiences learning Go.",
		},
		{
			Title:   "Is AI replacing developers?",
			Content: "An open discussion about AI and the future of software jobs.",
		},
	}

	for i, post := range posts {
		post.TopicID = topics[i%len(topics)].ID
		post.AuthorID = users[i%len(users)].ID

		var existing models.Post
		err := db.
			Where("title = ? AND topic_id = ?", post.Title, post.TopicID).
			First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			db.Create(&post)
			log.Println("📝 Created post:", post.Title)
		}
	}
}

/* ===================== COMMENTS ===================== */

func seedComments(db *gorm.DB) {
	var posts []models.Post
	var users []models.User

	db.Find(&posts)
	db.Find(&users)

	if len(posts) == 0 || len(users) == 0 {
		log.Println("⚠️ Skipping comment seeding")
		return
	}

	comments := []string{
		"Great post!",
		"I completely agree with this.",
		"Very insightful.",
		"Thanks for sharing your thoughts.",
	}

	for i, content := range comments {
		comment := models.Comment{
			PostID:   posts[i%len(posts)].ID,
			AuthorID: users[(i+1)%len(users)].ID,
			Content:  content,
		}

		var existing models.Comment
		err := db.
			Where("post_id = ? AND content = ?", comment.PostID, comment.Content).
			First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			db.Create(&comment)
			log.Println("💬 Created comment")
		}
	}
}

/* ===================== VOTES ===================== */

func seedVotes(db *gorm.DB) {
	var posts []models.Post
	var users []models.User

	db.Find(&posts)
	db.Find(&users)

	if len(posts) == 0 || len(users) == 0 {
		log.Println("⚠️ Skipping vote seeding")
		return
	}

	rand.Seed(time.Now().UnixNano())

	for _, post := range posts {
		for _, user := range users {

			// 50% chance user votes
			if rand.Intn(2) == 0 {
				continue
			}

			vote := models.Vote{
				UserID:      user.ID,
				VotableID:   post.ID,
				VotableType: "post",
				VoteType:    []string{"like", "dislike"}[rand.Intn(2)],
			}

			var existing models.Vote
			err := db.
				Where(
					"user_id = ? AND votable_id = ? AND votable_type = ?",
					user.ID, post.ID, "post",
				).
				First(&existing).Error

			if err == gorm.ErrRecordNotFound {
				db.Create(&vote)
				log.Println("👍 Created vote on post")
			}
		}
	}
}
