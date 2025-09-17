package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"gorm.io/gorm"
)

// CreatePost 创建帖子，可以独立发布（无需关联图片）
func CreatePost(post *models.Post) error {
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
func UpdatePost(post *models.Post, oldVersion uint) error {
	tx := database.DB.Model(&models.Post{}).
		Where("id = ? AND version = ?", post.ID, oldVersion).
		Updates(map[string]any{
			"content":         post.Content,
			"anonymous":       post.Anonymous,
			"visibility":      post.Visibility,
			"version":         gorm.Expr("version + 1"),
		})
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return tx.Error
}

func DeletePost(postID uint64) error {
	return database.DB.Delete(&models.Post{}, "id = ?", postID).Error
}

// GetVisiblePosts 获取未被拉黑的其他人发布的表白
func GetVisiblePosts(userID uint64, page, pageSize int) ([]models.Post, int64, error) {
	sub := database.DB.Model(&models.Blacklist{}).
		Where("user_id = ?", userID).
		Select("blocked_id")

	var total int64
	database.DB.Model(&models.Post{}).
		Where("visibility = ?", true).
		Where("user_id NOT IN (?)", sub).
		Count(&total)

	var list []models.Post
	err := database.DB.Preload("Images").
		Where("visibility = ?", true).
		Where("user_id NOT IN (?)", sub).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&list).Error

	return list, total, err
}