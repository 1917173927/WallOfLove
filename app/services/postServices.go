package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
)

type PostService struct{}

// CreatePost 创建帖子，可以独立发布（无需关联图片）
func (s *PostService) CreatePost(post *models.Post) error {
	return database.DB.Create(post).Error
}

func (s *PostService) UpdatePost(post *models.Post) error {
	return database.DB.Save(post).Error
}

func (s *PostService) DeletePost(postID string) error {
	return database.DB.Delete(&models.Post{}, "id = ?", postID).Error
}

// GetVisiblePosts 获取未被拉黑的其他人发布的表白
func (s *PostService) GetVisiblePosts(userID uint64) ([]models.Post, error) {
	var posts []models.Post
	err := database.DB.Where("user_id != ? AND id NOT IN (SELECT blocked_id FROM blacklists WHERE user_id = ?)", userID, userID).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}