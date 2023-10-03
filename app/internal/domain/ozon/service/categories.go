package service

import (
	"parser/internal/domain/ozon/dto"
)

type Categories interface {
	Iterate() <-chan *dto.Category
}
