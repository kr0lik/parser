package entity

import (
	"time"
)

type Category struct {
	ID        string `bson:"_id"`
	Url       string
	Path      []string
	ShortPath []string
	Active    bool
	Checked   time.Time
}

func NewCategory(uuid, url string, path, shortPath []string, time time.Time) *Category {
	return &Category{
		ID:        uuid,
		Path:      path,
		ShortPath: shortPath,
		Url:       url,
		Active:    true,
		Checked:   time,
	}
}
