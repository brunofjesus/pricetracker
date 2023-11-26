package product

import (
	"database/sql"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	"github.com/brunofjesus/pricetracker/catalog/util/pagination"
)

var metricsFinderOnce sync.Once
var metricsFinderInstance *metricsFinder

type ProductFinderFilters product_repository.ProductMetricsFilter

type MetricsFinder interface {
	FindProductById(productId int64) (*product_repository.ProductWithMetrics, error)
	FindProducts(
		pagination pagination.PaginatedQuery,
		filters ProductFinderFilters,
	) (*pagination.PaginatedData[[]product_repository.ProductWithMetrics], error)
}

type metricsFinder struct {
	db                       *sql.DB
	productMetricsRepository product_repository.ProductMetricsRepository
}

func GetMetricsFinder() *metricsFinder {
	metricsFinderOnce.Do(func() {
		metricsFinderInstance = &metricsFinder{
			db:                       repository.GetDatabaseConnection(),
			productMetricsRepository: product_repository.GetProductMetricsRepository(),
		}
	})
	return metricsFinderInstance
}

func (s *metricsFinder) FindProductById(productId int64) (*product_repository.ProductWithMetrics, error) {
	return s.productMetricsRepository.FindProductById(productId, nil)
}

func (s *metricsFinder) FindProducts(
	paginatedQuery pagination.PaginatedQuery,
	filters ProductFinderFilters,
) (*pagination.PaginatedData[[]product_repository.ProductWithMetrics], error) {

	tx, err := s.db.Begin()

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

	items, err := s.productMetricsRepository.FindProducts(
		paginatedQuery.Offset(), paginatedQuery.Limit(),
		sortField, paginatedQuery.SortDirection,
		(*product_repository.ProductMetricsFilter)(&filters),
		nil,
	)

	if err != nil {
		return nil, err
	}

	count, err := s.productMetricsRepository.CountProducts((*product_repository.ProductMetricsFilter)(&filters), nil)
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
