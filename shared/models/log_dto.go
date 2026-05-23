package dto

type LogRequest struct {
	Device  string `json:"device"`
	Ip      string `json:"ip"`
	Referer string `json:"referer"`
	Type    string `json:"type"`
}

type LogResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Code    int            `json:"code"`
	Data    map[string]any `json:"data"`
}
