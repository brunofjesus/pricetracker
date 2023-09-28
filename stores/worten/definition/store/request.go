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
