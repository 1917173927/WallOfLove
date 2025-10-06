package services

import (
	"time"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/robfig/cron/v3"
)

func StartScheduler() {
	c := cron.New()
	_, _ = c.AddFunc("@every 0.5m", checkAndPublishScheduledPosts)
	c.Start()
}

func checkAndPublishScheduledPosts() {
	var posts []models.Post
	now := time.Now()
	database.DB.Where("is_published = ? AND scheduled_at <= ?", false, now).Find(&posts)
	for _, post := range posts {
		post.IsPublished = true
		database.DB.Save(&post)
	}
}
