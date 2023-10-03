package repository

import (
	"parser/internal/domain/ozon/entity"
)

type Product interface {
	Get(uuid string) (*entity.Product, error)
	Upsert(*entity.Product) error
}
