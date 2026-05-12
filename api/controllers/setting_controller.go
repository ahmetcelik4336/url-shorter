package controllers

import (
	"api/internal/services"
	"encoding/json"
	dto "shared/models"

	beego "github.com/beego/beego/v2/server/web"
)

type SettingController struct {
	beego.Controller
	Service services.SettingService
}

// @Title Redirect
// @Description Yeni bir kısa URL oluşturur.
// @Success      302        {string}  string  "Hedef URL'ye yönlendirilir"
// @Failure 400 invalid request
// @Failure 409 {409} dto.GeneralResponse
// @Failure 500 server error
// @router /url/:shortcode [post]
func (c *SettingController) Save() {
	key := c.Ctx.Input.Param(":type")

	var content any
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &content); err != nil {
		json := &dto.GeneralResponse[any]{
			Status: false,
		}
		c.Data["json"] = json
		c.ServeJSON()
		return
	}

	setting, err := c.Service.Save(key, content)
	if setting != nil {
		json := &dto.GeneralResponse[any]{
			Status: true,
		}
		c.Data["json"] = json
		c.ServeJSON()
		return
	}
	res := make(map[string]string, 0)
	res["Err"] = err.Error()
	json := &dto.GeneralResponse[any]{
		Status: false,
		Errors: res,
	}
	c.Data["json"] = json
	c.ServeJSON()
}

func (c *SettingController) GetGeneralSettings() {

	setting, err := c.Service.GetGeneralSettings()
	if err != nil {
		json := &dto.GeneralResponse[any]{
			Status: false,
		}
		c.Data["json"] = json
		c.ServeJSON()
		return
	}

	c.Data["json"] = setting
	c.ServeJSON()
}
