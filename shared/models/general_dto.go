package dto

type GeneralResponse[T any] struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    *T                `json:"data,omitempty"`
	Code    int               `json:"code,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}
