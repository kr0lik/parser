package entity

import (
	"time"
)

type Category struct {
	ID        string `bson:"_id"`
	Url       string
	Api       string
	Path      []string
	ShortPath []string
	Active    bool
	Checked   time.Time
}

func NewCategory(uuid, url, api string, path, shortPath []string, time time.Time) *Category {
	return &Category{
		ID:        uuid,
		Path:      path,
		ShortPath: shortPath,
		Url:       url,
		Api:       api,
		Active:    true,
		Checked:   time,
	}
}
