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
	"strings"
	"time"
)

type Categories struct {
	chrome *selenium.Chrome
	config ozon.Categories
	logger *log.Logger
	ctx    context.Context
}

func NewCategories(chrome *selenium.Chrome, config ozon.Categories, logger *log.Logger, ctx context.Context) *Categories {
	return &Categories{
		chrome: chrome,
		config: config,
		logger: logger,
		ctx:    ctx,
	}
}

func (c *Categories) Iterate() <-chan *dto.Category {
	mainCategories, err := c.getRootCategories()
	if err != nil {
		c.logger.Fatal("Error while getting main categories:", err)
	}

	ch := make(chan *dto.Category)

	go func() {
		defer close(ch)

		for _, mainCategory := range mainCategories {
			if _, err := c.collectChildRecursively(mainCategory, ch); err != nil {
				c.logger.Error("Error while looping categories", err)
			}
		}
	}()

	return ch
}

func (c *Categories) getRootCategories() ([]*dto.Category, error) {
	c.logger.Debug("Getting root categories")

	wd, err := c.chrome.NewWindow(c.config.StartUrl)
	if err != nil {
		return nil, fmt.Errorf("error while getting window with root categoriesa: %v", err)
	}
	defer wd.Quit()

	if err := wd.WaitWithTimeout(func(wd seleniumBase.WebDriver) (bool, error) {
		return helper.HasElement(wd, c.config.ListMainCategories.Wait.Element)
	}, time.Duration(c.config.ListMainCategories.Wait.Timeout)*time.Second); err != nil {
		return nil, fmt.Errorf("error while witing load root categoriesa: %v", err)
	}

	categories, err := wd.FindElements(seleniumBase.ByCSSSelector, c.config.ListMainCategories.List.Element)
	if err != nil {
		return nil, fmt.Errorf("error while finding root categoriesa: %v", err)
	}

	res := make([]*dto.Category, 0)

	for _, category := range categories {
		nameContent, err := helper.GetText(category, c.config.ListMainCategories.List.Name)
		if err != nil {
			return nil, fmt.Errorf("error while while geting root category name: %v", err)
		}

		hrefContent, err := helper.GetHref(category, c.config.ListMainCategories.List.Url)
		if err != nil {
			return nil, fmt.Errorf("error while geting root category href: %v", err)
		}

		res = append(res, c.getCategoryDto([]string{nameContent}, []string{nameContent}, hrefContent))
	}

	return res, nil
}

func (c *Categories) collectChildRecursively(baseCategory *dto.Category, ch chan *dto.Category) ([]*dto.Category, error) {
	childCategories, err := c.getNameAndChildCategories(baseCategory)
	if err != nil {
		return nil, fmt.Errorf("error while geting child categories: %v", err)
	}

	res := make([]*dto.Category, 0)

	for _, childCategory := range childCategories {
		subChildrenCategories, err := c.collectChildRecursively(childCategory, ch)
		if err != nil {
			return nil, fmt.Errorf("error while geting child sub categories: %v", err)
		}

		if len(subChildrenCategories) > 0 {
			res = append(res, subChildrenCategories...)

			continue
		}

		ch <- childCategory
	}

	return res, nil
}

func (c *Categories) isBaseCategory(category string, baseCategory *dto.Category) bool {
	for _, categoryInPath := range baseCategory.Path {
		if strings.ToLower(category) == strings.ToLower(categoryInPath) {
			return true
		}
	}

	return false
}

func (c *Categories) getNameAndChildCategories(baseCategory *dto.Category) ([]*dto.Category, error) {
	c.logger.Debug(fmt.Sprintf("Getting child in `%s`", strings.Join(baseCategory.Path, "/")))

	wd, err := c.chrome.NewWindow(baseCategory.Url)
	if err != nil {
		return nil, fmt.Errorf("error while geting window with child categories: %v", err)
	}
	defer wd.Quit()

	title, err := wd.FindElement(seleniumBase.ByCSSSelector, c.config.LoopCategories.Title.Element)
	if err != nil {
		return nil, fmt.Errorf("error while getting child category title: %v", err)
	}

	name, err := title.Text()
	if err != nil {
		return nil, fmt.Errorf("error while getting child category title text: %v", err)
	}

	baseCategory.Path[len(baseCategory.Path)-1] = name

	if err := wd.WaitWithTimeout(func(wd seleniumBase.WebDriver) (bool, error) {
		return helper.HasElement(wd, c.config.LoopCategories.Wait.Element)
	}, time.Duration(c.config.LoopCategories.Wait.Timeout)*time.Second); err != nil {
		return nil, fmt.Errorf("error while waiting load child categories: %v", err)
	}

	categories, err := wd.FindElements(seleniumBase.ByCSSSelector, c.config.LoopCategories.List.Element)
	if err != nil {
		return nil, err
	}

	res := make([]*dto.Category, 0)

	for i, category := range categories {
		if i < len(baseCategory.Path) {
			continue
		}

		nameContent, err := helper.GetText(category, c.config.LoopCategories.List.Name)
		if err != nil {
			return nil, fmt.Errorf("error while getting child category name: %v", err)
		}

		if c.isBaseCategory(nameContent, baseCategory) {
			continue
		}

		hrefContent, err := helper.GetHref(category, c.config.LoopCategories.List.Url)
		if err != nil {
			return nil, fmt.Errorf("error while getting child category href: %v", err)
		}

		path := make([]string, 0)
		path = append(path, baseCategory.Path...)
		path = append(path, nameContent)

		shortPath := make([]string, 0)
		shortPath = append(shortPath, baseCategory.ShortPath...)
		shortPath = append(shortPath, nameContent)

		res = append(res, c.getCategoryDto(path, shortPath, hrefContent))
	}

	return res, nil
}

func (c *Categories) getCategoryDto(path, shortPath []string, hrefContent string) *dto.Category {
	idRegexp := regexp.MustCompile("([0-9]+)\\/{0,1}$")
	id := idRegexp.FindStringSubmatch(hrefContent)[1]

	baseUrl, _ := url.Parse(c.config.StartUrl)
	currentUrl, _ := url.Parse(hrefContent)
	prettyUrl := baseUrl.Scheme + "://" + baseUrl.Host + currentUrl.Path

	return &dto.Category{Id: id, Path: path, ShortPath: shortPath, Url: prettyUrl}
}
