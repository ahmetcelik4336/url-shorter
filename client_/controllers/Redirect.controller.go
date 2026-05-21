package controllers

import (
	"fmt"
	"net/http"

	// Takma ad (alias) vermeden direkt paket yollarını yazıyoruz
	dto "shared/models"
	"shared/utils"

	"github.com/mssola/user_agent"
)

type RedirectController struct {
	BaseController
}

func (c *RedirectController) Redirect() {
	shortcode := c.Ctx.Input.Param(":shortcode")
	password := c.Ctx.Input.Param(":password")

	// 1. Cihaz bilgisini al (Önceki adımımız)
	uaString := c.Ctx.Input.UserAgent()
	ua := user_agent.New(uaString)
	deviceType := "Desktop"
	if ua.Mobile() {
		deviceType = "Mobile"
	}
	fullDeviceName := fmt.Sprintf("%s (%s)", deviceType, ua.OSInfo().Name)

	// 2. IP Adresini ve Referer bilgisini al
	ipAddress := c.Ctx.Input.IP()
	referer := c.Ctx.Input.Referer()

	// 3. Request struct'ına atamaları yap
	var request dto.LogRequest
	request.Device = fullDeviceName
	request.Ip = ipAddress
	request.Referer = referer // Kullanıcı buraya nereden geldi? (Örn: twitter.com, google.com)

	// İsteği gönder
	status, _ := utils.SendRequest[dto.GeneralResponse[dto.UrlResponse]](request, "url/"+shortcode+"/"+password, "POST", c.Ctx, "")
	if status.Code == 200 {
		c.Ctx.Redirect(http.StatusFound, status.Data.LongUrl)
	} else {
		c.Data["json"] = status
		c.ServeJSON()
	}
}
