package dto

type SettingResponse struct {
	Key     string `form:"key" validate:"required" json:"key"`
	Content string `form:"content" validate:"required" json:"content"`
}

type SeoDetail struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type GeneralSettings struct {
	AppName   string               `json:"appname"`
	Social    map[string]string    `json:"social"`
	Seo       map[string]SeoDetail `json:"seo"`
	Maintence int                  `json:"maintence"`
}
