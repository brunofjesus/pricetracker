package store

type WortenCategoriesRequest struct {
	OperationName string                           `json:"operationName"`
	Variables     WortenCategoriesRequestVariables `json:"variables"`
	Query         string                           `json:"query"`
}

type WortenCategoriesRequestVariables struct {
	Debug           bool   `json:"debug"`
	URI             string `json:"uri"`
	FetchFullEntity bool   `json:"fetchFullEntity"`
}

var WortenCategoriesRequestQuery = "query solveURL($uri: String!, $debug: Boolean!, $params: SolveURLParams) {\n  solveURL(uri: $uri, params: $params) {\n    __typename\n    _id\n    code\n    success\n    message\n    contextId\n    layout {\n      __typename\n      _id\n      modules {\n        __typename\n        order\n        priority @include(if: $debug)\n        targetedBy\n        targets @include(if: $debug)\n        refs {\n          ...Ref\n          __typename\n        }\n      }\n    }\n  }\n}\n\nfragment Ref on ICMSRef {\n  __typename\n  _key\n  _type\n  valid\n  ... on CMSInternalLink {\n    id\n    url\n    __typename\n  }\n  ... on CMSExternalLink {\n    url\n    __typename\n  }\n}"

type WortenBrowseProductsRequest struct {
	OperationName string                               `json:"operationName"`
	Variables     WortenBrowseProductsRequestVariables `json:"variables"`
	Query         string                               `json:"query"`
}

type WortenBrowseProductsRequestVariables struct {
	Contexts    []string                                   `json:"contexts"`
	Params      WortenBrowseProductsRequestVariablesParams `json:"params"`
	HasVariants bool                                       `json:"hasVariants"`
}

type WortenBrowseProductsRequestVariablesParams struct {
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
	Filters    []any      `json:"filters"`
	Sort       WortenSort `json:"sort"`
	Collapse   bool       `json:"collapse"`
}

type WortenSort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

var WortenBrowseProductsRequestQuery = "query browseProducts($contexts: [String!]!, $params: SearchInput!, $hasVariants: Boolean!) {\n  browseProducts(contexts: $contexts, params: $params) {\n    ...Sr\n    hits {\n      totalProductVariants @include(if: $hasVariants)\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment Sr on SearchResults {\n  hits {\n    _id\n    product {\n      id\n      ean\n      sku\n      tags\n      name\n      brandName\n      description\n      type\n      url\n     image {\n        ...MA\n        __typename\n      }\n      defaultCategory {\n        url\n        __typename\n      }\n      __typename\n    }\n    totalOffers\n    data\n    campaigns {\n      data\n      id\n      __typename\n    }\n    winningOffer {\n      _id\n      channelId\n      offerId\n      sellerType\n      seller {\n        id\n        name\n        __typename\n      }\n      price {\n        ...Cur\n        __typename\n      }\n      isInStock\n      tags\n      __typename\n    }\n    secondOfferPrice {\n      ...Cur\n      __typename\n    }\n    __typename\n  }\n  hasNextPage\n  totalHits\n  filterErrors {\n    key\n    error\n    __typename\n  }\n  __typename\n}\n\nfragment Cur on CurrencyValue {\n  currency\n  value\n  __typename\n}\n\nfragment MA on ProductMediaAsset {\n  type\n  data\n  url\n  mimeType\n  __typename\n}\n"
