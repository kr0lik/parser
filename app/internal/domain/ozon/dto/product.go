package dto

const BestSaleFalse = 0
const BestSaleTrue = 1
const BestSaleUndefined = 2

type Product struct {
	Name        string
	Category    Category
	Url         string
	Price       float32
	BestSale    int
	Seller      string
	Rating      string
	ReviewCount int
}
