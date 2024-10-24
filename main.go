package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	Title       string
	Status      string
	LastChecked time.Time
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.marukyu-koyamaen.co.jp"),
	)

	products := []Product{}

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		title := e.ChildAttr("a.woocommerce-loop-product__link", "title")
		status := "Out of Stock"
		if !strings.Contains(e.Attr("class"), "outofstock") {
			status = "In Stock"
		}
		products = append(products, Product{
			Title:       title,
			Status:      status,
			LastChecked: time.Now(),
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://www.marukyu-koyamaen.co.jp/english/shop/products/category/matcha/principal/")
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range products {
		fmt.Printf("Title: %s, Status: %s, Last Checked: %s\n", product.Title, product.Status, product.LastChecked)
	}
}
