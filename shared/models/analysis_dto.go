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

type UrlTrackAnalysisResponse struct {
	Count int    `json:"count"`
	Type  string `json:"type"`
	Title string `json:"title"`
}
type GetLastReadingResponse struct {
	LastReading string
}
type UrlTrackAnalysisResponseBatch struct {
	Analysis    []*UrlTrackAnalysisResponse
	LastReading *GetLastReadingResponse
}

type UsageAnalysisResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ResponseClicked struct {
	Status bool `json:"status"`
	Count  int  `json:"count"`
}
