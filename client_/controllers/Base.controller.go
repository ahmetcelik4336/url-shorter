package controllers

import (
	"os"
	dto "shared/models"
	"shared/utils"

	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	lang := c.Ctx.Input.GetData("lang")

	if lang == nil {
		lang = c.Ctx.Input.Param(":lang")
	}

	if lang == "" || lang == nil {
		lang = os.Getenv("DEFAULTLANG")
	}

	c.Data["lang"] = lang
	c.Data["slug"] = lang
	c.Data["IsShowFooter"] = "show"
	c.Data["IsShowHeader"] = "show"
	generalSetting, _ := utils.SendRequest[*dto.GeneralSettings](nil, "setting/get/general", "GET", c.Ctx)
	c.Data["conf"] = generalSetting
}
