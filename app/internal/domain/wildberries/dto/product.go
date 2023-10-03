package dto

type Product struct {
	Name         string
	Category     Category
	Url          string
	Price        float32
	Seller       string
	Rating       string
	ReviewCount  int
	PopularScore int
}
