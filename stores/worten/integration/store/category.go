package store

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/brunofjesus/pricetracker/stores/worten/config"
	"github.com/brunofjesus/pricetracker/stores/worten/definition/store"
	"github.com/brunofjesus/pricetracker/stores/worten/integration/sitemap"
)

func FindCategories(logger *slog.Logger) (map[string]string, error) {
	var result = make(map[string]string)

	appConfig := config.GetApplicationConfiguration()

	logger.Debug("loading sitemap", slog.String("url", appConfig.CategoriesSitemap))
	sitemap, err := sitemap.GetSitemapGZ(appConfig.CategoriesSitemap)
	if err != nil {
		return result, err
	}

	logger.Debug("loading known categories")
	knownCategoriesMap, err := loadKnownCategories()
	if err != nil {
		return result, err
	}

	for _, element := range sitemap.Elements {
		if categoryId, exists := knownCategoriesMap[element.Loc]; exists {
			logger.Debug("Known Category", slog.String("id", categoryId), slog.String("url", element.Loc))
			result[categoryId] = element.Loc
		} else {
			categoryId, err := solveUrl(element.Loc)
			if err != nil {
				logger.Error("Cannot solve category", slog.String("url", element.Loc), slog.Any("error", err))
				continue
			}

			logger.Debug("New Category", slog.String("id", categoryId), slog.String("url", element.Loc))
			result[categoryId] = element.Loc

			knownCategoriesMap[element.Loc] = categoryId
			if err = addCategory(categoryId, element.Loc); err != nil {
				logger.Error("Cannot persist in known categories", slog.Any("error", err))
			}

			time.Sleep(time.Millisecond * time.Duration(appConfig.PolitenessDelayMs))
		}

	}

	return result, nil
}

func solveUrl(categoryUrl string) (string, error) {
	requestBody := store.WortenSolveURLRequest{
		OperationName: "solveURL",
		Variables: store.WortenSolveURLRequestVariables{
			Debug:           false,
			URI:             categoryUrl,
			FetchFullEntity: false,
		},
		Query: store.WortenSolveURLRequestQuery,
	}

	client := &http.Client{}

	byteBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", errors.New("error creating payload")
	}

	req, err := http.NewRequest("POST", "https://www.worten.pt/_/api/graphql?wOperationName=solveURL", bytes.NewBuffer(byteBody))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.9")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Host", "www.worten.pt")
	req.Header.Add("Origin", "https://www.worten.pt")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("x-worten-tenant", "pt")
	req.Header.Add("Priority", "u=3, i")
	req.Header.Add("x-forwarded-proto", "https:")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var response store.WortenSolveURLResponse
	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}

	var categoryId = response.Data.SolveURL.Item.ID

	return categoryId, nil
}

func loadKnownCategories() (map[string]string, error) {
	var result = make(map[string]string)

	_, err := os.Stat("categories.csv")
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return result, err
	}

	file, err := os.Open("categories.csv")
	if err != nil {
		return result, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("%v", err)
		}

		if len(record) == 2 {
			result[record[0]] = record[1]
		} else {
			fmt.Printf("Wrong length, expected 2, got %d", len(record))
		}
	}

	return result, nil
}

func addCategory(id string, url string) error {
	file, err := os.OpenFile("categories.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	return w.Write([]string{url, id})
}
