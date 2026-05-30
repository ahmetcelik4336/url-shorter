package controllers

import (
	"encoding/base64"
	dto "shared/models"
	"shared/utils"

	"github.com/skip2/go-qrcode"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	usersCountRes, _ := utils.SendRequest[dto.UrlCountAnalysisResponse](nil, "user/GetUserCount", "GET", c.Ctx, "")
	c.Data["usersCount"] = usersCountRes.Count

	clicked, _ := utils.SendRequest[dto.ResponseClicked](nil, "user/GetPerDayClick", "GET", c.Ctx, "")
	c.Data["clicked"] = clicked.Count

	urltotal, _ := utils.SendRequest[dto.UrlCountAnalysisResponse](nil, "user/GetURLCountTotal", "GET", c.Ctx, "")
	c.Data["urltotal"] = urltotal.Count
	c.Layout = "inc/layout.html"
	c.TplName = "home/index.html"

}

func (c *MainController) UrlShortenHandler() {
	request := &dto.ExcelUrlRequest{}
	if err := c.ParseForm(request); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["message"] = "Hata"
		c.TplName = "home/url.html"
		return
	}
	resultUrl, _ := utils.SendRequest[*dto.GeneralResponse[*dto.UrlResponse]](request, "create/url", "POST", c.Ctx, "")
	c.Data["url"] = resultUrl.Data
	c.Data["status"] = resultUrl.Status
	c.Data["message"] = resultUrl.Message
	c.TplName = "home/url.html"
}

func (c *MainController) QrGenerateHandler() {
	request := &dto.ExcelUrlRequest{}
	if err := c.ParseForm(request); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["message"] = "Hata"
		c.TplName = "home/qr.html"
		return
	}
	resultUrl_, _ := utils.SendRequest[*dto.GeneralResponse[dto.UrlResponse]](request, "create/qr", "POST", c.Ctx, "")
	c.Data["status"] = false
	if resultUrl_.Data != nil {
		var png []byte
		png, err := qrcode.Encode(resultUrl_.Data.ResultUrl, qrcode.Highest, 250)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			resultUrl_.Message = "QR kod üretilemedi"
			c.TplName = "home/qr.html"
			return
		}
		qrBase64 := base64.StdEncoding.EncodeToString(png)
		c.Data["qr"] = qrBase64
		c.Data["status"] = true
	}

	c.Data["message"] = resultUrl_.Message
	c.TplName = "home/qr.html"
}
