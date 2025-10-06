package services

import (
	"time"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
	"gorm.io/gorm/clause"
)

func CalculatePostHeat(postID uint64) uint64 {
	likeCnt := redis.GetPostLikeCount(postID, 0)
	viewCnt := redis.GetPostViewCount(postID)
	return uint64(likeCnt*3 + viewCnt*2)
}

func RefreshAllHeat() error {
	var postIDs []uint64
	if err := database.DB.Model(&models.Post{}).Pluck("id", &postIDs).Error; err != nil {
		return err
	}
	heatList := make([]models.Heat, 0, len(postIDs))
	for _, pid := range postIDs {
		val := CalculatePostHeat(pid)
		heatList = append(heatList, models.Heat{
			PostID:    pid,
			HeatValue: val,
			UpdatedAt: time.Now(),
		})
	}
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "post_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"heat_value", "updated_at"}),
	}).CreateInBatches(heatList, 1000).Error
}

type PopRanking struct {
	PostID       uint64 `json:"post_id"`
	HeatValue    uint64 `json:"heat_value"`
	UserNickname string `json:"user_nickname"`
}

func GetPopRanking(userID uint64) ([]PopRanking, error) {
	if err := RefreshAllHeat(); err != nil {
		return nil, err
	}
	blackList, _ := utils.GetBlackListIDs(userID)
	var rank []PopRanking
	base := database.DB.
		Table("heats").
		Select(`heats.post_id, heats.heat_value, CASE WHEN posts.anonymous THEN '?' ELSE users.nickname END AS user_nickname`).
		Joins("JOIN posts ON posts.id = heats.post_id").
		Joins("JOIN users ON users.id = posts.user_id").
		Where("posts.visibility = ?", true)
	if len(blackList) > 0 {
		base = base.Where("posts.user_id NOT IN (?)", blackList)
	}
	err := base.Order("heats.heat_value DESC").Limit(10).Scan(&rank).Error
	if err != nil {
		return nil, err
	}
	return rank, nil
}
