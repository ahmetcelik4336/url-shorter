package controllers

import (
	"api/internal/services"
	"encoding/json"
	dto "shared/models"

	beego "github.com/beego/beego/v2/server/web"
)

type AnalysisController struct {
	beego.Controller
	Service    services.AnalysisService
	LogService services.LogService
}

func (c *AnalysisController) GetURLCountTotal() {

	token, err := c.Service.GetURLCountTotal()
	if err != nil {
		c.CustomAbort(401, err.Error())
		return
	}

	c.Data["json"] = token
	c.ServeJSON()
}

// @Title Login
// @Description Kullanıcı giriş yapar ve JWT token döner
// @Param body body dto.LoginRequest true "Login bilgileri"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {string} string "invalid request"
// @router /user/login [post]
func (c *AnalysisController) GetURLCount() {

	var req dto.UrlCountAnalysisRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, "invalid request")
		return
	}

	userID := c.Ctx.Input.GetData("userID").(int)

	token, err := c.Service.GetUserURLCount(userID, req)
	if err != nil {
		c.CustomAbort(401, err.Error())
		return
	}

	c.Data["json"] = token
	c.ServeJSON()
}

// @Title Login
// @Description Kullanıcı giriş yapar ve JWT token döner
// @Param body body dto.LoginRequest true "Login bilgileri"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {string} string "invalid request"
// @router /user/login [post]
func (c *AnalysisController) GetURLStats() {

	var req dto.UsageAnalysisRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {

	}
	//fmt.Println(c.Ctx.Input.GetData("userID").(int))

	userID := c.Ctx.Input.GetData("userID").(int)
	token, err := c.Service.GetURLStats(userID, req)
	if err != nil {
		c.CustomAbort(401, err.Error())
		return
	}

	c.Data["json"] = token
	c.ServeJSON()
}

func (c *AnalysisController) GetPerDayClick() {
	token, err := c.LogService.GetPerDayClick()
	if err != nil {
		c.CustomAbort(401, err.Error())
		return
	}

	c.Data["json"] = token
	c.ServeJSON()
}

func (c *AnalysisController) TotalReading() {
	userID := c.Ctx.Input.GetData("userID").(int)
	totalRead, err := c.LogService.TotalReading(userID)
	if err != nil {
		c.CustomAbort(401, err.Error())
		return
	}
	c.Data["json"] = &totalRead
	c.ServeJSON()
}

func (c *AnalysisController) UrlTrackAnalysis() {
	var req dto.UsageAnalysisRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {

	}
	userID := c.Ctx.Input.GetData("userID").(int)
	totalRead, err := c.LogService.GetUrlTrackAnalysis(userID, &req)
	if err != nil {
		c.CustomAbort(401, err.Error())
		return
	}
	c.Data["json"] = totalRead
	c.ServeJSON()
}

func (c *AnalysisController) LogDatatable() {

	var req dto.DataTableRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {

	}
	userid := c.Ctx.Input.GetData("userID").(int)
	response, err := c.LogService.LogDatatable(req.Draw, req.Start, req.Length, userid, req.SearchValue, req.OrderColumnIdx, req.OrderDir, req.StartDate, req.EndDate)
	if err != nil {
		c.CustomAbort(401, err.Error())
		return
	}
	c.Data["json"] = response
	c.ServeJSON()
}
