package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brunofjesus/pricetracker/stores/pingodoce/definition"
)

func Crawl(store definition.Store, publishFunc func(definition.StoreProduct)) {
	index := 0
	total := 100
	for index < total {
		response, err := search(index, 100)

		if err != nil {
			log.Default().Print(err)
			break
		}

		total = response.Sections.Null.Total
		index = index + 100

		for _, product := range response.Sections.Null.Products {
			publishFunc(
				mapPingoDoceProductToStoreProduct(store, product.Source),
			)
		}

		time.Sleep(1 * time.Second) // be polite
	}
}

func search(offset, size int) (*definition.PingoDoceSearchResult, error) {

	urlFormat := "https://mercadao.pt/api/catalogues/6107d28d72939a003ff6bf51/products/search?from=%d&size=%d&esPreference=0.7998979678255991"
	url := fmt.Sprintf(urlFormat, offset, size)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response definition.PingoDoceSearchResult
	err = json.NewDecoder(res.Body).Decode(&response)

	return &response, err
}
