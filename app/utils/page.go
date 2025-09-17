package utils

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type Page struct {
	Page     int `json:"page"`     // 当前页
	PageSize int `json:"pageSize"` // 页大小
	Total    int `json:"total"`    // 数据量
}

// GetPaginationParams 获取分页参数
func GetPaginationParams(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	return page, pageSize
}

// Paginate 分页
func Paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := size
		if pageSize > 20 {
			pageSize = 20
		} else if pageSize <= 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}