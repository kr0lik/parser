package entity

import "time"

type Product struct {
	ID           string `bson:"_id"`
	Name         string
	CategoryPath []string
	Url          []string
	Price        []struct {
		Date  time.Time
		Price float32
	}
	BestSale []time.Time
	Seller   []string
	Rating   []struct {
		Date        time.Time
		Score       string
		ReviewCount int
	}
}

func NewProduct(uuid, name string, categoryPath []string) *Product {
	return &Product{
		ID:           uuid,
		Name:         name,
		CategoryPath: categoryPath,
		Price: make([]struct {
			Date  time.Time
			Price float32
		}, 0),
		Url:      make([]string, 0),
		BestSale: make([]time.Time, 0),
		Seller:   make([]string, 0),
		Rating: make([]struct {
			Date        time.Time
			Score       string
			ReviewCount int
		}, 0),
	}
}

func (p *Product) AddPrice(price float32) {
	if !(price > 0) {
		return
	}

	result := make([]struct {
		Date  time.Time
		Price float32
	}, 0)

	if len(p.Price) > 0 {
		for _, p := range p.Price {
			if time.Now().Format("2006-01-02") == p.Date.Format("2006-01-02") {
				continue
			}

			result = append(result, p)
		}
	}

	p.Price = append(result, struct {
		Date  time.Time
		Price float32
	}{Date: time.Now(), Price: price})
}

func (p *Product) AddUrl(url string) {
	if "" == url {
		return
	}

	isNotDuplicate := true

	if len(p.Url) > 0 {
		for _, u := range p.Url {
			if u != url {
				isNotDuplicate = false
				break
			}
		}
	}

	if isNotDuplicate {
		p.Url = append(p.Url, url)
	}
}

func (p *Product) AddRating(score string, reviewCount int) {
	if "" == score {
		return
	}

	result := make([]struct {
		Date        time.Time
		Score       string
		ReviewCount int
	}, 0)

	if len(p.Rating) > 0 {
		for _, r := range p.Rating {
			if time.Now().Format("2006-01-02") == r.Date.Format("2006-01-02") {
				continue
			}

			result = append(result, r)
		}
	}

	p.Rating = append(result, struct {
		Date        time.Time
		Score       string
		ReviewCount int
	}{Date: time.Now(), Score: score, ReviewCount: reviewCount})
}

func (p *Product) AddBestSale(bestSale bool) {
	if !bestSale {
		return
	}

	result := make([]time.Time, 0)

	if len(p.BestSale) > 0 {
		for _, bs := range p.BestSale {
			if time.Now().Format("2006-01-02") == bs.Format("2006-01-02") {
				continue
			}

			result = append(result, bs)
		}
	}

	p.BestSale = append(result, time.Now())
}

func (p *Product) AddSeller(seller string) {
	if "" == seller {
		return
	}

	isNotDuplicate := true

	if len(p.Seller) > 0 {
		for _, s := range p.Seller {
			if s != seller {
				isNotDuplicate = false
				break
			}
		}
	}

	if isNotDuplicate {
		p.Seller = append(p.Seller, seller)
	}
}
