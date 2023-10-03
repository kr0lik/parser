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

type Products struct {
	config wildberries.Products
	client *Client
	logger *log.Logger
	ctx    context.Context
}

func NewProducts(config wildberries.Products, client *Client, logger *log.Logger, ctx context.Context) *Products {
	return &Products{
		config: config,
		client: client,
		logger: logger,
		ctx:    ctx,
	}
}

func (p *Products) Iterate(category dto.Category, filter dto.ProductFilter) <-chan *dto.Product {
	ch := make(chan *dto.Product)

	go func() {
		defer close(ch)

		wg := &sync.WaitGroup{}

		p.collectProducts(&category, &filter, 1, ch, wg)

		wg.Wait()
	}()

	return ch
}

func (p *Products) getPageUrl(category *dto.Category, filter *dto.ProductFilter, page int) string {
	params := []string{p.config.SortParam}

	defaults := strings.Replace(p.config.FilterParam.Default, "<categoryId>", category.Id, 1)
	defaults = strings.Replace(defaults, "<page>", strconv.Itoa(page), 1)

	params = append(params, defaults)

	if filter.PriceMin > 0 || filter.PriceMax > 0 {
		from := "0"
		if filter.PriceMin > 0 {
			from = fmt.Sprintf("%d", filter.PriceMin*p.config.PriceFactor)
		}

		to := "999999999"
		if filter.PriceMax > 0 {
			to = fmt.Sprintf("%d", filter.PriceMax*p.config.PriceFactor)
		}

		priceParam := strings.Replace(p.config.FilterParam.PriceFromTo, "<from>", from, 1)
		priceParam = strings.Replace(priceParam, "<to>", to, 1)

		params = append(params, priceParam)
	}

	return category.Api + "?" + strings.Join(params, "&")
}

func (p *Products) collectProducts(category *dto.Category, filter *dto.ProductFilter, page int, ch chan *dto.Product, wg *sync.WaitGroup) {
	if category.Api == "" {
		p.logger.Warning("No api url", category.Path)
		return
	}

	apiUrl := p.getPageUrl(category, filter, page)

	data, err := p.client.get(apiUrl)
	if err != nil {
		p.logger.Error(fmt.Sprintf("error while calling products api page: %v", err), apiUrl)
		return
	}

	data, err = helper.InPath(data, p.config.List.Path)
	if err != nil {
		p.logger.Error(fmt.Sprintf("error while getting products: %v", err), apiUrl)
		return
	}

	items := data.([]interface{})

	if len(items) > 0 {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			page++

			p.collectProducts(category, filter, page, ch, wg)
		}(page)
	}

	for _, item := range items {
		idContent, err := helper.InPath(item, p.config.List.Attributes.Id)
		if err != nil {
			p.logger.Error(fmt.Sprintf("error while getting product id: %v", err))
			return
		}

		nameContent, err := helper.InPath(item, p.config.List.Attributes.Name)
		if err != nil {
			p.logger.Error(fmt.Sprintf("error while getting product name: %v", err))
			return
		}

		brandContent, err := helper.InPath(item, p.config.List.Attributes.Brand)
		if err != nil {
			p.logger.Error(fmt.Sprintf("error while getting product brand: %v", err), nameContent)
			return
		}

		priceContent, err := helper.InPath(item, p.config.List.Attributes.Price)
		if err != nil {
			p.logger.Error(fmt.Sprintf("error while getting product price: %v", err), nameContent)
			return
		}

		ratingContent, err := helper.InPath(item, p.config.List.Attributes.Rating)
		if err != nil {
			p.logger.Error(fmt.Sprintf("error while getting product rating: %v", err), nameContent)
			return
		}

		reviewCountContent, err := helper.InPath(item, p.config.List.Attributes.ReviewCount)
		if err != nil {
			p.logger.Error(fmt.Sprintf("error while getting product review count: %v", err), nameContent)
			return
		}

		popularScoreContent, err := helper.InPath(item, p.config.List.Attributes.PopularScore)
		if err != nil {
			p.logger.Error(fmt.Sprintf("error while getting product popular score: %v", err), nameContent)
			return
		}

		ch <- p.getProductDto(category, idContent, nameContent, brandContent, priceContent, ratingContent, reviewCountContent, popularScoreContent)
	}
}

func (p *Products) getProductDto(category *dto.Category, idContent, nameContent, brandContent, priceContent, ratingContent, reviewCountContent, popularScoreContent interface{}) *dto.Product {
	prettyId := strconv.Itoa(int(idContent.(float64)))

	prettyName := strings.Replace(p.config.List.Name, "<name>", fmt.Sprint(nameContent), 1)
	prettyName = strings.Replace(prettyName, "<brand>", fmt.Sprint(brandContent), 1)

	prettyUrl := strings.Replace(p.config.List.Url, "<id>", prettyId, 1)

	prettyPrice := float32(priceContent.(float64))
	prettyPrice = prettyPrice / float32(p.config.PriceFactor)

	prettyRating := fmt.Sprint(ratingContent)

	prettyReviewCount := int(reviewCountContent.(float64))

	prettyPopularScore := int(popularScoreContent.(float64))

	return &dto.Product{
		Name:         prettyName,
		Category:     *category,
		Url:          prettyUrl,
		Price:        prettyPrice,
		Rating:       prettyRating,
		ReviewCount:  prettyReviewCount,
		PopularScore: prettyPopularScore,
	}
}
