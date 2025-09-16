package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"gorm.io/gorm"
	"time"
)

type PostService struct {}

func (ps *PostService) CreatePost(post *models.Post) error {
	if post.Anonymous {
		if post.UserID == 0 {
			return gorm.ErrInvalidData
		}
		post.UserID = 0
	}

	if post.Images == nil || len(post.Images) > 9 {
		return gorm.ErrInvalidData
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	result := database.DB.Create(post)
	return result.Error
}

func (ps *PostService) UpdatePost(post *models.Post) error {
	post.UpdatedAt = time.Now()
	result := database.DB.Model(&models.Post{}).
		Where("id = ?", post.ID).
		Updates(post)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (ps *PostService) DeletePost(postID string) error {
	result := database.DB.Where("id = ?", postID).Delete(&models.Post{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (ps *PostService) CheckPostOwnership(postID string, userID uint) bool {
	var post models.Post
	result := database.DB.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return false
	}
	return uint(post.UserID) == userID
}