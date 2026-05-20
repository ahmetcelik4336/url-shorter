package dto

type PanelMenu struct {
	Key        string `json:"key"`
	Url        string `json:"content"`
	Title      string
	SubsExists bool
	Subs       []PanelMenu
	Order      int
}
