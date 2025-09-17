package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/app/utils"
)

//create review2
func CreateReview2(review2 *models.Review2) error {
	return database.DB.Create(review2).Error
}

func GetReview2sByPostID(postID uint64) error {
	var reviews []models.Review2
	return database.DB.Where("post_id = ?", postID).Find(&reviews).Error
}
// Get Review2s
func GetVisibleReview2s(reviewID uint64, page, pageSize int) ([]models.Review2, int64, error) {
	sub := database.DB.Model(&models.Blacklist{}).
		Where("user_id = ?", reviewID).
		Select("blocked_id")

	var total int64
	database.DB.Model(&models.Review2{}).
		Where("user_id NOT IN (?)", sub).
		Count(&total)

	var list []models.Review2
	err := database.DB.
		Where("user_id NOT IN (?)", sub).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&list).Error

	return list, total, err
}