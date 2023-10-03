package ozon

import (
	"gopkg.in/yaml.v3"
	"os"
)

var parser *Parser

type Parser struct {
	Url        string     `yaml:"url"`
	Region     Region     `yaml:"region"`
	Categories Categories `yaml:"categories"`
	Products   Products   `yaml:"products"`
}

type Region struct {
	CheckRegion struct {
		Wait struct {
			Element string `yaml:"element"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"wait"`
		Check struct {
			Element string `yaml:"element"`
			IsName  string `yaml:"isName"`
		} `yaml:"check"`
	} `yaml:"checkRegion"`

	CallChoosePopup struct {
		Click string `yaml:"click"`
	} `yaml:"callChoosePopup"`

	CallChooser struct {
		Wait struct {
			Element string `yaml:"element"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"wait"`
		Click string `yaml:"click"`
	} `yaml:"callChooser"`

	ChooseRegion struct {
		Wait struct {
			Element string `yaml:"element"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"wait"`
		List struct {
			Element string `yaml:"element"`
			Name    string `yaml:"name"`
		} `yaml:"list"`
	} `yaml:"chooseRegion"`
}

type Categories struct {
	StartUrl string `yaml:"startUrl"`

	ListMainCategories struct {
		Wait struct {
			Element string `yaml:"element"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"wait"`
		List struct {
			Element string `yaml:"element"`
			Name    string `yaml:"name"`
			Url     string `yaml:"url"`
		} `yaml:"list"`
	} `yaml:"listMainCategories"`

	LoopCategories struct {
		Wait struct {
			Element string `yaml:"element"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"wait"`
		List struct {
			Element string `yaml:"element"`
			Name    string `yaml:"name"`
			Url     string `yaml:"url"`
		} `yaml:"list"`
		Title struct {
			Element string `yaml:"element"`
		} `yaml:"title"`
	} `yaml:"loopCategories"`
}

type Products struct {
	SortParam   string `yaml:"sortParam"`
	FilterParam struct {
		PriceFromTo  string `yaml:"priceFromTo"`
		OnlyBestSale string `yaml:"onlyBestSale"`
	} `yaml:"filterParam"`
	ListProducts struct {
		Wait struct {
			Element string `yaml:"element"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"wait"`
		List struct {
			Element     string   `yaml:"element"`
			Name        []string `yaml:"name"`
			Url         []string `yaml:"url"`
			Price       []string `yaml:"price"`
			Rating      []string `yaml:"rating"`
			ReviewCount []string `yaml:"reviewCount"`
		} `yaml:"list"`
	} `yaml:"listProducts"`
	Pagination struct {
		Wait struct {
			Element string `yaml:"element"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"wait"`
		List struct {
			Element string `yaml:"element"`
			Name    string `yaml:"name"`
			Url     string `yaml:"url"`
		} `yaml:"list"`
	} `yaml:"pagination"`
}

func ProvideCategories() Categories {
	return parser.Categories
}

func ProvideProducts() Products {
	return parser.Products
}

func ReadParser() {
	parser = new(Parser)

	data, err := os.ReadFile("config/ozon/parser.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, parser)
	if err != nil {
		panic(err)
	}
}
