package useCase

import (
	"context"
	"fmt"
	"parser/internal/domain/wildberries/dto"
	"parser/internal/domain/wildberries/entity"
	"parser/internal/domain/wildberries/repository"
	"parser/internal/domain/wildberries/service"
	"strings"
	"time"
)

type CollectCategories struct {
	filter             dto.CategoryFilter
	categories         service.Categories
	categoryRepository repository.Category
	logger             service.Log
	ctx                context.Context
}

func NewCollectCategories(filter dto.CategoryFilter, categories service.Categories, categoryRepository repository.Category, logger service.Log, ctx context.Context) *CollectCategories {
	return &CollectCategories{
		filter:             filter,
		categories:         categories,
		categoryRepository: categoryRepository,
		logger:             logger,
		ctx:                ctx,
	}
}

func (c *CollectCategories) Run() {
	checkTime := time.Now()

	categoryCh := c.categories.Iterate()

Loop:
	for {
		select {
		case <-c.ctx.Done():
			return
		case category, ok := <-categoryCh:
			if !ok {
				break Loop
			}

			if !c.isIncluded(category) || c.isExcluded(category) {
				break
			}

			if err := c.saveCategory(category, checkTime); err != nil {
				c.logger.Error(fmt.Sprintf("Error while saving category `%s`", strings.Join(category.Path, "/")), err)
			}
		}
	}

	if err := c.categoryRepository.DisableOld(checkTime); err != nil {
		c.logger.Fatal("Error while disabling old categories:", err)
	}
}

func (c *CollectCategories) isIncluded(category *dto.Category) bool {
	for _, pathPart := range category.Path {
		for _, inc := range c.filter.Include {
			if strings.ToLower(pathPart) == strings.ToLower(inc) {
				return true
			}
		}
	}

	return false
}

func (c *CollectCategories) isExcluded(category *dto.Category) bool {
	for _, pathPart := range category.Path {
		for _, inc := range c.filter.Exclude {
			if strings.ToLower(pathPart) == strings.ToLower(inc) {
				return true
			}
		}
	}

	return false
}

func (c *CollectCategories) saveCategory(category *dto.Category, checkTime time.Time) error {
	return c.categoryRepository.Upsert(
		entity.NewCategory(
			category.Id,
			category.Url,
			category.Api,
			category.Path,
			category.ShortPath,
			checkTime,
		),
	)
}
