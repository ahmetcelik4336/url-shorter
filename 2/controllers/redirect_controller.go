package controllers

import (
	"2/services"
	"encoding/json"
	dto "shared/models"

	beego "github.com/beego/beego/v2/server/web"
)

type RedirectController struct {
	beego.Controller
	Service    services.UrlService
	LogService services.LogService
}

// @Title Redirect
// @Description Yeni bir kısa URL oluşturur.
// @Success      302        {string}  string  "Hedef URL'ye yönlendirilir"
// @Failure 400 invalid request
// @Failure 409 {409} dto.GeneralResponse
// @Failure 500 server error
// @router /url/:shortcode [post]
func (c *RedirectController) Redirect() {

	shortcode := c.Ctx.Input.Param(":shortcode")
	password := c.Ctx.Input.Header("password")
	response, err := c.Service.Redirect(shortcode, password)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	if response.Code == 200 {
		var req dto.LogRequest
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
			c.CustomAbort(400, err.Error())
			return
		}
		c.LogService.Create(response.Data.Id, req)
		c.Data["json"] = response.Data
		c.ServeJSON()
		//c.Ctx.Redirect(http.StatusFound, response.Data.LongUrl)
		return
	}

	c.CustomAbort(response.Code, response.Message)
}
