package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/brunofjesus/pricetracker/stores/worten/definition/store"
)

func FindCategories() (map[string]string, error) {
	requestBody := []store.WortenCategoriesRequest{
		{
			OperationName: "solveURL",
			Variables: store.WortenCategoriesRequestVariables{
				Debug:           true,
				URI:             "/diretorio-de-categorias",
				FetchFullEntity: false,
			},
			Query: store.WortenCategoriesRequestQuery,
		},
	}

	client := &http.Client{}

	byteBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.New("error creating payload")
	}

	req, err := http.NewRequest("POST", "https://www.worten.pt/_/api/graphql?wOperationName=solveURL", bytes.NewBuffer(byteBody))

	if err != nil {
		fmt.Println(err)
		return nil, err
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
		return nil, err
	}
	defer res.Body.Close()

	var response []store.WortenCategoriesResponse
	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, errors.New("cannot fetch categories: empty response")
	}

	var modules = response[0].Data.SolveURL.Layout.Modules

	var result map[string]string = make(map[string]string)

	for _, module := range modules {
		if module.Targets == "pt-diretorio-categorias-category-links" {
			for _, ref := range module.Refs {
				if len(ref.URL) > 1 {
					result[ref.ID] = ref.URL
				}
			}
		}
	}

	if len(result) == 0 {
		return result, errors.New("cannot find categories")
	}

	return result, nil
}
