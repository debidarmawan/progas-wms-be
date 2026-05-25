package helper

import "progas-wms-be/dto"

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

func NormalizePagination(query *dto.ListQuery) (page, limit, offset int) {
	page = query.Page
	if page < 1 {
		page = DefaultPage
	}

	limit = query.Limit
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset = (page - 1) * limit
	return page, limit, offset
}

func BuildPaginationMeta(page, limit int, total int64) dto.PaginationMeta {
	totalPages := 0
	if total > 0 {
		totalPages = int(total) / limit
		if int(total)%limit > 0 {
			totalPages++
		}
	}

	return dto.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}
