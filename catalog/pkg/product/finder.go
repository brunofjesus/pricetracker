package product

import (
	"database/sql"
	"errors"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
	"github.com/brunofjesus/pricetracker/catalog/util/nulltype"
)

type FinderFilters struct {
	ProductId []int64

	StoreId    *int
	MinPrice   *float64
	MaxPrice   *float64
	NameLike   string
	BrandLike  string
	Available  *bool
	ProductUrl string

	MinDifference      *float64
	MaxDifference      *float64
	MinDiscountPercent *float64
	MaxDiscountPercent *float64
	MinAveragePrice    *float64
	MaxAveragePrice    *float64
	MinMinimumPrice    *float64
	MaxMinimumPrice    *float64
	MinMaximumPrice    *float64
	MaxMaximumPrice    *float64
}

type Finder struct {
	DB                         *sql.DB
	ProductWithStatsRepository *product_repository.ProductWithStatsRepository
	ProductRepository          *product_repository.Repository
	ProductMetaRepository      *product_repository.MetaRepository
}

func (s *Finder) FindProductById(productId int64) (*product_repository.ProductWithStats, error) {
	return s.ProductWithStatsRepository.FindProductById(productId, nil)
}

func (s *Finder) FindDetailedProducts(
	paginatedQuery pagination.PaginatedQuery,
	filters FinderFilters,
) (*pagination.PaginatedData[[]product_repository.ProductWithStats], error) {
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

	repositoryFilters := product_repository.ProductWithStatsFilter{
		ProductId:          filters.ProductId,
		StoreId:            nulltype.IntPtrToSqlNullInt32(filters.StoreId),
		MinPrice:           nulltype.Float64PtrToSqlNullFloat64(filters.MinPrice),
		MaxPrice:           nulltype.Float64PtrToSqlNullFloat64(filters.MaxPrice),
		NameLike:           nulltype.StringToSqlNullString(filters.NameLike),
		BrandLike:          nulltype.StringToSqlNullString(filters.BrandLike),
		Available:          nulltype.BooleanPrtToSqlNullBool(filters.Available),
		ProductUrl:         nulltype.StringToSqlNullString(filters.ProductUrl),
		MinDifference:      nulltype.Float64PtrToSqlNullFloat64(filters.MinDifference),
		MaxDifference:      nulltype.Float64PtrToSqlNullFloat64(filters.MaxDifference),
		MinDiscountPercent: nulltype.Float64PtrToSqlNullFloat64(filters.MinDiscountPercent),
		MaxDiscountPercent: nulltype.Float64PtrToSqlNullFloat64(filters.MaxDiscountPercent),
		MinAveragePrice:    nulltype.Float64PtrToSqlNullFloat64(filters.MinAveragePrice),
		MaxAveragePrice:    nulltype.Float64PtrToSqlNullFloat64(filters.MaxAveragePrice),
		MinMinimumPrice:    nulltype.Float64PtrToSqlNullFloat64(filters.MinMinimumPrice),
		MaxMinimumPrice:    nulltype.Float64PtrToSqlNullFloat64(filters.MaxMinimumPrice),
		MinMaximumPrice:    nulltype.Float64PtrToSqlNullFloat64(filters.MinMaximumPrice),
		MaxMaximumPrice:    nulltype.Float64PtrToSqlNullFloat64(filters.MaxMaximumPrice),
	}

	items, err := s.ProductWithStatsRepository.FindProducts(
		paginatedQuery.Offset(), paginatedQuery.Limit(),
		sortField, paginatedQuery.SortDirection,
		&repositoryFilters,
		nil,
	)
	if err != nil {
		return nil, err
	}

	count, err := s.ProductWithStatsRepository.CountProducts(&repositoryFilters, nil)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return pagination.NewPaginatedData[[]product_repository.ProductWithStats](
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
