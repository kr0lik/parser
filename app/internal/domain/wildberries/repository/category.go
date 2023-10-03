package repository

import (
	"parser/internal/domain/wildberries/entity"
	"time"
)

type Category interface {
	IterateActive() <-chan *entity.Category
	Upsert(category *entity.Category) error
	DisableOld(time time.Time) error
}
