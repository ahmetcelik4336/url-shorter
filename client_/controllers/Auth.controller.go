package controllers

import (
	"net/http"
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

	if resp := validator.CheckStruct(req); resp != nil {
		c.Data["IsShowHeader"] = ""
		c.Data["IsShowFooter"] = ""
		c.Data["Errors"] = resp.Errors
		c.Data["Message"] = resp.Message
		c.Data["OldEmail"] = req.Email
		c.Layout = "inc/layout.html"
		c.TplName = "auth/login.html"
		return
	}

	apiResp, _ := utils.SendRequest[dto.LoginResponse](req, "user/login", "POST", c.Ctx)

	if apiResp.Token != "" {
		c.SetSession("token", apiResp.Token)
		c.Redirect("/", 302)
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

	if resp := validator.CheckStruct(req); resp != nil {
		c.Data["IsShowHeader"] = ""
		c.Data["IsShowFooter"] = ""
		c.Data["Errors"] = resp.Errors
		c.Data["Message"] = resp.Message
		c.Data["OldEmail"] = req.Email
		c.Layout = "inc/layout.html"
		c.TplName = "auth/register.html"
		return
	}

	apiResp, _ := utils.SendRequest[dto.RegisterResponse](req, "user/register", "POST", c.Ctx)

	if apiResp.Status == true {
		flash := beego.NewFlash()
		flash.Data["success"] = "Kayıt başarılı"
		flash.Store(&c.Controller)
		c.Redirect("/auth/login", 302)
		return
	}

	c.Data["Message"] = apiResp.Message
	c.Data["IsShowHeader"] = ""
	c.Data["IsShowFooter"] = ""
	c.Layout = "inc/layout.html"
	c.TplName = "auth/register.html"
}
