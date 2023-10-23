package crawler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/brunofjesus/pricetracker/stores/pingodoce/config"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/definition"
)

func Crawl(logger *slog.Logger, store definition.Store, publishFunc func(definition.StoreProduct)) {
	appConfig := config.GetApplicationConfiguration()

	categories, err := fetchCategories()
	if err != nil {
		logger.Error("Cannot get categories", slog.Any("error", err))
		return
	}

	currentCategoryIdx := 0
	totalCategories := len(categories.Tree)

	for _, v := range categories.Tree {
		currentCategoryIdx = currentCategoryIdx + 1
		slog.Info(
			"Crawling category",
			slog.Group("category",
				slog.String("slug", v.Slug),
				slog.String("id", v.ID),
			),
			slog.Group("progress",
				slog.Int("current", currentCategoryIdx),
				slog.Int("total", totalCategories),
			),
		)

		index := 0
		total := 100
		for index < total {
			response, err := search(v.ID, index, 100)

			if err != nil {
				logger.Error(
					"Error searching",
					slog.Group("category",
						slog.String("slug", v.Slug),
						slog.String("id", v.ID),
					),
					slog.Any("error", err),
				)
				break
			}

			total = response.Sections.Null.Total
			index = index + 100

			for _, product := range response.Sections.Null.Products {
				publishFunc(
					mapPingoDoceProductToStoreProduct(store, product.Source),
				)
			}

			time.Sleep(time.Duration(appConfig.PolitenessDelay.PageMs) * time.Millisecond)
		}
		time.Sleep(time.Duration(appConfig.PolitenessDelay.CategoryMs) * time.Millisecond)
	}
}

func search(category string, offset, size int) (*definition.PingoDoceSearchResult, error) {
	urlFormat := "https://mercadao.pt/api/catalogues/6107d28d72939a003ff6bf51/products/search?mainCategoriesIds=[\"%s\"]&from=%d&size=%d&esPreference=0.7998979678255991"
	url := fmt.Sprintf(urlFormat, category, offset, size)

	res, err := request(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response definition.PingoDoceSearchResult
	err = json.NewDecoder(res.Body).Decode(&response)

	return &response, err
}

func fetchCategories() (*definition.PingoDoceCategories, error) {
	url := "https://mercadao.pt/api/catalogues/6107d28d72939a003ff6bf51/with-descendants"

	res, err := request(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response definition.PingoDoceCategories
	err = json.NewDecoder(res.Body).Decode(&response)

	return &response, err
}

func request(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("X-Version", "3.16.0")
	req.Header.Add("X-Name", "webapp")
	req.Header.Add("ngsw-bypass", "true")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("TE", "trailers")

	return client.Do(req)
}
