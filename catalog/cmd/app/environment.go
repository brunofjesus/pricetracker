package main

import (
	"github.com/brunofjesus/pricetracker/catalog/internal/app"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	price_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	stats_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/stats"
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
	Finder  *product.Finder
}

type Price struct {
	Finder *price.Finder
}

type Store struct {
	Handler *store.Handler
	Finder  *store.Finder
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
	statsRepository := stats_repository.NewRepository(db)

	productCreator := product.Creator{
		DB:                    db,
		StoreRepository:       storeRepository,
		ProductRepository:     productRepository,
		ProductMetaRepository: productMetaRepository,
		PriceRepository:       priceRepository,
		StatsRepository:       statsRepository,
	}

	productUpdater := product.Updater{
		DB:                    db,
		StoreRepository:       storeRepository,
		ProductRepository:     productRepository,
		ProductMetaRepository: productMetaRepository,
		PriceRepository:       priceRepository,
		StatsRepository:       statsRepository,
	}

	productFinder := product.Finder{
		DB:                    db,
		MetricsRepository:     metricsRepository,
		ProductRepository:     productRepository,
		ProductMetaRepository: productMetaRepository,
	}

	productHandler := product.Handler{
		Matchers: []product.Matcher{
			&product.UrlMatcher{Finder: &productFinder},
			&product.EanMatcher{Finder: &productFinder},
			&product.SkuMatcher{Finder: &productFinder},
		},
		Creator: &productCreator,
		Updater: &productUpdater,
	}

	priceFinder := price.Finder{
		DB:         db,
		Repository: priceRepository,
	}

	storeHandler := store.Handler{
		Repository: storeRepository,
	}

	storeFinder := store.Finder{
		StoreRepository: storeRepository,
	}

	env := Environment{
		Product: Product{
			Creator: &productCreator,
			Handler: &productHandler,
			Finder:  &productFinder,
		},
		Price: Price{
			Finder: &priceFinder,
		},
		Store: Store{
			Handler: &storeHandler,
			Finder:  &storeFinder,
		},
	}

	return env
}
