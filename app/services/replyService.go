package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
)

//创建回复
func CreateReply(reply *models.Reply) error {
	return database.DB.Create(reply).Error
}