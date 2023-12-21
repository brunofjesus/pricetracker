package product

import (
	"database/sql"
	"errors"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
)

type FinderFilters product_repository.ProductMetricsFilter

type Finder struct {
	DB                    *sql.DB
	MetricsRepository     *product_repository.MetricsRepository
	ProductRepository     *product_repository.Repository
	ProductMetaRepository *product_repository.MetaRepository
}

func (s *Finder) FindProductById(productId int64) (*product_repository.ProductWithMetrics, error) {
	return s.MetricsRepository.FindProductById(productId, nil)
}

func (s *Finder) FindDetailedProducts(
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

	items, err := s.MetricsRepository.FindProducts(
		paginatedQuery.Offset(), paginatedQuery.Limit(),
		sortField, paginatedQuery.SortDirection,
		(*product_repository.ProductMetricsFilter)(&filters),
		nil,
	)
	if err != nil {
		return nil, err
	}

	count, err := s.MetricsRepository.CountProducts((*product_repository.ProductMetricsFilter)(&filters), nil)
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

func (s *Finder) FindProductByUrl(url string) (*product_repository.Product, error) {
	return s.ProductRepository.FindProductByUrl(url, nil)
}

func (s *Finder) FindProductIdByStoreSlugAndSKUs(storeSlug string, skus []string) (int64, error) {
	productId, err := s.ProductMetaRepository.FindProductIdByStoreSlugAndSKUs(storeSlug, skus, nil)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return productId, nil
}

func (s *Finder) FindProductIdByStoreSlugAndEANs(storeSlug string, eans []string) (int64, error) {
	validEANs := filterNumbers(eans)

	if len(validEANs) > 0 {
		productId, err := s.ProductMetaRepository.FindProductIdByStoreSlugAndEANs(storeSlug, validEANs, nil)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else if err == nil {
			return productId, nil
		}
	}

	return 0, nil
}
