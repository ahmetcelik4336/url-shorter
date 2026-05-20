package controllers

import (
	"api/internal/services"
	"encoding/json"
	dto "shared/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type UrlController struct {
	beego.Controller
	Service services.UrlService
}

// @Title Create
// @Description Yeni bir kısa URL oluşturur.
// @Param   body    body    dto.CreateUrlRequest    true    "Uzun URL ve opsiyonel bilgiler"
// @Success 200 {object} dto.UrlResponse
// @Failure 400 invalid request
// @Failure 500 server error
// @Security ApiKeyAuth
// @router /panel/create [post]
func (c *UrlController) Create() {
	var req dto.CreateUrlRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, "invalid request")
		return
	}

	urlResult, err := c.Service.Create(c.Ctx.Input.GetData("userID").(int), req)
	if err != nil {
		c.Data["json"] = dto.GeneralResponse[any]{
			Status:  false,
			Message: err.Error(),
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = urlResult
	c.ServeJSON()
}

// @Title History
// @Description Giriş yapmış kullanıcının URL oluşturma geçmişini listeler.
// @Success 200 {array} dto.UrlResponse
// @Failure 500 server error
// @Security ApiKeyAuth
// @Router /panel/history [get]
func (c *UrlController) History() {
	urls, err := c.Service.History(c.Ctx.Input.GetData("userID").(int))
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}
	c.Data["json"] = urls
	c.ServeJSON()
}

func (c *UrlController) HistoryById() {
	id := c.Ctx.Input.Param(":id")
	num, err := strconv.Atoi(id)
	urls, err := c.Service.HistoryById(c.Ctx.Input.GetData("userID").(int), num)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}
	c.Data["json"] = urls
	c.ServeJSON()
}

// @Title Update
// @Description Url bilgilerini güncelle
// @Param   body    body    dto.UpdateUrlRequest    true    "Uzun URL ve opsiyonel bilgiler"
// @Success 200 {object} dto.UrlResponse
// @Failure 500 server error
// @Security ApiKeyAuth
// @Router /panel/update [put]
func (c *UrlController) Update() {
	var req dto.UpdateUrlRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, "invalid request")
		return
	}

	user, err := c.Service.Update(req)
	if err != nil {
		c.Data["json"] = dto.GeneralResponse[any]{
			Status:  false,
			Message: err.Error(),
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = user
	c.ServeJSON()
}

// Bulkcreate godoc
// @Title Toplu URL Oluşturma
// @Description Birden fazla uzun URL'yi aynı anda kısaltmak için kullanılır
// @Param   body    body    []dto.CreateUrlRequest    true    "Oluşturulacak URL listesi"
// @Success 201 {object}  dto.GeneralResponse
// @Failure 400 {string} string "Geçersiz JSON formatı"
// @Failure 500 {string} string "Urller oluşturulurken hata oluştu"
// @Security ApiKeyAuth
// @Router /panel/bulk-create [post]
func (c *UrlController) Bulkcreate() {
	var req []dto.CreateUrlRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, err.Error())
		return
	}

	// 2026-06-01T10:00:00Z

	user, err := c.Service.Bulkcreate(c.Ctx.Input.GetData("userID").(int), req)
	if err != nil {
		c.Data["json"] = dto.GeneralResponse[any]{
			Status:  false,
			Message: err.Error(),
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = user
	c.ServeJSON()
}
