package dto

import "time"

type UsageAnalysisRequest struct {
	Start time.Time `json:"start,omitempty" form:"start"`
	End   time.Time `json:"end,omitempty" form:"end"`
}

type UrlCountAnalysisRequest struct {
	Date string `json:"date"`
}

type UrlCountAnalysisResponse struct {
	Count int `json:"count"`
}

type UsageAnalysisResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ResponseClicked struct {
	Status bool `json:"status"`
	Count  int  `json:"count"`
}
