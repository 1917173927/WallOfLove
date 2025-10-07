package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/config"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/gin-gonic/gin"
)

func toURLPath(p string) string {
	s := filepath.ToSlash(p)
	s = strings.TrimPrefix(s, "./")
	return strings.TrimPrefix(s, "/")
}

// 尺寸校验 -> 去重(hash) -> 类型识别(magic number) -> 落盘 -> 事务处理头像唯一 -> 返回记录
func UploadImage(c *gin.Context, userID uint64, username, postID, isAvatar string, file *multipart.FileHeader) (*models.Image, error) {
	maxSize := config.Config.GetInt64("image.max_size")
	if maxSize <= 0 {
		maxSize = 2 * 1024 * 1024 // 默认 2MB
	}
	if file.Size > maxSize {
		return nil, apiException.ImageSizeExceeded
	}

	// 计算哈希
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, fileContent); err != nil {
		return nil, err
	}
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// 检查物理文件是否已存在（数据库中是否有相同哈希的图片）
	var existingImage models.Image
	fileAlreadyExists := false
	if err := database.DB.Where("checksum = ?", fileHash).First(&existingImage).Error; err == nil {
		fileAlreadyExists = true
	}

	// 重置指针
	if _, err := fileContent.Seek(0, 0); err != nil {
		return nil, err
	}

	// Magicnumber检测文件类型
	buffer := make([]byte, 512)
	n, err := fileContent.Read(buffer)
	if err != nil {
		return nil, err
	}
	buffer = buffer[:n]
	mtype := mimetype.Detect(buffer)
	if !strings.HasPrefix(mtype.String(), "image/") {
		return nil, apiException.NotImage
	}
	ext := strings.ToLower(mtype.Extension())

	userFolder := fmt.Sprintf("%d-%s", userID, username)
	userFolderPath := filepath.Join("images", userFolder)
	if err := os.MkdirAll(userFolderPath, os.ModePerm); err != nil {
		return nil, err
	}

	fileName := fileHash + ext
	dst := filepath.Join(userFolderPath, fileName)

	// 只有物理文件不存在时才保存
	if !fileAlreadyExists {
		if err := c.SaveUploadedFile(file, dst); err != nil {
			return nil, err
		}
	} else {
		// 若已存在，确保路径存在（如不同用户目录下）
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			if err := c.SaveUploadedFile(file, dst); err != nil {
				return nil, err
			}
		}
	}

	// 解析 postID（可选）
	var postIDUint *uint64
	if postID != "" {
		postIDParsed, err := strconv.ParseUint(postID, 10, 64)
		if err == nil {
			postIDUint = &postIDParsed
		}
	}

	isAvatarLower := strings.ToLower(strings.TrimSpace(isAvatar))
	avatarFlag := isAvatarLower == "1" || isAvatarLower == "true" || isAvatarLower == "yes"

	image := &models.Image{
		UserID:    userID,
		PostID:    postIDUint,
		IsAvatar:  avatarFlag,
		FilePath:  dst,
		Size:      file.Size,
		Checksum:  fileHash,
		CreatedAt: time.Now(),
	}

	// 头像：事务内保证唯一
	if avatarFlag {
		tx := database.DB.Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}
		if err := tx.Model(&models.Image{}).Where("user_id = ? AND is_avatar = ?", userID, true).Update("is_avatar", false).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Create(image).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("avatar_path", dst).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}
	} else {
		if err := database.DB.Create(image).Error; err != nil {
			return nil, err
		}
	}

	// 相对路径
	image.FilePath = toURLPath(image.FilePath)
	return image, nil
}

func DeleteImage(imageID uint64, userID uint64) error {
	var image models.Image
	if err := database.DB.Where("id = ?", imageID).First(&image).Error; err != nil {
		return apiException.ImageNotFound
	}
	if image.UserID != userID {
		return apiException.NotPermission
	}
	if err := os.Remove(image.FilePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := database.DB.Delete(&image).Error; err != nil {
		return err
	}
	return nil
}
