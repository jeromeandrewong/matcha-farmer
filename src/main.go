package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly/v2"
)

type Product struct {
	Title       string
	Status      string
	LastChecked time.Time
}

func sendTeleNotification(message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return fmt.Errorf("telegram bot token or chat ID is not set")
	}

	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id": {chatID},
			"text":    {message},
		})

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", response.Status)
	}

	return nil

}

func HandleRequest(ctx context.Context) (string, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.marukyu-koyamaen.co.jp"),
	)

	products := []Product{}

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		title := e.ChildAttr("a.woocommerce-loop-product__link", "title")
		status := "❌ Out of Stock"
		if !strings.Contains(e.Attr("class"), "outofstock") {
			status = "✅ In Stock"
		}
		products = append(products, Product{
			Title:       title,
			Status:      status,
			LastChecked: time.Now(),
		})
	})

	err := c.Visit("https://www.marukyu-koyamaen.co.jp/english/shop/products/category/matcha/principal/")
	if err != nil {
		return "", err
	}

	singaporeLocation, _ := time.LoadLocation("Asia/Singapore")
	singaporeTime := time.Now().In(singaporeLocation)

	message := "Marukyu Koyamaen Stock Check:\n\n"

	message += fmt.Sprintf("🕜 Last Checked: %s\n\n", singaporeTime.Format("Mon, 2 Jan 3:04 PM"))

	for _, product := range products {
		message += fmt.Sprintf("🍵 Name: %s\n📦 Status: %s\n\n",
			product.Title,
			product.Status)
	}

	err = sendTeleNotification(message)
	if err != nil {
		return "", fmt.Errorf("failed to send Telegram notification: %v", err)
	}

	return "OK", nil
}

func main() {
	lambda.Start(HandleRequest)
}
