package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
)

// 创建帖子
func CreatePost(post *models.Post) error {
	// 如果未设置发布时间，则立即发布
	if post.ScheduledAt == nil {
		post.IsPublished = true
	}
	return database.DB.Create(post).Error
}
func GetPostDataByID(postID uint64) (*models.Post, error) {
	var post models.Post
	result := database.DB.
		Where("id = ?", postID).
		First(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}
func UpdatePost(post *models.Post) error {
	return database.DB.Model(post).
		Select("content", "anonymous", "visibility").
		Updates(post).Error
}

func DeletePost(postID uint64) error {
	return database.DB.Delete(&models.Post{}, "id = ?", postID).Error
}

//获取未被拉黑的其他人发布的表白
func GetVisiblePosts(userID uint64, page, pageSize int) ([]models.Post, int64,error) {
	sub, _ := utils.GetBlackListIDs(userID)
	var total int64
	database.DB.Model(&models.Post{}).
		Where("visibility = ?", true).
		Where("user_id NOT IN (?)", sub).
		Count(&total)

	var list []models.Post
	err := database.DB.Table("posts").
		Select(`posts.*,(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id AND likes.review_id = 0) AS like_count,
		                 EXISTS(SELECT 1 FROM likes WHERE likes.post_id = posts.id AND likes.review_id = 0 AND likes.user_id = ?) AS liked_by_me`,userID).//查询帖子的点赞数以及该用户是否点赞过该帖子
		Preload("Images").//图片预加载
		Where("visibility = ?", true).
		Where("user_id NOT IN (?)", sub).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Scan(&list).Error
	return list, total, err
}

