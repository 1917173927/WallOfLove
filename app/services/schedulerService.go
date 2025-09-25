package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/robfig/cron/v3"
	"time"
)

// StartScheduler 启动定时任务服务
func StartScheduler() {
	c := cron.New()
	// 每分钟检查一次
	_, _ = c.AddFunc("@every 0.5m", checkAndPublishScheduledPosts)
	c.Start()
}

// checkAndPublishScheduledPosts 检查并发布到期的帖子
func checkAndPublishScheduledPosts() {
	var posts []models.Post
	now := time.Now()
	// 查询所有未发布且发布时间已过的帖子
	database.DB.Where("is_published = ? AND scheduled_at <= ?", false, now).Find(&posts)
	for _, post := range posts {
		post.IsPublished = true
		database.DB.Save(&post)
	}
}