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

type DataTableResponse struct {
	Draw            int         `json:"draw"`            // Frontend'den gelen draw değeri aynen dönmeli
	RecordsTotal    int64       `json:"recordsTotal"`    // Filtresiz toplam kayıt sayısı
	RecordsFiltered int64       `json:"recordsFiltered"` // Filtrelenmiş toplam kayıt sayısı
	Data            interface{} `json:"data"`            // Tabloda listelenecek veri dizisi
}

type DataTableRowsResponse struct {
	Id         int    `json:"id"`         // Frontend'den gelen draw değeri aynen dönmeli
	Code       string `json:"code"`       // Filtresiz toplam kayıt sayısı
	Created_at string `json:"created_at"` // Filtrelenmiş toplam kayıt sayısı
	Device     string `json:"device"`     // Tabloda listelenecek veri dizisi
	Ip         string `json:"ip"`
	Type       string `json:"type"`
	Reading_at string `json:"reading_at"`
}

type DataTableRequest struct {
	Draw           int       `json:"draw"`
	Start          int       `json:"start"`
	Length         int       `json:"length"`
	SearchValue    string    `json:"search[value]"`    // Karşı API'nin beklediği JSON key'leri
	OrderColumnIdx string    `json:"order[0][column]"` // Karşı API form mu JSON mı bekliyor buna göre isimlendirilmeli
	OrderDir       string    `json:"order[0][dir]"`
	StartDate      time.Time `json:"startdate"`
	EndDate        time.Time `json:"enddate"`
}
