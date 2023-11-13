package pagination

import (
	"errors"
	"net/http"
	"strconv"
)

type PaginatedQuery struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"page_size"`
	SortField     string `json:"sort_field"`
	SortDirection string `json:"sort_direction"`
}

func NewPaginatedQuery(page, pageSize int, sortField, sortDirection string) (*PaginatedQuery, error) {
	if page < 1 {
		return nil, errors.New("page needs to be 1 or greater")
	}
	if pageSize > 100 {
		return nil, errors.New("page size cannot be that large")
	}

	if sortDirection == "" {
		sortDirection = "asc"
	} else if sortDirection != "asc" && sortDirection != "desc" {
		return nil, errors.New("sort direction can be either 'asc' or 'desc'")
	}

	return &PaginatedQuery{
		Page:          page,
		PageSize:      pageSize,
		SortField:     sortField,
		SortDirection: sortDirection,
	}, nil
}

func FromHttpRequest(r *http.Request) (*PaginatedQuery, error) {
	query := r.URL.Query()

	var err error

	page := 1
	size := 10

	if query.Get("page") != "" {
		page, err = strconv.Atoi(query.Get("page"))
		if err != nil {
			return nil, err
		}
	}

	if query.Get("size") != "" {
		size, err = strconv.Atoi(query.Get("size"))
		if err != nil {
			return nil, err
		}
	}

	sortField := query.Get("sortField")
	sortDirection := query.Get("sortDirection")

	return NewPaginatedQuery(page, size, sortField, sortDirection)
}

func (p *PaginatedQuery) Offset() int64 {
	return int64((p.Page - 1) * p.PageSize)
}

func (p *PaginatedQuery) Limit() int {
	return p.PageSize
}

func (p *PaginatedQuery) HasSortField() bool {
	return p.SortField != ""
}

func (p *PaginatedQuery) HasValidSortField(validSortFields []string) bool {
	for _, field := range validSortFields {
		if p.SortField == field {
			return true
		}
	}

	return false
}

func (p *PaginatedQuery) GetSortFieldIfValid(validSortFields []string, fallback string) string {
	if p.HasSortField() && p.HasValidSortField(validSortFields) {
		return p.SortField
	}

	return fallback
}
