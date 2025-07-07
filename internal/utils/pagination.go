package utils

import (
	"math"
	"strconv"
	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

type PaginationResult struct {
	Offset   int
	Limit    int
	Page     int
	PageSize int
}

func (p *PaginationParams) Calculate() PaginationResult {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}

	offset := (p.Page - 1) * p.PageSize

	return PaginationResult{
		Offset:   offset,
		Limit:    p.PageSize,
		Page:     p.Page,
		PageSize: p.PageSize,
	}
}

func CalculateHasMore(totalCount int64, page, pageSize int) bool {
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	return int64(page) < totalPages
}

// GetQueryInt lấy giá trị int từ query, có giá trị mặc định nếu không hợp lệ
func GetQueryInt(c *gin.Context, key string, defaultVal int) int {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
