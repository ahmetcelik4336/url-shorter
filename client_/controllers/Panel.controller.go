package controllers

import (
	"fmt"
	dto "shared/models"
	"shared/utils"

	beego "github.com/beego/beego/v2/server/web"
)

type PanelController struct {
	beego.Controller
}

func (c *PanelController) Panel() {
	req := &dto.UrlCountAnalysisRequest{Date: ""}
	res, err := utils.SendRequest[dto.UrlCountAnalysisResponse](req, "panel/GetURLCount", "POST", c.Ctx)

	if err != nil {
		fmt.Println("Hata oluştu:", err)
		return
	}

	fmt.Println(res.Count)
	c.TplName = "panel.html"
}
