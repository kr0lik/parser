package service

import (
	"parser/internal/domain/wildberries/dto"
)

type Categories interface {
	Iterate() <-chan *dto.Category
}
