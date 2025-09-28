// Package services 包含所有业务逻辑处理，负责与数据库交互和核心业务逻辑的实现。
package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"golang.org/x/crypto/bcrypt"
)

// CheckUsername 检查用户名是否已存在，执行以下操作：
// 1. 查询数据库中是否存在指定用户名的记录
// 2. 返回查询结果或错误
func CheckUsername(username string) error {
	result := database.DB.Where("username=?", username).First(&models.User{})
	return result.Error
}

// GetUser 根据用户名获取用户信息，执行以下操作：
// 1. 查询数据库中指定用户名的记录
// 2. 返回用户信息或错误

// 2. 返回用户信息或错误
func GetUser(username string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username=?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 注册用户
func Register(user models.User) error {
	result := database.DB.Create(&user)
	return result.Error
}

// 密码加密
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 密码比对
func CompareHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

// 根据ID获取用户信息
func GetUserDataByID(userID uint64) (*models.User, error) {
	var user models.User
	result := database.DB.
		Where("id = ?", userID).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 更新用户信息
func UpdateProfile(user *models.User) error {
	return database.DB.Model(user).
		Select("nickname", "username", "password", "avatar_path").
		Updates(user).Error
}

// 拉黑用户
func BlackUser(userID, blockedID uint64) error {
	return database.DB.Create(&models.Blacklist{UserID: userID, BlockedID: blockedID}).Error
}

// 取消拉黑用户
func UnblackUser(userID, blockedID uint64) error {
	return database.DB.Where("user_id = ? AND blocked_id = ?", userID, blockedID).Delete(&models.Blacklist{}).Error
}

// 获取拉黑用户信息列表
type BlackedUser struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func GetBlackedUsers(userID uint64) ([]BlackedUser, error) {
	// 拿被拉黑人ID 列表
	ids, err := utils.GetBlackListIDs(userID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []BlackedUser{}, nil // 没人被拉黑，直接空列表
	}
	var list []BlackedUser
	err = database.DB.
		Table("Users").
		Select("id as user_id, username, nickname").
		Where("id IN ?", ids).
		Scan(&list).Error

	return list, err
}
