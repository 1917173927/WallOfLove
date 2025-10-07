package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
)

// 创建回复
func CreateReply(reply *models.Reply) error {
	redis.IncrReply(reply.ReviewID)
	return database.DB.Create(reply).Error
}

func GetRepliesByReviewID(reviewID uint64, userID uint64, page int, pageSize int) ([]models.ReplyWithNickname, int64, error) {
	sub, _ := utils.GetBlackListIDs(userID)
	var total int64

	countDB := database.DB.Model(&models.Reply{}).
		Where("review_id = ?", reviewID)
	if len(sub) > 0 {
		countDB = countDB.Where("user_id NOT IN (?)", sub)
	}
	countDB.Count(&total)

	// 获取原始回复列表
	var replies []models.Reply
	q := database.DB.
		Where("review_id = ?", reviewID)
	if len(sub) > 0 {
		q = q.Where("user_id NOT IN (?)", sub)
	}
	err := q.
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&replies).Error
	if err != nil {
		return nil, 0, err
	}

	// 收集所有回复的用户ID
	userIDs := make([]uint64, 0, len(replies))
	for _, r := range replies {
		userIDs = append(userIDs, r.UserID)
	}

	// 批量获取用户信息
	userInfos := make(map[uint64]struct {
		Nickname   string
		AvatarPath string
	})
	for _, id := range userIDs {
		user, err := GetUserDataByID(id)
		if err == nil && user != nil {
			userInfos[id] = struct {
				Nickname   string
				AvatarPath string
			}{
				Nickname:   user.Nickname,
				AvatarPath: user.AvatarPath,
			}
		}
	}

	// 创建返回结果
	list := make([]models.ReplyWithNickname, 0, len(replies))
	for _, r := range replies {
		list = append(list, models.ReplyWithNickname{
			Reply:      r,
			Nickname:   userInfos[r.UserID].Nickname,
			AvatarPath: userInfos[r.UserID].AvatarPath,
		})
	}

	return list, total, nil
}

// 删除回复
func DeleteReply(replyID uint64) error {
	redis.DecrReply(replyID)
	return database.DB.Delete(&models.Reply{}, replyID).Error
}

func GetReplyByReplyID(replyID uint64) (*models.Reply, error) {
	var reply models.Reply
	result := database.DB.
		Where("id = ?", replyID).
		First(&reply)
	if result.Error != nil {
		return nil, result.Error
	}
	return &reply, nil
}
