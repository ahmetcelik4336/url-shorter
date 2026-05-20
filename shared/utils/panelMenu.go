package utils

import (
	"shared/helpers"
	dto "shared/models"

	"github.com/beego/i18n"
)

func GetPanelMenus(lang string, data dto.GeneralSettings) []dto.PanelMenu {
	return []dto.PanelMenu{
		{
			Key:        "home",
			Title:      i18n.Tr(lang, "home"),
			Url:        helpers.Baseurl(lang, "panel"),
			SubsExists: false,
			Order:      1,
			Subs: []dto.PanelMenu{
				{
					Key:   "home2",
					Title: "Home 1",
				},
			},
		},
		{
			Key:        "profile",
			Title:      i18n.Tr(lang, "profile"),
			Url:        helpers.Baseurl(lang, "panel/profile"),
			SubsExists: false,
			Order:      2,
		},
		{
			Key:        "urls",
			Title:      i18n.Tr(lang, "urls"),
			Url:        helpers.Baseurl(lang, "panel/urls"),
			SubsExists: false,
			Order:      3,
		},
		{
			Key:        "logout",
			Title:      i18n.Tr(lang, "logout"),
			Url:        helpers.Baseurl(lang, "panel/logout"),
			SubsExists: false,
			Order:      4,
		},
	}
}

/*
func GetPanelActiveMenu(active, lang string, data dto.GeneralSettings) dto.PanelMenu {
	menus := GetPanelMenus(lang, data)
	return menus[active]
}*/
