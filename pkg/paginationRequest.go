package pkg

import (
	"net/http"
	"strconv"
)

type PaginationRequest struct {
	Limit int `json:"limit" validate:"required,gte=1,lte=100"`

	Page int `json:"page" validate:"required,gte=1"`

	SortKey string `json:"sort_key" validate:"omitempty,alphanum"`

	Sort string `json:"sort" validate:"omitempty,oneof=asc desc"`
}

func NewPaginationRequest(r *http.Request) *PaginationRequest {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	page, _ := strconv.Atoi(query.Get("page"))
	sortKey := query.Get("sortkey")
	sort := query.Get("sort")

	pagReq := PaginationRequest{
		Limit:   limit,
		Page:    page,
		SortKey: sortKey,
		Sort:    sort,
	}
	return &pagReq
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

type PaginationResponse[T any] struct {
	Data []T `json:"data"`

	Size int64 `json:"length"`
}
