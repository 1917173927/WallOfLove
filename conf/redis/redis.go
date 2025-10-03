package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

func Init() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Try to ping Redis with retries
	var err error
	for i := 0; i < 3; i++ {
		if err = Rdb.Ping(Ctx).Err(); err == nil {
			return nil
		}
		// Wait before retrying
		time.Sleep(time.Second * 2)
	}
	
	return fmt.Errorf("failed to connect to Redis after 3 attempts: %v", err)
}
// 点赞数 +1 并返回最新值
func IncrPostLike(postID, reviewID uint64) int64 {
	count, _ := Rdb.Incr(Ctx, fmt.Sprintf("post:%d:review:%d:likes", postID, reviewID)).Result()
	return count
}

// 点赞数 -1 并返回最新值
func DecrPostLike(postID, reviewID uint64) int64 {
	count, _ := Rdb.Decr(Ctx, fmt.Sprintf("post:%d:review:%d:likes", postID, reviewID)).Result()
	return count
}

// 标记用户已赞
func SetUserLiked(postID, userID, reviewID uint64) {
	Rdb.Set(Ctx, fmt.Sprintf("post:%d:review:%d:liked:uid:%d", postID, reviewID, userID), 1, 0)
}

// 取消已赞标记
func DelUserLiked(postID, userID, reviewID uint64) {
	Rdb.Del(Ctx, fmt.Sprintf("post:%d:review:%d:liked:uid:%d", postID, reviewID, userID))
}

// 读点赞数
func GetPostLikeCount(postID,reviewID uint64) int64 {
	key := fmt.Sprintf("post:%d:review:%d:likes", postID,reviewID)
	//读缓存
	if val, err := Rdb.Get(Ctx, key).Int64(); err == nil {
		return val
	}
	//错误：读库
	var cnt int64
	database.DB.Model(&models.Like{}).Where("post_id = ? AND review_id = ?", postID,reviewID).Count(&cnt)
	//写回缓存
	Rdb.Set(Ctx, key, cnt, 0)
	return cnt
}

// 读是否已赞
func IsUserLiked(postID, userID,reviewID uint64) bool {
	key := fmt.Sprintf("post:%d:review:%d:liked:uid:%d", postID, userID,reviewID)
	//读缓存
	if val, err := Rdb.Get(Ctx, key).Result(); err == nil {
		return val == "1"
	}
	//错误：读库
	var like int64
	database.DB.Model(&models.Like{}).
		Where("user_id = ? AND post_id = ? AND review_id = ?", userID,reviewID,postID).
		Count(&like)
	//写回缓存
	if like > 0 {
		Rdb.Set(Ctx, key, 1, 0)
		return true
	}
	Rdb.Set(Ctx, key, 0, 0)
	return false
}

// 评论数+1 并返回最新值
func IncrReview(postID uint64) int64 {
	count, _ := Rdb.Incr(Ctx, fmt.Sprintf("post:%d:reviews", postID)).Result()
	return count
}

// 评论数-1 并返回最新值
func DecrReview(postID uint64) int64 {
	count, _ := Rdb.Incr(Ctx, fmt.Sprintf("post:%d:reviews", postID)).Result()
	return count
}

// 读评论数
func GetPostReviewCount(postID uint64) int64 {
	key := fmt.Sprintf("post:%d:reviews", postID)
	//读缓存
	if val, err := Rdb.Get(Ctx, key).Int64(); err == nil {
		return val
	}
	//错误：读库
	var cnt int64
	database.DB.Model(&models.Review{}).Where("post_id = ? ", postID).Count(&cnt)
	//写回缓存
	Rdb.Set(Ctx, key, cnt, 0)
	return cnt
}

// 回复数+1 并返回最新值
func IncrReply(reviewID uint64) int64 {
	count, _ := Rdb.Incr(Ctx, fmt.Sprintf("review:%d:replies", reviewID)).Result()
	return count
}

// 回复数-1 并返回最新值
func DecrReply(reviewID uint64) int64 {
	count, _ := Rdb.Decr(Ctx, fmt.Sprintf("review:%d:replies", reviewID)).Result()
	return count
}

// 读回复数
func GetReviewReplyCount(reviewID uint64) int64 {
	key := fmt.Sprintf("review:%d:replies", reviewID)
	//读缓存
	if val, err := Rdb.Get(Ctx, key).Int64(); err == nil {
		return val
	}
	//错误：读库
	var cnt int64
	database.DB.Model(&models.Reply{}).Where("review_id = ? ", reviewID).Count(&cnt)
	//写回缓存
	Rdb.Set(Ctx, key, cnt, 0)
	return cnt
}

//浏览数+1 并返回最新值
func IncrView(postID uint64) int64 {
	count, _ := Rdb.Incr(Ctx, fmt.Sprintf("post:%d:views", postID)).Result()
	return count
}
// 读浏览数
func GetPostViewCount(postID uint64) int64 {
	key := fmt.Sprintf("post:%d:views", postID)
	//读缓存
	if val, err := Rdb.Get(Ctx, key).Int64(); err == nil {
		return val
	}
	//错误：读库
	var cnt int64
	database.DB.Model(&models.View{}).Where("post_id = ? ", postID).Count(&cnt)
	//写回缓存
	Rdb.Set(Ctx, key, cnt, 0)
	return cnt
}