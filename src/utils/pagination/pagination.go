package pagination

import (
	"database/sql"
	"fmt"
	"github.com/development-raul/footy-predictor/src/clients/mysql/footy_db"
	"github.com/development-raul/footy-predictor/src/utils/constants"
	"math"
	"reflect"
)

type PaginatedResponse struct {
	From        int64       `json:"from"`
	Data        interface{} `json:"data"`
	CurrentPage int64       `json:"current_page"`
	LastPage    int64       `json:"last_page"`
	PerPage     int64       `json:"per_page"`
	To          int64       `json:"to"`
	Total       int64       `json:"total"`
}

func GeneratePaginationQuery(page, perPage int64) string {
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = constants.DefaultPerPage
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", perPage, (page-1)*perPage)
}

func GeneratePaginatedResponse(data interface{}, page, perPage, total int64) PaginatedResponse {
	if page <= 1 {
		page = 1
	}

	if perPage == 0 {
		perPage = constants.DefaultPerPage
	}

	var lastPage int64
	var from int64
	var to int64
	if float64(total) > 0 && reflect.Slice == reflect.TypeOf(data).Kind() {
		from = (perPage * (page - 1)) + 1
		to = from + int64(reflect.ValueOf(data).Len()) - 1
		lastPage = int64(math.Ceil(float64(total) / float64(perPage)))
	} else {
		from = 0
		to = 0
		lastPage = 1
	}

	return PaginatedResponse{
		From:        from,
		Data:        data,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     perPage,
		To:          to,
		Total:       total,
	}
}

func GeneratePaginationSort(defaultSort string, sortBy string, sortOrder string) string {
	if sortBy == "" && sortOrder == "" {
		return defaultSort
	}
	return fmt.Sprintf("%s %s", sortBy, sortOrder)
}

func GetTableTotalRowsArgs(query string, params ...interface{}) (int64, error) {
	var total int64
	err := footy_db.Client.Get(&total, query, params...)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return total, nil
}
