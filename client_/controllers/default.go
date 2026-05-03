package controllers

import (
	dto "shared/models"
	"shared/utils"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	usersCountRes, _ := utils.SendRequest[dto.UrlCountAnalysisResponse](nil, "user/GetUserCount", "GET", c.Ctx)
	c.Data["usersCount"] = usersCountRes.Count
	c.Layout = "inc/layout.html"
	c.TplName = "home/index.html"

}
