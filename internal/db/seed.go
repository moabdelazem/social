package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/moabdelazem/social/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe",
}

var titles = []string{
	"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
	"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
	"Home Office Setup", "Digital Detox", "Gardening Basics",
	"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
	"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
	"Fitness at Home", "Personal Finance Tips", "Creative Writing",
	"Mental Health Awareness", "Learning New Skills",
}

var contents = []string{
	"In this post, we'll explore how to develop good habits that stick and transform your life.",
	"Discover the benefits of a minimalist lifestyle and how to declutter your home and mind.",
	"Learn practical tips for eating healthy on a budget without sacrificing flavor.",
	"Traveling doesn't have to be expensive. Here are some tips for seeing the world on a budget.",
	"Mindfulness meditation can reduce stress and improve your mental well-being. Here's how to get started.",
	"Increase your productivity with these simple and effective strategies.",
	"Set up the perfect home office to boost your work-from-home efficiency and comfort.",
	"A digital detox can help you reconnect with the real world and improve your mental health.",
	"Start your gardening journey with these basic tips for beginners.",
	"Transform your home with these fun and easy DIY projects.",
	"Yoga is a great way to stay fit and flexible. Here are some beginner-friendly poses to try.",
	"Sustainable living is good for you and the planet. Learn how to make eco-friendly choices.",
	"Master time management with these tips and get more done in less time.",
	"Nature has so much to offer. Discover the benefits of spending time outdoors.",
	"Whip up delicious meals with these simple and quick cooking recipes.",
	"Stay fit without leaving home with these effective at-home workout routines.",
	"Take control of your finances with these practical personal finance tips.",
	"Unleash your creativity with these inspiring writing prompts and exercises.",
	"Mental health is just as important as physical health. Learn how to take care of your mind.",
	"Learning new skills can be fun and rewarding. Here are some ideas to get you started.",
}

var tags = []string{
	"Self Improvement", "Minimalism", "Health", "Travel", "Mindfulness",
	"Productivity", "Home Office", "Digital Detox", "Gardening", "DIY",
	"Yoga", "Sustainability", "Time Management", "Nature", "Cooking",
	"Fitness", "Personal Finance", "Writing", "Mental Health", "Learning",
}

var comments = []string{
	"Great post! Thanks for sharing.",
	"I completely agree with your thoughts.",
	"Thanks for the tips, very helpful.",
	"Interesting perspective, I hadn't considered that.",
	"Thanks for sharing your experience.",
	"Well written, I enjoyed reading this.",
	"This is very insightful, thanks for posting.",
	"Great advice, I'll definitely try that.",
	"I love this, very inspirational.",
	"Thanks for the information, very useful.",
}

// Seed populates the database with test data
func Seed(s store.Storage, db *sql.DB) error {
	ctx := context.Background()

	// Create users
	var userIDs []int64
	for _, username := range usernames {
		user := store.User{
			Username: username,
			Email:    fmt.Sprintf("%s@example.com", username),
			Password: "hashed_password_" + username,
		}
		if err := s.UsersRepo.Create(ctx, &user); err != nil {
			return fmt.Errorf("failed to create user %s: %w", username, err)
		}
		userIDs = append(userIDs, user.ID)
		log.Printf("Created user: %s (ID: %d)", user.Username, user.ID)
	}

	// Create posts
	var postIDs []int64
	for i, title := range titles {
		// Assign random user to each post
		userID := userIDs[rand.Intn(len(userIDs))]

		// Create random tags (2-4 tags per post)
		numTags := rand.Intn(3) + 2
		postTags := make([]string, numTags)
		usedTags := make(map[int]bool)
		for j := 0; j < numTags; j++ {
			tagIndex := rand.Intn(len(tags))
			for usedTags[tagIndex] {
				tagIndex = rand.Intn(len(tags))
			}
			usedTags[tagIndex] = true
			postTags[j] = tags[tagIndex]
		}

		post := store.Post{
			UserID:  userID,
			Title:   title,
			Content: contents[i%len(contents)],
			Tags:    postTags,
		}
		if err := s.PostsRepo.Create(ctx, &post); err != nil {
			return fmt.Errorf("failed to create post %s: %w", title, err)
		}
		postIDs = append(postIDs, post.ID)
		log.Printf("Created post: %s (ID: %d)", post.Title, post.ID)
	}

	// Create comments (2-5 comments per post)
	for _, postID := range postIDs {
		numComments := rand.Intn(4) + 2
		for i := 0; i < numComments; i++ {
			userID := userIDs[rand.Intn(len(userIDs))]
			commentText := comments[rand.Intn(len(comments))]

			query := `
				INSERT INTO comments (post_id, user_id, content, created_at)
				VALUES ($1, $2, $3, NOW())
			`
			_, err := db.ExecContext(ctx, query, postID, userID, commentText)
			if err != nil {
				return fmt.Errorf("failed to create comment for post %d: %w", postID, err)
			}
		}
	}

	log.Println("✅ Database seeding completed successfully!")
	log.Printf("Created %d users, %d posts, and multiple comments", len(userIDs), len(postIDs))
	return nil
}
