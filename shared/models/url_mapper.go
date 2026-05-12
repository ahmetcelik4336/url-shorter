package dto

import (
	"api/ent"
)

func ToUrlResponse(p *ent.Url) *UrlResponse {

	//loc, _ := time.LoadLocation("Europe/Istanbul")

	response := &UrlResponse{
		Id:        p.ID,
		ShortCode: p.ShortCode,
		LongUrl:   p.LongURL,
		CreateAt:  p.CreatedAt.Format("2.1.2006 15:04"),
	}

	return response
}

func ToUrlResponseList(posts []*ent.Url) []*UrlResponse {

	var list []*UrlResponse

	for _, p := range posts {
		list = append(list, ToUrlResponse(p))
	}

	return list
}
