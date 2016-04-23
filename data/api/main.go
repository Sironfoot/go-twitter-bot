package api

import (
	"net/http"
	"strconv"

	"github.com/sironfoot/go-twitter-bot/data/db"
	"github.com/sironfoot/go-twitter-bot/data/models"
)

const ok = "OK"

type messageResponse struct {
	Message string `json:"message"`
}

type pagedResponse struct {
	messageResponse
	Page           int `json:"page"`
	RecordsPerPage int `json:"recordPerPage"`
	TotalRecords   int `json:"totalRecords"`
}

type createResponse struct {
	messageResponse
	Errors []models.ValidationError `json:"errors"`
	ID     *string                  `json:"id"`
}

type updateResponse struct {
	messageResponse
	Errors []models.ValidationError `json:"errors"`
}

func extractAndValidatePagingInfo(req *http.Request) (paging db.PagingInfo, errResponse *messageResponse) {
	qs := req.URL.Query()

	recordsPerPage, err := strconv.Atoi(qs.Get("recordsPerPage"))
	if err != nil {
		recordsPerPage = 20
	}

	page, err := strconv.Atoi(qs.Get("page"))
	if err != nil {
		page = 1
	}

	if recordsPerPage < 1 || recordsPerPage > 100 {
		errResponse = &messageResponse{
			Message: "'recordsPerPage' must be between 1 and 100",
		}
		return
	}

	if page < 1 || page > 100 {
		errResponse = &messageResponse{
			Message: "'page' must be between 1 and 100",
		}
		return
	}

	paging = db.PagingInfo{
		OrderBy:        db.UsersOrderByDateCreated,
		Asc:            false,
		Page:           page,
		RecordsPerPage: recordsPerPage,
	}

	return
}
