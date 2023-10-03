package wildberries

import (
	"gopkg.in/yaml.v3"
	"os"
	"parser/internal/domain/wildberries/dto"
)

var filter *Filter

type Filter struct {
	Category dto.CategoryFilter  `yaml:"category"`
	Product  []dto.ProductFilter `yaml:"product"`
}

func ProvideCategoryFilter() dto.CategoryFilter {
	return filter.Category
}

func ProvideProductFilter() []dto.ProductFilter {
	return filter.Product
}

func ReadFilter() {
	filter = new(Filter)

	data, err := os.ReadFile("config/wildberries/filter.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, filter)
	if err != nil {
		panic(err)
	}
}
