package service

import (
	"parser/internal/domain/ozon/dto"
)

type Products interface {
	Iterate(category dto.Category, filter dto.ProductFilter) <-chan *dto.Product
}
