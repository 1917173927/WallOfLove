package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
)

//Create review
func CreateReview(review *models.Review) error {
	return database.DB.Create(review).Error
}

// Get Reviews
func GetReviewsByPostID(postID uint64) ([]models.Review, error) {	
	var reviews []models.Review
	err := database.DB.Where("post_id = ?", postID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
