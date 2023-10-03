//go:build wireinject
// +build wireinject

package wildberries

import (
	"context"
	"github.com/google/wire"
	"parser/internal/config"
	"parser/internal/config/wildberries"
	"parser/internal/domain/wildberries/repository"
	"parser/internal/domain/wildberries/service"
	"parser/internal/domain/wildberries/useCase"
	"parser/internal/infrastructure/log"
	"parser/internal/infrastructure/mongodb"
	"parser/internal/infrastructure/wildberries/api"
	"parser/internal/infrastructure/wildberries/storage"
)

var categoriesParserSet = wire.NewSet(
	api.NewCategories,
	api.NewClient,
	wildberries.ProvideCategories,
	wire.Bind(new(service.Categories), new(*api.Categories)),
)

var productsParserSet = wire.NewSet(
	api.NewProducts,
	api.NewClient,
	wildberries.ProvideProducts,
	wire.Bind(new(service.Products), new(*api.Products)),
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
		wildberries.ProvideCategoryFilter,
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
		wildberries.ProvideProductFilter,
		productsParserSet,
		categoryRepositorySet,
		productRepositorySet,
		mongoDbSet,
		loggerSet,
	)
	return &useCase.CollectProducts{}, nil
}
