//go:build wireinject
// +build wireinject

package ozon

import (
	"context"
	"github.com/google/wire"
	"parser/internal/config"
	"parser/internal/config/ozon"
	"parser/internal/domain/ozon/repository"
	"parser/internal/domain/ozon/service"
	"parser/internal/domain/ozon/useCase"
	"parser/internal/infrastructure/log"
	"parser/internal/infrastructure/mongodb"
	"parser/internal/infrastructure/ozon/parser"
	"parser/internal/infrastructure/ozon/storage"
	"parser/internal/infrastructure/selenium"
)

var chromeSet = wire.NewSet(
	selenium.NewChrome,
	config.ProvideSeleniumServerOptions,
)

var categoriesParserSet = wire.NewSet(
	parser.NewCategories,
	ozon.ProvideCategories,
	chromeSet,
	wire.Bind(new(service.Categories), new(*parser.Categories)),
)

var productsParserSet = wire.NewSet(
	parser.NewProducts,
	ozon.ProvideProducts,
	chromeSet,
	wire.Bind(new(service.Products), new(*parser.Products)),
)

var mongoDbSet = wire.NewSet(
	mongodb.NewClient,
	config.ProvideMongodbClientOptions,
)

var categoryRepositorySet = wire.NewSet(
	storage.NewCategory,
	wire.Bind(new(repository.Category), new(*storage.Category)),
)

var productRepositorySet = wire.NewSet(
	storage.NewProduct,
	wire.Bind(new(repository.Product), new(*storage.Product)),
)

var loggerSet = wire.NewSet(
	log.NewLogger,
	wire.Bind(new(service.Log), new(*log.Logger)),
)

func InitialiseCollectCategories(ctx context.Context) (*useCase.CollectCategories, error) {
	wire.Build(
		useCase.NewCollectCategories,
		ozon.ProvideCategoryFilter,
		categoriesParserSet,
		categoryRepositorySet,
		mongoDbSet,
		loggerSet,
	)
	return &useCase.CollectCategories{}, nil
}

func InitialiseCollectProducts(ctx context.Context) (*useCase.CollectProducts, error) {
	wire.Build(
		useCase.NewCollectProducts,
		ozon.ProvideProductFilter,
		productsParserSet,
		categoryRepositorySet,
		productRepositorySet,
		mongoDbSet,
		loggerSet,
	)
	return &useCase.CollectProducts{}, nil
}
