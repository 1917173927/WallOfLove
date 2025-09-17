package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
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
func Paginate(data []interface{}, page, pageSize, total int) (result []interface{}, pagination Page) {
	// 偏移量
	offset := (page - 1) * pageSize
	if offset >= len(data) {
		return []interface{}{}, Page{Page: page, PageSize: pageSize, Total: total}
	}
	// 结束位置
	end := offset + pageSize
	if end > len(data) {
		end = len(data)
	}
	// 分页数据
	result = data[offset:end]
	// 分页信息
	pagination = Page{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	return result, pagination
}