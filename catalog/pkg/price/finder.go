package price

import (
	"database/sql"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	"time"
)

type Finder struct {
	DB         *sql.DB
	Repository *price.Repository
}

type Price struct {
	DateTime time.Time `json:"date_time,omitempty"`
	Value    int       `json:"price,omitempty"`
}

func (s *Finder) FindPriceHistoryBetween(productId int64, from time.Time, to time.Time, tx *sql.Tx) ([]Price, error) {
	//TODO: limit interval

	prices, err := s.Repository.FindPricesBetween(productId, from, to, tx)

	if err != nil {
		return nil, err
	}

	result := make([]Price, len(prices))
	for i, productPrice := range prices {
		result[i] = Price{
			DateTime: productPrice.DateTime,
			Value:    productPrice.Price,
		}
	}

	return result, nil
}
