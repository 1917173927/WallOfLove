package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/app/utils"
)

//创建回复
func CreateReply(reply *models.Reply) error {
	return database.DB.Create(reply).Error
}

func GetRepliesByReviewID(reviewID uint64,userID uint64,page int,pageSize int) ([]models.Reply,int64, error) {
	sub, _ := utils.GetBlackListIDs(userID)
	var total int64

	database.DB.Model(&models.Reply{}).
		Where("review_id = ?", reviewID).
		Where("user_id NOT IN (?)", sub).
		Count(&total)

	var list []models.Reply
	err := database.DB.
		Where("review_id = ?", reviewID).
		Where("user_id NOT IN (?)", sub).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&list).Error
	return list, total, err
}