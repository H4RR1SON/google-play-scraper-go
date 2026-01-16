package gplay

import (
	"net/http"
	"time"
)

type CallOptions struct {
	Throttle   int
	Headers    http.Header
	RetryCount int
	RetryWait  time.Duration
}

type AppOptions struct {
	CallOptions
	AppID   string
	Lang    string
	Country string
}

type ListOptions struct {
	CallOptions
	Collection Collection
	Category   Category
	Age        *Age
	Num        int
	Lang       string
	Country    string
	FullDetail bool
}

type SearchPrice string

const (
	SearchPriceAll  SearchPrice = "all"
	SearchPriceFree SearchPrice = "free"
	SearchPricePaid SearchPrice = "paid"
)

type SearchOptions struct {
	CallOptions
	Term       string
	Num        int
	Lang       string
	Country    string
	FullDetail bool
	Price      SearchPrice
}

type DeveloperOptions struct {
	CallOptions
	DevID      string
	Num        int
	Lang       string
	Country    string
	FullDetail bool
}

type SuggestOptions struct {
	CallOptions
	Term    string
	Lang    string
	Country string
}

type ReviewsOptions struct {
	CallOptions
	AppID               string
	Lang                string
	Country             string
	Sort                Sort
	Num                 int
	Paginate            bool
	NextPaginationToken *string
}

type SimilarOptions struct {
	CallOptions
	AppID      string
	Lang       string
	Country    string
	FullDetail bool
	Num        int
}

type PermissionsOptions struct {
	CallOptions
	AppID   string
	Lang    string
	Country string
	Short   bool
}

type DataSafetyOptions struct {
	CallOptions
	AppID string
	Lang  string
}

type CategoriesOptions struct {
	CallOptions
}
