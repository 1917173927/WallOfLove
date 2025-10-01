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
	likeCnt := redis.GetPostLikeCount(postID, 0) // 帖子点赞
	viewCnt := redis.GetPostViewCount(postID)    // 帖子浏览
	return uint64(likeCnt*3 + viewCnt*2)
}

// 计算全部帖子热度并落库
func RefreshAllHeat() error {
	// 拿全部帖子 ID
	var postIDs []uint64
	if err := database.DB.Model(&models.Post{}).Pluck("id", &postIDs).Error; err != nil {
		return err
	}

	// 内存里计算，批量插入
	heatList := make([]models.Heat, 0, len(postIDs))
	for _, pid := range postIDs {
		val := CalculatePostHeat(pid)
		heatList = append(heatList, models.Heat{
			PostID:    pid,
			HeatValue: val,
			UpdatedAt: time.Now(),
		})
	}
	// 批量插入（存在则更新热度，不存在则插入）
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "post_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"heat_value", "updated_at"}),
	}).CreateInBatches(heatList, 1000).Error
}

//获取排行榜
type PopRanking struct {
	PostID    uint64    `json:"post_id"`
	HeatValue uint64    `json:"heat_value"`
	UserNickname string `json:"user_nickname"`
}
func GetPopRanking(userID uint64) ([]PopRanking, error) {
	// 先刷新热度（保证实时）
	if err := RefreshAllHeat(); err != nil {
		return nil, err
	}

	// 拿用户黑名单
	blackList, _ := utils.GetBlackListIDs(userID)

	// 连表查：热度前 10 + 发布人信息 + 匿名打码
	var rank []PopRanking
	err := database.DB.
		Table("heats").
		Select(`heats.post_id, heats.heat_value, 
		        CASE WHEN posts.anonymous THEN '?' ELSE users.nickname END AS user_nickname`).// 根据bool值匿名打码
		Joins("JOIN posts ON posts.id = heats.post_id").
		Joins("JOIN users ON users.id = posts.user_id").
		Where("posts.visibility = ?", true).                      // 公开帖
		Where("posts.user_id NOT IN (?)", blackList).             // 过滤拉黑
		Order("heats.heat_value DESC").                           // 热度降序
		Limit(10).                                                // 只要前 10
		Scan(&rank).Error
	if err != nil {
		return nil, err
	}

	return rank, nil
}