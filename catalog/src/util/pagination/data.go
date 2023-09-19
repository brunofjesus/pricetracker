package pagination

import (
	"math"
)

type PaginatedData[R any] struct {
	CurrentPage   int    `json:"current_page"`
	TotalPages    int64  `json:"total_pages"`
	PageSize      int    `json:"page_size"`
	ItemCount     int    `json:"item_count"`
	TotalResults  int64  `json:"total_results"`
	SortField     string `json:"sort_field"`
	SortDirection string `json:"sort_direction"`
	Items         R      `json:"items"`
}

func NewPaginatedData[R any](items R, itemLen int, currentPage int, pageSize int, totalResults int64, sortField string, sortDirection string) *PaginatedData[R] {
	return &PaginatedData[R]{
		CurrentPage:   currentPage,
		TotalPages:    calculateTotalPages(totalResults, pageSize),
		PageSize:      pageSize,
		ItemCount:     itemLen,
		TotalResults:  totalResults,
		SortField:     sortField,
		SortDirection: sortDirection,
		Items:         items,
	}
}

func calculateTotalPages(totalResults int64, pageSize int) int64 {
	pages := float64(totalResults) / float64(pageSize)
	return int64(math.Ceil(pages))
}
