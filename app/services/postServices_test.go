package services

import (
	"testing"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestDeletePost 测试删除帖子的功能
// 测试场景包括：
// 1. 删除存在的帖子
// 2. 删除不存在的帖子
// 3. 删除帖子后验证数据库记录是否被删除
func TestDeletePost(t *testing.T) {
	// 初始化测试数据库
	db, err := database.InitTestDB()
	assert.NoError(t, err)
	defer database.CleanupTestDB(db)

	// 测试用例1：删除存在的帖子
	t.Run("删除存在的帖子", func(t *testing.T) {
		// 创建测试数据
		post := &models.Post{ID: 1, Content: "Test Post"}
		err := db.Create(post).Error
		assert.NoError(t, err)

		// 执行删除操作
		err = DeletePost(post.ID)
		assert.NoError(t, err)

		// 验证帖子是否被删除
		var deletedPost models.Post
		err = db.First(&deletedPost, "id = ?", post.ID).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// 测试用例2：删除不存在的帖子
	t.Run("删除不存在的帖子", func(t *testing.T) {
		err := DeletePost(999)
		assert.NoError(t, err) // 删除不存在的帖子不应返回错误
	})
}