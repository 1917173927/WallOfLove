package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
)

//创建回复
func CreateReview2(review2 *models.Review2) error {
	return database.DB.Create(review2).Error
}