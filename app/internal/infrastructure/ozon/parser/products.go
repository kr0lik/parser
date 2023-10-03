package parser

import (
	"context"
	"fmt"
	seleniumBase "github.com/tebeka/selenium"
	"net/url"
	"parser/internal/config/ozon"
	"parser/internal/domain/ozon/dto"
	"parser/internal/infrastructure/log"
	"parser/internal/infrastructure/ozon/parser/helper"
	"parser/internal/infrastructure/selenium"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Products struct {
	chrome *selenium.Chrome
	config ozon.Products
	logger *log.Logger
	ctx    context.Context
}

func NewProducts(chrome *selenium.Chrome, config ozon.Products, logger *log.Logger, ctx context.Context) *Products {
	return &Products{
		chrome: chrome,
		config: config,
		logger: logger,
		ctx:    ctx,
	}
}

func (p *Products) Iterate(category dto.Category, filter dto.ProductFilter) <-chan *dto.Product {
	ch := make(chan *dto.Product)

	page := 1
	pages := map[int]string{
		page: p.getPageUrl(&category, &filter),
	}

	go func() {
		defer close(ch)

		total := 0
		retry := 0

		for {
			count, err := p.collectProducts(&category, &filter, page, &pages, ch)
			if err != nil {
				p.logger.Error(fmt.Sprintf("Error while getting products at page %d", page), err)
				if retry < 3 {
					p.logger.Debug(fmt.Sprintf("Retry getting products at page %d", page), err)
					continue
				}
			}

			total += count
			retry = 0
			page++

			if page > len(pages) {
				break
			}
		}

		p.logger.Debug(fmt.Sprintf("Total in `%s`", strings.Join(category.Path, "/")), total)
	}()

	return ch
}

func (p *Products) getPageUrl(category *dto.Category, filter *dto.ProductFilter) string {
	params := []string{p.config.SortParam}

	if filter.PriceMin > 0 || filter.PriceMax > 0 {
		from := "0"
		if filter.PriceMin > 0 {
			from = fmt.Sprintf("%d", filter.PriceMin)
		}

		to := "9999999"
		if filter.PriceMax > 0 {
			to = fmt.Sprintf("%d", filter.PriceMax)
		}

		priceParam := strings.Replace(p.config.FilterParam.PriceFromTo, "<from>", from, 1)
		priceParam = strings.Replace(priceParam, "<to>", to, 1)

		params = append(params, priceParam)
	}

	if filter.OnlyBestSale {
		params = append(params, p.config.FilterParam.OnlyBestSale)
	}

	return category.Url + "?" + strings.Join(params, "&")
}

func (p *Products) collectProducts(category *dto.Category, filter *dto.ProductFilter, currentPage int, pages *map[int]string, ch chan *dto.Product) (int, error) {
	p.logger.Debug(fmt.Sprintf("Loading products in `%s`", strings.Join(category.Path, "/")), currentPage, "/", len(*pages))

	currentUrl, isUrlExist := (*pages)[currentPage]
	if !isUrlExist {
		p.logger.Warning("End of pagination")
		return 0, nil
	}

	wd, err := p.chrome.NewWindow(currentUrl)
	if err != nil {
		return 0, fmt.Errorf("error while getting window with products: %v", err)
	}
	defer wd.Quit()

	p.collectPages(wd, pages)

	if err := wd.WaitWithTimeout(func(wd seleniumBase.WebDriver) (bool, error) {
		return helper.HasElement(wd, p.config.ListProducts.Wait.Element)
	}, time.Duration(p.config.ListProducts.Wait.Timeout)*time.Second); err != nil {
		return 0, fmt.Errorf("error while waiting load products: %v", err)
	}

	products, err := wd.FindElements(seleniumBase.ByCSSSelector, p.config.ListProducts.List.Element)
	if err != nil {
		return 0, fmt.Errorf("error while finding products: %v", err)
	}

	for i, product := range products {
		errors := make([]error, 0)

		nameContent, err := p.getNameContent(product)
		if err != nil {
			errors = append(errors, err)
		}

		hrefContent, err := p.getHrefContent(product)
		if err != nil {
			errors = append(errors, err)
		}

		priceContent, err := p.getPriceContent(product)
		if err != nil {
			errors = append(errors, err)
		}

		if len(errors) > 0 {
			return i + 1, fmt.Errorf("error while getting product attriburtes: %v", errors)
		}

		ratingContent, err := p.getRatingContent(product)
		if err != nil {
			p.logger.Warning(fmt.Sprintf("Error while get rating for product `%s` at `%s`", nameContent, currentUrl), err)
		}

		reviewCountContent, err := p.getReviewCountContent(product)
		if err != nil {
			p.logger.Warning(fmt.Sprintf("Error while get review count for product `%s` at `%s`", nameContent, currentUrl), err)
		}

		ch <- p.getProductDto(category, nameContent, hrefContent, priceContent, ratingContent, reviewCountContent, filter.OnlyBestSale)
	}

	return len(products), nil
}

func (p *Products) getNameContent(product seleniumBase.WebElement) (string, error) {
	for _, selector := range p.config.ListProducts.List.Name {
		text, err := helper.GetText(product, selector)
		if "" != text && err == nil {
			return text, nil
		}
	}

	return "", fmt.Errorf("no name found")
}

func (p *Products) getHrefContent(product seleniumBase.WebElement) (string, error) {
	for _, selector := range p.config.ListProducts.List.Url {
		text, err := helper.GetHref(product, selector)
		if "" != text && err == nil {
			hrefPattern := "([\\w\\-.,@?^=%&:/~+#]*[\\w\\-@?^=%&/~+#])"
			matched, _ := regexp.MatchString(hrefPattern, text)
			if matched {
				return text, nil
			}
		}
	}

	return "", fmt.Errorf("no href found")
}

func (p *Products) getPriceContent(product seleniumBase.WebElement) (string, error) {
	for _, selector := range p.config.ListProducts.List.Price {
		text, err := helper.GetText(product, selector)
		if "" != text && err == nil {
			pricePattern := "\\d+.\\x{20BD}"
			matched, _ := regexp.MatchString(pricePattern, text)
			if matched {
				return text, nil
			}
		}
	}

	return "", fmt.Errorf("no price found")
}

func (p *Products) getRatingContent(product seleniumBase.WebElement) (string, error) {
	for _, selector := range p.config.ListProducts.List.Rating {
		text, err := helper.GetText(product, selector)
		if "" != text && err == nil {
			ratingPattern := "^[.0-9]{1}[.,]{0,1}[0-9]{0,1}[\\s]{0,}$"
			matched, _ := regexp.MatchString(ratingPattern, text)
			if matched {
				return text, nil
			}
		}
	}

	return "", fmt.Errorf("no rating found")
}

func (p *Products) getReviewCountContent(product seleniumBase.WebElement) (string, error) {
	for _, selector := range p.config.ListProducts.List.ReviewCount {
		text, err := helper.GetText(product, selector)
		if "" != text && err == nil {
			ratingPattern := "^\\d+"
			matched, _ := regexp.MatchString(ratingPattern, text)
			if matched {
				return text, nil
			}
		}
	}

	return "", fmt.Errorf("no review count found")
}

func (p *Products) collectPages(wd seleniumBase.WebDriver, pages *map[int]string) {
	currentUrl, err := wd.CurrentURL()
	if err != nil {
		p.logger.Warning("Error while getting current url", err)
	}

	if err := wd.WaitWithTimeout(func(wd seleniumBase.WebDriver) (bool, error) {
		return helper.HasElement(wd, p.config.Pagination.Wait.Element)
	}, time.Duration(p.config.Pagination.Wait.Timeout)*time.Second); err != nil {
		p.logger.Debug(fmt.Sprintf("Pagination not loaded at `%s`", currentUrl), err)
		return
	}

	pagination, err := wd.FindElements(seleniumBase.ByCSSSelector, p.config.Pagination.List.Element)
	if err != nil {
		p.logger.Debug(fmt.Sprintf("No pagination found at `%s`", currentUrl), err)
		return
	}

	for _, page := range pagination {
		nameContent, err := helper.GetText(page, p.config.Pagination.List.Name)
		if err != nil {
			p.logger.Warning("Error while getting page num", err)
			continue
		}

		num, err := strconv.Atoi(nameContent)
		if err != nil {
			continue
		}

		hrefContent, err := helper.GetHref(page, p.config.Pagination.List.Url)
		if err != nil {
			p.logger.Warning("Error while getting page href", err)
			continue
		}

		baseUrl, _ := url.Parse(currentUrl)
		pageUrl, _ := url.Parse(hrefContent)
		prettyUrl := baseUrl.Scheme + "://" + baseUrl.Host + pageUrl.Path + "?" + pageUrl.RawQuery

		(*pages)[num] = prettyUrl
	}
}

func (p *Products) getProductDto(category *dto.Category, nameContent, hrefContent, priceContent, ratingContent, reviewCountContent string, onlyBestSale bool) *dto.Product {
	baseUrl, _ := url.Parse(category.Url)
	productUrl, _ := url.Parse(hrefContent)
	prettyUrl := baseUrl.Scheme + "://" + baseUrl.Host + productUrl.Path

	priceRegexp := regexp.MustCompile("[0-9,.]+")
	priceContentClear := strings.Join(priceRegexp.FindAllString(priceContent, -1), "")
	prettyPrice, _ := strconv.ParseFloat(priceContentClear, 8)

	prettyRating := strings.Replace(strings.Trim(ratingContent, " "), ",", ".", 1)

	prettyBestSale := dto.BestSaleUndefined
	if onlyBestSale {
		prettyBestSale = dto.BestSaleTrue
	}

	reviewCountPretty, _ := strconv.Atoi(strings.Trim(reviewCountContent, " "))

	return &dto.Product{
		Name:        nameContent,
		Category:    *category,
		Url:         prettyUrl,
		Price:       float32(prettyPrice),
		BestSale:    prettyBestSale,
		Rating:      prettyRating,
		ReviewCount: reviewCountPretty,
	}
}
