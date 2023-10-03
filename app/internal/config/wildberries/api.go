package wildberries

import (
	"gopkg.in/yaml.v3"
	"os"
)

var api *Api

type Api struct {
	Url        string     `yaml:"url"`
	Categories Categories `yaml:"categories"`
	Products   Products   `yaml:"products"`
}

type Categories struct {
	SiteUrl            string `yaml:"siteUrl"`
	Url                string `yaml:"url"`
	ListMainCategories struct {
		Path       string `yaml:"path"`
		Attributes struct {
			Id    string `yaml:"id"`
			Name  string `yaml:"name"`
			Title string `yaml:"title"`
			Url   string `yaml:"url"`
			Shard string `yaml:"shard"`
			Child string `yaml:"child"`
		} `yaml:"attributes"`
	} `yaml:"listMainCategories"`
	LoopChildCategories struct {
		Path       string `yaml:"path"`
		Attributes struct {
			Id    string `yaml:"id"`
			Name  string `yaml:"name"`
			Title string `yaml:"title"`
			Url   string `yaml:"url"`
			Shard string `yaml:"shard"`
			Child string `yaml:"child"`
		} `yaml:"attributes"`
	} `yaml:"loopChildCategories"`
	ProductsApi string `yaml:"productsApi"`
}

type Products struct {
	SortParam   string `yaml:"sortParam"`
	FilterParam struct {
		Default     string `yaml:"default"`
		PriceFromTo string `yaml:"priceFromTo"`
	} `yaml:"filterParam"`
	List struct {
		Path       string `yaml:"path"`
		Attributes struct {
			Id           string `yaml:"id"`
			Name         string `yaml:"name"`
			Brand        string `yaml:"brand"`
			Price        string `yaml:"price"`
			Rating       string `yaml:"rating"`
			ReviewCount  string `yaml:"reviewCount"`
			PopularScore string `yaml:"popularScore"`
		} `yaml:"attributes"`
		Url  string `yaml:"url"`
		Name string `yaml:"name"`
	} `yaml:"list"`
	PriceFactor int `yaml:"priceFactor"`
}

func ProvideCategories() Categories {
	return api.Categories
}

func ProvideProducts() Products {
	return api.Products
}

func ReadApi() {
	api = new(Api)

	data, err := os.ReadFile("config/wildberries/api.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, api)
	if err != nil {
		panic(err)
	}
}
