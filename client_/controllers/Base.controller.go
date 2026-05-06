package controllers

import (
	"fmt"

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
		lang = beego.AppConfig.DefaultString("defaultlang", "en")
	}

	// Verileri set et
	c.Data["lang"] = lang
	c.Data["slug"] = lang
	c.Data["IsShowFooter"] = "show"
	c.Data["IsShowHeader"] = "show"

	fmt.Println("Seçilen Dil:", lang)
}
