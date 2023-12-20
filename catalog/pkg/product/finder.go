package product

import (
	"database/sql"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
)

type FinderFilters product_repository.ProductMetricsFilter

type Finder struct {
	DB         *sql.DB
	Repository *product_repository.MetricsRepository
}

func (s *Finder) FindProductById(productId int64) (*product_repository.ProductWithMetrics, error) {
	return s.Repository.FindProductById(productId, nil)
}

func (s *Finder) FindProducts(
	paginatedQuery pagination.PaginatedQuery,
	filters FinderFilters,
) (*pagination.PaginatedData[[]product_repository.ProductWithMetrics], error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	sortField := paginatedQuery.GetSortFieldIfValid(
		[]string{
			"product_id", "store_id", "name", "brand", "price", "available",
			"diff", "discount_percent", "average", "minimum", "maximum",
		},
		"name",
	)

	items, err := s.Repository.FindProducts(
		paginatedQuery.Offset(), paginatedQuery.Limit(),
		sortField, paginatedQuery.SortDirection,
		(*product_repository.ProductMetricsFilter)(&filters),
		nil,
	)
	if err != nil {
		return nil, err
	}

	count, err := s.Repository.CountProducts((*product_repository.ProductMetricsFilter)(&filters), nil)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return pagination.NewPaginatedData[[]product_repository.ProductWithMetrics](
		items, len(items),
		paginatedQuery.Page, paginatedQuery.PageSize, count,
		paginatedQuery.SortField, paginatedQuery.SortDirection,
	), nil
}
