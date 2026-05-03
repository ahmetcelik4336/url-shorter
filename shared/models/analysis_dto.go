package dto

type UsageAnalysisRequest struct {
	Date string `json:"date"`
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
