package controllers

import (
	"log"
	"net/http"
	"shared/helpers"
	dto "shared/models"
	"shared/utils"
	"shared/validator"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Login() {

	flash := beego.ReadFromRequest(&c.Controller)
	if success, ok := flash.Data["success"]; ok {
		c.Data["SuccessMessage"] = success
	}

	if success, ok := flash.Data["errmessage"]; ok {
		c.Data["Message"] = success
	}

	id, err := Cpt.CreateCaptcha()
	if err != nil {
		log.Println("Captcha oluşturma hatası:", err)
	}

	c.Data["captchaId"] = id
	c.Data["IsShowHeader"] = ""
	c.Data["IsShowFooter"] = ""
	c.Layout = "inc/layout.html"
	c.TplName = "auth/login.html"
}

func (c *AuthController) Register() {
	c.Data["IsShowHeader"] = ""
	c.Data["IsShowFooter"] = ""
	c.Layout = "inc/layout.html"
	c.TplName = "auth/register.html"
}

func (c *AuthController) LoginHandler() {
	var req dto.LoginRequest

	if err := c.ParseForm(&req); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Form verisi işlenemedi")
		return
	}

	captchaId := c.GetString("captcha_id")
	userInput := c.GetString("captcha_value")

	langVal := c.Ctx.Input.GetData("lang")
	lang, ok := langVal.(string)
	if !ok || lang == "" {
		lang = "en"
	}

	// Doğrulama kontrolü
	if !Cpt.Verify(captchaId, userInput) {
		flash := beego.NewFlash()
		flash.Data["errmessage"] = "Doğrulama Hatalı!"
		flash.Store(&c.Controller)
		c.Redirect(helpers.Baseurl(lang, "auth/login"), 302)
		return
	}

	if resp := validator.CheckStruct(req, lang); resp != nil {
		c.Data["IsShowHeader"] = ""
		c.Data["IsShowFooter"] = ""
		c.Data["Errors"] = resp.Errors
		c.Data["Message"] = resp.Message
		c.Data["OldEmail"] = req.Email
		c.Layout = "inc/layout.html"
		c.TplName = "auth/login.html"
		return
	}

	apiResp, _ := utils.SendRequest[dto.LoginResponse](req, "user/login", "POST", c.Ctx, "")

	if apiResp.Token != "" {
		c.SetSession("token", apiResp.Token)
		c.Redirect(helpers.Baseurl(lang, "panel"), 302)
		return
	}

	c.Data["Message"] = "E-posta veya şifre hatalı"
	c.Data["IsShowHeader"] = ""
	c.Data["IsShowFooter"] = ""
	c.Data["OldEmail"] = req.Email
	c.Layout = "inc/layout.html"
	c.TplName = "auth/login.html"
}

func (c *AuthController) RegisterHandler() {
	var req dto.RegisterRequest

	if err := c.ParseForm(&req); err != nil {
		c.CustomAbort(http.StatusBadRequest, "Form verisi işlenemedi")
		return
	}
	langVal := c.Ctx.Input.GetData("lang")
	lang, ok := langVal.(string)
	if !ok || lang == "" {
		lang = "en"
	}

	if resp := validator.CheckStruct(req, lang); resp != nil {
		c.Data["IsShowHeader"] = ""
		c.Data["IsShowFooter"] = ""
		c.Data["Errors"] = resp.Errors
		c.Data["Message"] = resp.Message
		c.Data["OldEmail"] = req.Email
		c.Layout = "inc/layout.html"
		c.TplName = "auth/register.html"
		return
	}

	apiResp, _ := utils.SendRequest[dto.RegisterResponse](req, "user/register", "POST", c.Ctx, "")

	if apiResp.Status == true {
		flash := beego.NewFlash()
		flash.Data["success"] = "Kayıt başarılı"
		flash.Store(&c.Controller)
		c.Redirect(helpers.Baseurl(lang, "auth/login"), 302)
		return
	}

	c.Data["Message"] = apiResp.Message
	c.Data["IsShowHeader"] = ""
	c.Data["IsShowFooter"] = ""
	c.Layout = "inc/layout.html"
	c.TplName = "auth/register.html"
}
func (c *AuthController) LogOut() {
	c.DelSession("token")
	c.Redirect(helpers.Baseurl(c.Data["slug"], "auth/login"), 302)
}
