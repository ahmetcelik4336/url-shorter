package controllers

import (
	dto "shared/models"
	"shared/utils"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	usersCountRes, _ := utils.SendRequest[dto.UrlCountAnalysisResponse](nil, "user/GetUserCount", "GET", c.Ctx, "")
	c.Data["usersCount"] = usersCountRes.Count

	clicked, _ := utils.SendRequest[dto.ResponseClicked](nil, "user/GetPerDayClick", "GET", c.Ctx, "")
	c.Data["clicked"] = clicked.Count

	urltotal, _ := utils.SendRequest[dto.UrlCountAnalysisResponse](nil, "user/GetURLCountTotal", "GET", c.Ctx, "")
	c.Data["urltotal"] = urltotal.Count
	c.Layout = "inc/layout.html"
	c.TplName = "home/index.html"

}
