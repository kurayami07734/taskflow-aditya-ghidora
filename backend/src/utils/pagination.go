package utils

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func GetPagination(r *http.Request) (page, limit, offset int) {
	page = 1
	limit = 10

	if params := r.URL.Query(); len(params) > 0 {
		if p := params["page"]; len(p) > 0 {
			if parsed, err := strconv.Atoi(p[0]); err == nil && parsed > 0 {
				page = parsed
			}
		}
		if l := params["limit"]; len(l) > 0 {
			if parsed, err := strconv.Atoi(l[0]); err == nil && parsed > 0 {
				limit = parsed
				if limit > 30 {
					limit = 30
				}
			}
		}
	}

	offset = (page - 1) * limit
	return
}

func CalculatePagination(page, limit, total int) Pagination {
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
