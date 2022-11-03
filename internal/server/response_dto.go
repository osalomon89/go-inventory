package server

import "github.com/osalomon89/go-inventory/internal/domain"

// swagger:model ResponseInfo
type ResponseInfo struct {
	Status int         `json:"status"`
	Data   domain.Book `json:"data"`
}

// swagger:model ResponseAllInfo
type ResponseAllInfo struct {
	Status int           `json:"status"`
	Data   []domain.Book `json:"data"`
}

// swagger:model ResponseDeleteInfo
type ResponseDeleteInfo struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

// swagger:model ResponseError
type ResponseError struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}
