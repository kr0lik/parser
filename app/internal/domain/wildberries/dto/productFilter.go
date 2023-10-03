package dto

type ProductFilter struct {
	InCategory string `yaml:"inCategory"`
	PriceMin   int    `yaml:"priceMin"`
	PriceMax   int    `yaml:"priceMax"`
}
