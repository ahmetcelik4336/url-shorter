package controllers

import (
	"encoding/json"
	"log"
	"shared/helpers"
	dto "shared/models"
	"shared/utils"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type PanelController struct {
	BaseController
}

func (c *PanelController) Panel() {
	general, _ := c.Data["conf"].(dto.GeneralSettings)
	c.Data["menus"] = utils.GetPanelMenus(c.Lang, general)
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	req := dto.UsageAnalysisRequest{}
	if err := c.ParseForm(&req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Geçersiz tarih formatı"}
		c.ServeJSON()
		return
	}

	usage, _ := utils.SendRequest[[]dto.UsageAnalysisResponse](req, "panel/GetURLStats", "POST", c.Ctx, tokenn)
	jsons, _ := json.Marshal(usage)
	c.Data["usage"] = string(jsons)
	if !req.Start.IsZero() {
		c.Data["start"] = req.Start.Format("2006-01-02")
	}
	if !req.End.IsZero() {
		c.Data["end"] = req.End.Format("2006-01-02")
	}

	c.Data["headerActive"] = "active"
	c.Data["active"] = "home"
	c.Layout = "inc/layout.html"
	c.TplName = "panel/index.html"
}

func (c *PanelController) Url() {
	general, _ := c.Data["conf"].(dto.GeneralSettings)
	c.Data["menus"] = utils.GetPanelMenus(c.Lang, general)
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)

	c.Data["headerActive"] = "active"
	c.Data["active"] = "urls"
	urls, err := utils.SendRequest[[]dto.UrlResponse](nil, "panel/history", "GET", c.Ctx, tokenn)
	log.Println("urls", err)
	c.Data["urls"] = urls
	c.Layout = "inc/layout.html"
	c.TplName = "panel/url/index.html"
}

func (c *PanelController) UrlAdd() {
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	id := c.Ctx.Input.Param(":id")
	c.Data["action"] = "panel/urls/save"
	if id != "" {
		urls, _ := utils.SendRequest[[]dto.UrlResponse](nil, "panel/history/"+id, "GET", c.Ctx, tokenn)
		c.Data["urls"] = urls
		c.Data["action"] = "panel/urls/save/" + id
	}

	c.TplName = "panel/url/add.html"
}

func (c *PanelController) UrlSaveHandler() {
	id := c.Ctx.Input.Param(":id")
	idd, err := strconv.Atoi(id)
	if idd > 0 && err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.Body([]byte("Invalid ID format"))
		return
	}
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	var req any
	var action string
	if id == "" {
		req = &dto.CreateUrlRequest{}
		action = "create"
	} else {
		req = &dto.UpdateUrlRequest{
			ID: idd,
		}
		action = "update"
	}

	if err := c.ParseForm(req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	log.Println("post", req)

	status, err := utils.SendRequest[dto.GeneralResponse[any]](req, "panel/"+action, "POST", c.Ctx, tokenn)

	log.Println("errs", err)

	flash := beego.NewFlash()
	if status.Status {
		flash.Data["success"] = status.Message
	} else {
		flash.Data["err"] = status.Message
	}

	flash.Store(&c.Controller)
	c.Redirect(helpers.Baseurl(c.Data["lang"], "panel/urls"), 302)
}

func (c *PanelController) UrlDelete() {

}
