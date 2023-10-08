package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/brunofjesus/pricetracker/stores/worten/definition/store"
)

func FindProducts(pageNumber int, categoryId, slug string, callback func(store.WortenProductHit) error) (bool, error) {
	requestBody := []store.WortenBrowseProductsRequest{
		{
			OperationName: "browseProducts",
			Variables: store.WortenBrowseProductsRequestVariables{
				Contexts: []string{categoryId},
				Params: store.WortenBrowseProductsRequestVariablesParams{
					PageNumber: pageNumber,
					PageSize:   48,
					Filters: []store.WortenBrowseProductsRequestVariablesParamsFilter{
						{
							Key:     "seller_name",
							Virtual: false,
							Value:   []string{"Worten"},
						},
					},
					Sort: store.WortenSort{
						Field: "rank1",
						Order: "ASC",
					},
					Collapse: false,
				},
				HasVariants: false,
			},
			Query: store.WortenBrowseProductsRequestQuery,
		},
	}

	client := &http.Client{}

	byteBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, errors.New("error creating payload")
	}

	req, err := http.NewRequest("POST", "https://www.worten.pt/_/api/graphql?wOperationName=browseProducts", bytes.NewBuffer(byteBody))

	if err != nil {
		return false, err
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
		return false, err
	}
	defer res.Body.Close()

	var response []store.WortenBrowseProductResponse
	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		return false, err
	}

	if len(response) == 0 {
		return false, errors.New("cannot fetch products: empty response")
	}

	var hits = response[0].Data.BrowseProducts.Hits

	for _, hit := range hits {
		callback(hit)
	}
	return response[0].Data.BrowseProducts.HasNextPage, nil
}
