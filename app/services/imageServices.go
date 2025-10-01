// Package services 包含所有业务逻辑处理，负责与数据库交互和核心业务逻辑的实现。
package services

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/config"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/gin-gonic/gin"
)

// UploadImage 处理图片上传的核心逻辑，包括以下步骤：
// 1. 验证文件类型和大小
// 2. 生成唯一的文件名并保存到指定目录
// 3. 将图片信息保存到数据库
// 4. 返回图片的访问路径或错误信息

// UploadImage 处理图片上传逻辑，执行以下操作：
// 1. 验证文件类型和大小
// 2. 计算文件哈希值
// 3. 检查是否已存在相同图片
// 4. 如果不存在则保存新文件
// 5. 记录文件信息到数据库
// 6. 返回文件信息或错误
func UploadImage(c *gin.Context, userID uint64, username string, postID string, isAvatar string, file *multipart.FileHeader) (*models.Image, error) {
	// 检查图片大小
	maxSize := config.Config.GetInt64("image.max_size")
	if maxSize <= 0 {
		maxSize = 2 * 1024 * 1024 // 默认 2MB
	}
	if file.Size > maxSize {
		return nil, apiException.ImageSizeExceeded
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" {
		return nil, apiException.ImageTypeInvalid
	}

	// 读取文件内容并计算哈希
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	// 计算文件哈希
	hasher := sha256.New()
	if _, err := io.Copy(hasher, fileContent); err != nil {
		return nil, err
	}
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// 检查是否已存在相同图片
	var existingImage models.Image
	if err := database.DB.Where("checksum = ?", fileHash).First(&existingImage).Error; err == nil {
		// 图片已存在，直接返回已有记录
		return &existingImage, nil
	}

	// 重置文件指针
	if _, err := fileContent.Seek(0, 0); err != nil {
		return nil, err
	}

	// 验证文件内容类型
	buffer := make([]byte, 512)
	if _, err := fileContent.Read(buffer); err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, apiException.NotImage
	}

	// 创建用户专属文件夹
	userFolder := fmt.Sprintf("%d-%s", userID, username)
	userFolderPath := filepath.Join("images", userFolder)
	if err := os.MkdirAll(userFolderPath, os.ModePerm); err != nil {
		return nil, err
	}

	// 生成唯一文件名
	fileName := fileHash + ext
	dst := filepath.Join(userFolderPath, fileName)

	// 保存图片到服务器
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return nil, err
	}

	// 保存图片信息到数据库
	var postIDUint *uint64
	if postID != "" {
		postIDParsed, err := strconv.ParseUint(postID, 10, 64)
		if err == nil {
			postIDUint = &postIDParsed
		}
	}
	image := &models.Image{
		UserID:    userID,
		PostID:    postIDUint,
		FilePath:  dst,
		Size:      file.Size,
		Checksum:  fileHash,
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(image).Error; err != nil {
		return nil, err
	}
	return image, nil
}

// DeleteImage 删除图片及其数据库记录
// 1. 根据ID获取图片信息
// 2. 验证操作权限
// 3. 删除物理文件
// 4. 删除数据库记录
func DeleteImage(imageID uint64, userID uint64) error {
	// 获取图片信息
	var image models.Image
	if err := database.DB.Where("id = ?", imageID).First(&image).Error; err != nil {
		return apiException.ImageNotFound
	}

	// 验证操作权限
	if image.UserID != userID {
		return apiException.NotPermission
	}

	// 删除物理文件
	if err := os.Remove(image.FilePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 删除数据库记录
	if err := database.DB.Delete(&image).Error; err != nil {
		return err
	}

	return nil
}

