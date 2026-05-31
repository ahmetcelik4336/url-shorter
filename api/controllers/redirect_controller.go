package controllers

import (
	"api/internal/services"
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
	password := c.Ctx.Input.Param(":password")
	response, err := c.Service.Redirect(shortcode, password)
	if err != nil || response.Code != 200 {
		c.Data["json"] = dto.GeneralResponse[any]{
			Message: response.Message,
			Status:  false,
		}
		c.ServeJSON()
		return
	}

	if response.Code == 200 {
		var req dto.LogRequest
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
			c.Data["json"] = dto.GeneralResponse[any]{
				Message: "Hata oluştu!",
			}
			return
		}
		url, err1 := c.Service.GetById(response.Data.Id)
		userid := 0
		if err1 == nil && url != nil && url.Edges.User != nil {
			userid = url.Edges.User.ID
		}
		_, err := c.LogService.Create(response.Data.Id, userid, req)

		if err == nil {
			c.Data["json"] = response
		} else {
			c.Data["json"] = dto.GeneralResponse[any]{
				Message: "Hata oluştu!",
			}
		}

		c.ServeJSON()
		//c.Ctx.Redirect(http.StatusFound, response.Data.LongUrl)
		return
	}

	c.CustomAbort(response.Code, response.Message)
}
