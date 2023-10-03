package api

import (
	"context"
	"fmt"
	"parser/internal/config/wildberries"
	"parser/internal/domain/wildberries/dto"
	"parser/internal/infrastructure/log"
	"parser/internal/infrastructure/wildberries/api/helper"
	"strconv"
	"strings"
	"sync"
)

type Categories struct {
	config wildberries.Categories
	client *Client
	logger *log.Logger
	ctx    context.Context
}

func NewCategories(config wildberries.Categories, client *Client, logger *log.Logger, ctx context.Context) *Categories {
	return &Categories{
		config: config,
		client: client,
		logger: logger,
		ctx:    ctx,
	}
}

func (c *Categories) Iterate() <-chan *dto.Category {
	ch := make(chan *dto.Category)

	go func() {
		defer close(ch)

		wg := &sync.WaitGroup{}

		if err := c.collectCategories(ch, wg); err != nil {
			c.logger.Error("Error while looping categories", err)
		}

		wg.Wait()
	}()

	return ch
}

func (c *Categories) collectCategories(ch chan *dto.Category, wg *sync.WaitGroup) error {
	data, err := c.client.get(c.config.Url)
	if err != nil {
		return fmt.Errorf("error while calling main categories api page: %v", err)
	}

	data, err = helper.InPath(data, c.config.ListMainCategories.Path)
	if err != nil {
		return fmt.Errorf("error while getting main categories: %v", err)
	}

	for _, item := range data.([]interface{}) {
		id, err := helper.InPath(item, c.config.ListMainCategories.Attributes.Id)
		if err != nil {
			return fmt.Errorf("error while getting main category id: %v", err)
		}

		name, err := helper.InPath(item, c.config.ListMainCategories.Attributes.Name)
		if err != nil {
			return fmt.Errorf("error while getting main category name: %v", err)
		}

		title, err := helper.InPath(item, c.config.ListMainCategories.Attributes.Title)
		if err != nil {
			return fmt.Errorf("error while getting main category (%s) title: %v", name, err)
		}

		url, err := helper.InPath(item, c.config.LoopChildCategories.Attributes.Url)
		if err != nil {
			return fmt.Errorf("error while getting main category (%s) url: %v", name, err)
		}

		shard, err := helper.InPath(item, c.config.ListMainCategories.Attributes.Shard)
		if err != nil {
			return fmt.Errorf("error while getting main category (%s) shard: %v", name, err)
		}

		child, err := helper.InPath(item, c.config.ListMainCategories.Attributes.Child)
		if err != nil {
			return fmt.Errorf("error while getting main category (%s) child: %v", name, err)
		}

		category := dto.Category{
			Id:        fmt.Sprint(id),
			Path:      []string{fmt.Sprint(title)},
			ShortPath: []string{fmt.Sprint(name)},
			Url:       fmt.Sprint(url),
			Api:       strings.Replace(c.config.ProductsApi, "<categoryShard>", fmt.Sprint(shard), 1),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			c.loopChild(ch, category, child, wg)
		}()
	}

	return nil
}

func (c *Categories) loopChild(ch chan *dto.Category, currentCategory dto.Category, child interface{}, wg *sync.WaitGroup) {
	data, err := helper.InPath(child, c.config.LoopChildCategories.Path)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error while getting child categories in %s: %v", currentCategory.Url, err))
		return
	}

	for _, item := range data.([]interface{}) {
		id, err := helper.InPath(item, c.config.LoopChildCategories.Attributes.Id)
		if err != nil {
			c.logger.Error(fmt.Sprintf("error while getting child category id: %v", err))
			return
		}

		name, err := helper.InPath(item, c.config.LoopChildCategories.Attributes.Name)
		if err != nil {
			c.logger.Error(fmt.Sprintf("error while getting child category name: %v", err))
			return
		}

		title, err := helper.InPath(item, c.config.LoopChildCategories.Attributes.Title)
		if err != nil {
			title = name
			c.logger.Warning(fmt.Sprintf("error while getting child category (%v) title: %v", append(currentCategory.Path, fmt.Sprint(name)), err))
		}

		url, err := helper.InPath(item, c.config.LoopChildCategories.Attributes.Url)
		if err != nil {
			c.logger.Error(fmt.Sprintf("error while getting child category (%v) url: %v", append(currentCategory.Path, fmt.Sprint(name)), err))
			return
		}

		shard, err := helper.InPath(item, c.config.LoopChildCategories.Attributes.Shard)
		if err != nil {
			c.logger.Error(fmt.Sprintf("error while getting child category (%v) shard: %v", append(currentCategory.Path, fmt.Sprint(name)), err))
		}

		api := ""
		if shard != "" {
			api = strings.Replace(c.config.ProductsApi, "<categoryShard>", fmt.Sprint(shard), 1)
		}

		category := dto.Category{
			Id:        strconv.Itoa(int(id.(float64))),
			Path:      append(currentCategory.Path, fmt.Sprint(title)),
			ShortPath: append(currentCategory.ShortPath, fmt.Sprint(name)),
			Url:       fmt.Sprint(url),
			Api:       api,
		}

		child, err := helper.InPath(item, c.config.LoopChildCategories.Attributes.Child)
		if err != nil {
			ch <- &category
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			c.loopChild(ch, category, child, wg)
		}()
	}
}
