package product

import (
	"database/sql"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	"github.com/brunofjesus/pricetracker/catalog/util/nulltype"
	"github.com/brunofjesus/pricetracker/catalog/util/pagination"
)

var metricsFinderOnce sync.Once
var metricsFinderInstance *metricsFinder

type ProductMetricsSearch struct {
	StoreId    int
	MinPrice   float64
	MaxPrice   float64
	NameLike   string
	BrandLike  string
	Available  nulltype.NullBoolean
	ProductUrl string

	MinDifference      float64
	MaxDifference      float64
	MinDiscountPercent float64
	MaxDiscountPercent float64
	MinAveragePrice    float64
	MaxAveragePrice    float64
	MinMaximumPrice    float64
	MaxMaximumPrice    float64
	MinMinimumPrice    float64
	MaxMinimumPrice    float64

	Page pagination.PaginatedQuery
}

type MetricsFinder interface {
	FindProductById(productId int64) (*product_repository.ProductWithMetrics, error)
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

func (s *metricsFinder) FindProducts(search ProductMetricsSearch) {

}
