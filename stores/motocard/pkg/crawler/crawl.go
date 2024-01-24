package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"github.com/brunofjesus/pricetracker/stores/motocard/config"
	"github.com/brunofjesus/pricetracker/stores/motocard/pkg/definition"
	"log/slog"
	"net/http"
	"time"
)

func Crawl(logger *slog.Logger, politenessDelay int64, publishFunc func(dto.StoreProduct)) {
	appConfig := config.GetApplicationConfiguration()

	pages := 1

	for i := 1; i <= pages; i++ {
		response, err := search(
			appConfig.UserAgent,
			appConfig.Store.Country,
			appConfig.Store.Department,
			i,
		)
		if err != nil {
			logger.Error("could not fetch page. Aborting", slog.Any("error", err))
			break
		}

		pages := response.View.Paginator.Pages
		rawProducts := response.View.Results
		currency := response.Context.Locale.Currency.Code
		slog.Info("got page ",
			slog.Int("currentPage", i),
			slog.Int("totalPages", pages),
			slog.Int("rawProducts", len(rawProducts)))

		for _, product := range rawProducts {
			product := mapMotocardProductToStoreProduct(logger, appConfig.Store.Slug, product, currency)
			publishFunc(product)
		}

		time.Sleep(time.Millisecond * time.Duration(politenessDelay))
	}
}

func search(
	userAgent string,
	country string,
	department string,
	page int,
) (*definition.MotocardPageResponse, error) {

	url := fmt.Sprintf("https://www.motocard.com/view/results?url=/%s?p=%d&nidx", country, page)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Host", "www.motocard.com")
	req.Header.Add("Referer", fmt.Sprintf("https://www.motocard.com/%s/", country))
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("X-Requested-Department", department)
	req.Header.Add("X-Requested-Language", country)
	req.Header.Add("X-Requested-Mode", "clientside")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response definition.MotocardPageResponse
	err = json.NewDecoder(res.Body).Decode(&response)

	return &response, err
}
