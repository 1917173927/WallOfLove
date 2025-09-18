package utils

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
)

func GetBlackListIDs(userID uint64) ([]uint64, error) {
	var ids []uint64
	err := database.DB.Model(&models.Blacklist{}).
		Where("user_id = ?", userID).
		Pluck("blocked_id", &ids).Error
	return ids, err
}