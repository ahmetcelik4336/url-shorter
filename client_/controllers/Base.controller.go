package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	c.Data["IsShowFooter"] = "show"
	c.Data["IsShowHeader"] = "show"
}
