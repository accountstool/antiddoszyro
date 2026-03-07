package util

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/domain"
)

func PaginationFromRequest(c *gin.Context, defaultPageSize int) (limit int, offset int, meta domain.Pagination) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(defaultPageSize)))
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize, (page - 1) * pageSize, domain.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func WithPagination(meta domain.Pagination, total int64) domain.Pagination {
	meta.TotalItems = total
	meta.TotalPages = int64(math.Ceil(float64(total) / float64(meta.PageSize)))
	if meta.TotalPages == 0 {
		meta.TotalPages = 1
	}
	return meta
}
