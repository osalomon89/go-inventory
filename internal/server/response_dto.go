package server

import "github.com/osalomon89/go-inventory/internal/domain"

// swagger:model ResponseInfo
type ResponseInfo struct {
	// The Response Status code
	// example: 200
	Status int `json:"status"`
	// The Response Data
	// example: 200
	Data domain.Book `json:"data"`
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
	// The Response Status code
	// example: 400
	Status int `json:"status"`
	// The Response Data
	// example: error getting book
	Data string `json:"data"`
}
