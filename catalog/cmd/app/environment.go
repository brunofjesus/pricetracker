package main

import (
	"github.com/brunofjesus/pricetracker/catalog/internal/app"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	price_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	store_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
	"github.com/brunofjesus/pricetracker/catalog/pkg/price"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/store"
)

type Environment struct {
	Product Product
	Price   Price
	Store   Store
}

type Product struct {
	Creator *product.Creator
	Handler *product.Handler
	Matcher *product.Matcher
	Finder  *product.Finder
}

type Price struct {
	Finder *price.Finder
}

type Store struct {
	Handler *store.Handler
}

func loadEnvironment(appConfig *app.ApplicationConfiguration) Environment {
	db := repository.Connect(
		appConfig.Database.DSN, appConfig.Database.Attempts,
	)

	storeRepository := store_repository.NewRepository(db)
	productRepository := product_repository.NewRepository(db)
	productMetaRepository := product_repository.NewMetaRepository(db)
	priceRepository := price_repository.NewRepository(db)
	metricsRepository := product_repository.NewMetricsRepository(db)

	productMatcher := product.Matcher{
		StoreRepository:       storeRepository,
		ProductRepository:     productRepository,
		ProductMetaRepository: productMetaRepository,
		PriceRepository:       priceRepository,
	}

	productCreator := product.Creator{
		DB:                    db,
		StoreRepository:       storeRepository,
		ProductRepository:     productRepository,
		ProductMetaRepository: productMetaRepository,
		PriceRepository:       priceRepository,
	}

	productUpdater := product.Updater{
		DB:                    db,
		StoreRepository:       storeRepository,
		ProductRepository:     productRepository,
		ProductMetaRepository: productMetaRepository,
		PriceRepository:       priceRepository,
	}

	productHandler := product.Handler{
		Matcher: &productMatcher,
		Creator: &productCreator,
		Updater: &productUpdater,
	}

	productFinder := product.Finder{
		DB:         db,
		Repository: metricsRepository,
	}

	priceFinder := price.Finder{
		DB:         db,
		Repository: priceRepository,
	}

	storeHandler := store.Handler{
		Repository: storeRepository,
	}

	env := Environment{
		Product: Product{
			Creator: &productCreator,
			Handler: &productHandler,
			Matcher: &productMatcher,
			Finder:  &productFinder,
		},
		Price: Price{
			Finder: &priceFinder,
		},
		Store: Store{
			Handler: &storeHandler,
		},
	}

	return env
}
