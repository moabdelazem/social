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
	"walter", "xenia", "yasmin", "zoe", "alex", "bella", "chris", "diana",
	"ethan", "felicia", "gary", "holly", "isaac", "julia", "kyle", "lucy",
	"marcus", "nora", "owen", "paula", "quincy", "rita", "sam", "tina",
	"ulysses", "vera", "will", "xavier", "yara", "zachary", "aaron", "beth",
	"carlos", "debbie", "eli", "faith", "greg", "heather", "irene", "jack",
	"kate", "leo", "maya", "nick", "olivia", "paul", "queenie", "ruby",
	"scott", "teresa", "unity", "vincent", "willow", "xena", "york", "zelda",
}

var titles = []string{
	"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
	"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
	"Home Office Setup", "Digital Detox", "Gardening Basics",
	"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
	"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
	"Fitness at Home", "Personal Finance Tips", "Creative Writing",
	"Mental Health Awareness", "Learning New Skills", "Building Confidence",
	"Photography for Beginners", "The Art of Negotiation", "Public Speaking Tips",
	"Morning Routines That Work", "Evening Wind-Down Rituals", "Meal Prep Made Easy",
	"Zero Waste Living", "Indoor Plant Care", "Pet Training Basics",
	"Career Development Strategies", "Side Hustle Ideas", "Retirement Planning 101",
	"Building Strong Relationships", "Conflict Resolution", "Active Listening Skills",
	"Tech Gadgets Review", "Software Development Best Practices", "Cybersecurity Basics",
	"Data Science Fundamentals", "Machine Learning Introduction", "Cloud Computing Explained",
	"Leadership Principles", "Team Management", "Remote Work Success",
	"Book Recommendations", "Movie Reviews", "Music Discovery",
	"Art Appreciation", "Cultural Diversity", "Historical Insights",
	"Philosophy Basics", "Critical Thinking", "Problem Solving Techniques",
}

var contents = []string{
	"In this post, we'll explore how to develop good habits that stick and transform your life. Small changes lead to big results.",
	"Discover the benefits of a minimalist lifestyle and how to declutter your home and mind. Less is truly more.",
	"Learn practical tips for eating healthy on a budget without sacrificing flavor. Nutrition doesn't have to be expensive.",
	"Traveling doesn't have to be expensive. Here are some tips for seeing the world on a budget. Adventure awaits!",
	"Mindfulness meditation can reduce stress and improve your mental well-being. Here's how to get started with just 5 minutes a day.",
	"Increase your productivity with these simple and effective strategies. Work smarter, not harder.",
	"Set up the perfect home office to boost your work-from-home efficiency and comfort. Your workspace matters.",
	"A digital detox can help you reconnect with the real world and improve your mental health. Unplug to recharge.",
	"Start your gardening journey with these basic tips for beginners. Growing your own food is rewarding.",
	"Transform your home with these fun and easy DIY projects. Create something beautiful today.",
	"Yoga is a great way to stay fit and flexible. Here are some beginner-friendly poses to try. Namaste!",
	"Sustainable living is good for you and the planet. Learn how to make eco-friendly choices every day.",
	"Master time management with these tips and get more done in less time. Productivity is about priorities.",
	"Nature has so much to offer. Discover the benefits of spending time outdoors. Fresh air heals.",
	"Whip up delicious meals with these simple and quick cooking recipes. Anyone can cook!",
	"Stay fit without leaving home with these effective at-home workout routines. No gym needed.",
	"Take control of your finances with these practical personal finance tips. Financial freedom is possible.",
	"Unleash your creativity with these inspiring writing prompts and exercises. Every story starts somewhere.",
	"Mental health is just as important as physical health. Learn how to take care of your mind. You're not alone.",
	"Learning new skills can be fun and rewarding. Here are some ideas to get you started. Never stop growing.",
	"Building confidence takes practice. Here are strategies that actually work. Believe in yourself.",
	"Photography basics for beginners. Capture moments that matter with these simple techniques.",
	"Negotiation is an essential life skill. Learn how to advocate for yourself effectively.",
	"Public speaking doesn't have to be scary. Tips to overcome stage fright and deliver with confidence.",
	"Morning routines can set the tone for your entire day. Find what works for you.",
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
	"This changed my perspective completely!",
	"Bookmarking this for later reference.",
	"Could you elaborate more on this point?",
	"This is exactly what I needed to hear today.",
	"Brilliant explanation, thank you!",
	"I've been doing this wrong all along!",
	"Can't wait to try this out.",
	"This deserves more attention!",
	"You just earned a new follower.",
	"Mind blown! Thanks for sharing.",
	"This is gold, thank you!",
	"I wish I had known this earlier.",
	"Sharing this with my friends!",
	"This made my day better.",
	"Perfect timing for this post!",
}

// Seed populates the database with test data using a transaction
func Seed(s store.Storage, db *sql.DB) error {
	ctx := context.Background()

	return store.WithTx(db, ctx, func(tx *sql.Tx) error {
		// Create users
		userIDs, err := seedUsers(ctx, s)
		if err != nil {
			return err
		}
		log.Printf("Created %d users", len(userIDs))

		// Create posts
		postIDs, err := seedPosts(ctx, s, userIDs)
		if err != nil {
			return err
		}
		log.Printf("Created %d posts", len(postIDs))

		// Create comments
		if err := seedComments(ctx, tx, postIDs, userIDs); err != nil {
			return err
		}
		log.Printf("Created comments for all posts")

		// Create followers (users follow random other users)
		if err := seedFollowers(ctx, tx, userIDs); err != nil {
			return err
		}
		log.Printf("Created follower relationships")

		log.Println("Database seeding completed successfully!")
		return nil
	})
}

func seedUsers(ctx context.Context, s store.Storage) ([]int64, error) {
	userIDs := make([]int64, 0, len(usernames))

	for _, username := range usernames {
		user := store.User{
			Username: username,
			Email:    fmt.Sprintf("%s@example.com", username),
		}

		// Hash the password properly
		if err := user.Password.Set("password123"); err != nil {
			return nil, fmt.Errorf("failed to hash password for user %s: %w", username, err)
		}

		if err := s.UsersRepo.Create(ctx, &user); err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", username, err)
		}
		userIDs = append(userIDs, user.ID)
	}

	return userIDs, nil
}

func seedPosts(ctx context.Context, s store.Storage, userIDs []int64) ([]int64, error) {
	postIDs := make([]int64, 0, len(titles))

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
			return nil, fmt.Errorf("failed to create post %s: %w", title, err)
		}
		postIDs = append(postIDs, post.ID)
	}

	return postIDs, nil
}

func seedComments(ctx context.Context, tx *sql.Tx, postIDs, userIDs []int64) error {
	for _, postID := range postIDs {
		numComments := rand.Intn(4) + 2
		for i := 0; i < numComments; i++ {
			userID := userIDs[rand.Intn(len(userIDs))]
			commentText := comments[rand.Intn(len(comments))]

			query := `
				INSERT INTO comments (post_id, user_id, content, created_at)
				VALUES ($1, $2, $3, NOW())
			`
			if _, err := tx.ExecContext(ctx, query, postID, userID, commentText); err != nil {
				return fmt.Errorf("failed to create comment for post %d: %w", postID, err)
			}
		}
	}
	return nil
}

func seedFollowers(ctx context.Context, tx *sql.Tx, userIDs []int64) error {
	// Each user follows 3-8 random other users
	for _, followerID := range userIDs {
		numFollows := rand.Intn(6) + 3 // 3 to 8 follows
		followed := make(map[int64]bool)

		for i := 0; i < numFollows; i++ {
			// Pick a random user to follow (not themselves)
			userID := userIDs[rand.Intn(len(userIDs))]
			for userID == followerID || followed[userID] {
				userID = userIDs[rand.Intn(len(userIDs))]
			}
			followed[userID] = true

			query := `
				INSERT INTO followers (follower_id, user_id, created_at)
				VALUES ($1, $2, NOW())
				ON CONFLICT (follower_id, user_id) DO NOTHING
			`
			if _, err := tx.ExecContext(ctx, query, followerID, userID); err != nil {
				return fmt.Errorf("failed to create follower relationship: %w", err)
			}
		}
	}
	return nil
}
