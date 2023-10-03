package useCase

import (
	"context"
	"fmt"
	"parser/internal/domain/wildberries/dto"
	"parser/internal/domain/wildberries/entity"
	"parser/internal/domain/wildberries/helper"
	"parser/internal/domain/wildberries/repository"
	"parser/internal/domain/wildberries/service"
)

type CollectProducts struct {
	filter             []dto.ProductFilter
	products           service.Products
	categoryRepository repository.Category
	productRepository  repository.Product
	logger             service.Log
	ctx                context.Context
}

func NewCollectProducts(filter []dto.ProductFilter, products service.Products, categoryRepository repository.Category, productRepository repository.Product, logger service.Log, ctx context.Context) *CollectProducts {
	return &CollectProducts{
		filter:             filter,
		products:           products,
		categoryRepository: categoryRepository,
		productRepository:  productRepository,
		logger:             logger,
		ctx:                ctx,
	}
}

func (p *CollectProducts) Run() {
	categoryCh := p.categoryRepository.IterateActive()

	for {
		select {
		case <-p.ctx.Done():
			return
		case category, ok := <-categoryCh:
			if !ok {
				return
			}

			p.processCategory(*category, p.ctx)
		}
	}
}

func (p *CollectProducts) processCategory(category entity.Category, ctx context.Context) {
	productsCh := p.products.Iterate(
		dto.Category{
			Id:   category.ID,
			Url:  category.Url,
			Path: category.Path,
			Api:  category.Api,
		},
		p.getProductFilter(category.Path),
	)

	for {
		select {
		case <-ctx.Done():
			return
		case product, ok := <-productsCh:
			if !ok {
				return
			}

			if err := p.saveProduct(product); err != nil {
				p.logger.Error(fmt.Sprintf("Error while saving product `%s`", product.Name), err)
			}
		}
	}
}

func (p *CollectProducts) getProductFilter(categoryPath []string) dto.ProductFilter {
	res := dto.ProductFilter{}

	for _, productFilter := range p.filter {
		if "" == productFilter.InCategory {
			res = productFilter
		}

		for _, path := range categoryPath {
			if path == productFilter.InCategory {
				return productFilter
			}
		}
	}

	return res
}

func (p *CollectProducts) saveProduct(product *dto.Product) error {
	targetEntity := p.getEntity(product)

	targetEntity.AddPrice(product.Price)
	targetEntity.AddUrl(product.Url)
	targetEntity.AddRating(product.Rating, product.ReviewCount)
	targetEntity.AddSeller(product.Seller)
	targetEntity.AddPopularScore(product.PopularScore)

	return p.productRepository.Upsert(targetEntity)
}

func (p *CollectProducts) getEntity(product *dto.Product) *entity.Product {
	uuid := helper.GenerateProductUuid(*product)

	targetEntity, err := p.productRepository.Get(uuid.String())
	if err != nil {
		targetEntity = entity.NewProduct(uuid.String(), product.Name, product.Category.ShortPath)
	}

	return targetEntity
}
