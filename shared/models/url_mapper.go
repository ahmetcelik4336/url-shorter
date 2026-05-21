package dto

import (
	"api/ent"
	"os"
)

func ToUrlResponse(p *ent.Url) *UrlResponse {
	var IsEncrypted string = "Şifresiz"
	if p.IsEncrypted {
		IsEncrypted = "Şifreli"
	}

	var aliass string

	if p.Alias == "" {
		aliass = ""
	} else {
		aliass = p.Alias
	}

	// ExpirationDate için nil kontrolü yapıyoruz
	var formattedExpiration string
	if !p.ExpirationDate.IsZero() {
		formattedExpiration = p.ExpirationDate.Format("2.1.2006 15:04")
	} else {
		formattedExpiration = "Süresiz" // Veya "" (boş string) yapabilirsiniz
	}
	var base string
	var base2 string
	if p.Alias == "" {
		base = p.ShortCode
	} else {
		base = p.Alias
	}

	if p.IsEncrypted {
		base2 = "/" + p.Password
	} else {
		base2 = ""
	}

	siteurl := os.Getenv("SITEURL") + "url/"

	response := &UrlResponse{
		Id:              p.ID,
		ShortCode:       p.ShortCode,
		LongUrl:         p.LongURL,
		CreateAt:        p.CreatedAt.Format("2.1.2006 15:04"),
		Alias:           aliass,
		IsEncrypted:     p.IsEncrypted,
		IsEncryptedText: IsEncrypted,
		// Çökmeyi önleyen güvenli değişkeni buraya veriyoruz:
		ExpirationDate:          formattedExpiration,
		ExpirationDateFormatted: p.ExpirationDate.Format("2006-01-02 15:04"),
		ResultUrl:               siteurl + base + base2,
		Password:                p.Password,
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
