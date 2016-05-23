package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
)

// AppContext is an app specific context struct to hold per request variables
type AppContext struct {
	Settings Config
	Response interface{}
}

// Config represents a configuration settings for the app
type Config struct {
	Database    Database    `json:"database"`
	AppSettings AppSettings `json:"appSettings"`
}

// Database represents database configuration settings for the app
type Database struct {
	DriverName       string `json:"driverName"`
	ConnectionString string `json:"connectionString"`
}

// AppSettings represents general application settings for the app
type AppSettings struct {
	ServerAddress    string `json:"serverAddress"`
	EncryptionKey    string `json:"encryptionKey"`
	BCryptWorkFactor int    `json:"bcryptWorkFactor"`
}

// MessageResponse represents a standard JSON message response
type MessageResponse struct {
	Message string `json:"message"`
}

const ok = "OK"
const (
	defaultRecordsPerPage = 20
	maxPage               = 100
	maxRecordsPerPage     = 100
)

type pagedResponse struct {
	Message        string `json:"message"`
	Page           int    `json:"page"`
	RecordsPerPage int    `json:"recordPerPage"`
	TotalRecords   int    `json:"totalRecords"`
}

type createResponse struct {
	Message string                   `json:"message"`
	Errors  []models.ValidationError `json:"errors"`
	ID      *string                  `json:"id"`
}

type updateResponse struct {
	Message string                   `json:"message"`
	Errors  []models.ValidationError `json:"errors"`
}

func getPagingDefaults(orderBy string, asc bool, allowedOrderByColumns []string) PagingDefaults {
	return PagingDefaults{
		RecordsPerPage:        defaultRecordsPerPage,
		OrderBy:               orderBy,
		Asc:                   asc,
		MaxPage:               maxPage,
		MaxRecordsPerPage:     maxRecordsPerPage,
		AllowedOrderByColumns: allowedOrderByColumns,
	}
}

// PagingDefaults specifies default querysting paging arguments when they aren't present,
// and also restricts max 'page' and 'recordsPerPage' values as well as allowed orderBy column names.
type PagingDefaults struct {
	RecordsPerPage        int
	OrderBy               string
	Asc                   bool
	MaxPage               int
	MaxRecordsPerPage     int
	AllowedOrderByColumns []string
}

// ExtractAndValidatePagingInfo extracts paging information from a URL querystring and
// validates it for correctness and allowed ranges/column names.
func ExtractAndValidatePagingInfo(req *http.Request, defaults PagingDefaults) (db.PagingInfo, error) {
	var paging db.PagingInfo
	qs := req.URL.Query()

	// extract recordsPerPage
	recordsPerPage, err := strconv.Atoi(qs.Get("recordsPerPage"))
	if err != nil {
		recordsPerPage = defaults.RecordsPerPage
	}

	if recordsPerPage < 1 || recordsPerPage > defaults.MaxRecordsPerPage {
		return paging, fmt.Errorf("'recordsPerPage' must be between 1 and %d", defaults.MaxRecordsPerPage)
	}

	// extract page number
	page, err := strconv.Atoi(qs.Get("page"))
	if err != nil {
		page = 1
	}

	if page < 1 || page > defaults.MaxPage {
		return paging, fmt.Errorf("'page' must be between 1 and %d", defaults.MaxPage)
	}

	// extract orderBy column
	orderBy := strings.TrimSpace(qs.Get("orderBy"))
	if orderBy == "" {
		orderBy = defaults.OrderBy
	}

	allowed := false
	for _, allowedOrderBy := range defaults.AllowedOrderByColumns {
		if orderBy == allowedOrderBy {
			allowed = true
		}
	}

	if !allowed {
		return paging, fmt.Errorf("'%s' is not a valid orderBy column. Must be one of: %s", orderBy, strings.Join(defaults.AllowedOrderByColumns, ", "))
	}

	// extract sort order direction
	ascQs := strings.ToLower(qs.Get("asc"))
	asc := defaults.Asc
	if ascQs == "true" {
		asc = true
	} else if ascQs == "false" {
		asc = false
	}

	paging.OrderBy = orderBy
	paging.Asc = asc
	paging.Page = page
	paging.RecordsPerPage = recordsPerPage

	return paging, nil
}
