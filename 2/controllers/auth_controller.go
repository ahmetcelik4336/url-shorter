package controllers

import (
	"encoding/json"
	dto "shared/models"
	"shared/utils"
	"shared/validator"

	"2/services"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	beego.Controller
	Service services.UserService
}

// @Title Login
// @Description Kullanıcı giriş yapar ve JWT token döner
// @Param body body dto.LoginRequest true "Login bilgileri"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {string} string "invalid request"
// @router /user/login [post]
func (c *AuthController) Login() {

	var req dto.LoginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, "invalid request")
		return
	}

	if resp := validator.CheckStruct(req, "tr"); resp != nil {
		c.Data["json"] = dto.GeneralResponse[any]{
			Errors:  resp.Errors,
			Status:  false,
			Message: "Validation errors",
		}
		c.ServeJSON()
		return
	}

	token, err := c.Service.Login(req.Email, req.Password)
	if err != nil {
		c.CustomAbort(401, "email veya şifre bulunamadı")
		return
	}

	c.Data["json"] = token
	c.ServeJSON()
}

// @Title Register
// @Description Yeni kullanıcı kaydı oluşturur
// @Param body body dto.RegisterRequest true "Register bilgileri"
// @Success 200 {object} dto.UserResponse "başarılı kayıt"
// @Failure 400 {string} string "invalid request"
// @Failure 409 {object} dto.UserResponse "email exists"
// @Failure 401 {string} string "error"
// @router /user/register [post]
func (c *AuthController) Register() {

	var req dto.RegisterRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, "invalid request")
		return
	}

	status := c.Service.EmailExists(req.Email)

	if status {
		c.Data["json"] = dto.RegisterResponse{
			Status:  false,
			Message: "Email Kullanılıyor",
		}
		c.ServeJSON()
		return
	}

	register, err := c.Service.Register(req)
	if err != nil {
		c.CustomAbort(401, "error")
		return
	}

	c.Data["json"] = register
	c.ServeJSON()
}

// @Title Profile
// @Description Giriş yapan kullanıcının profil bilgisi
// @Success 200 {object} dto.UserResponse
// @Failure 500 {string} string "server error"
// @Security ApiKeyAuth
// @router /user/profile [get]
func (c *AuthController) Profile() {

	user, err := c.Service.FindUser(c.Ctx.Input.GetData("userID").(int))
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = user
	c.ServeJSON()
}

func (c *AuthController) ValidateToken() {
	token := utils.GetTokenFromHeader(c.Ctx.Input.Header("Authorization"))
	user, err := c.Service.ValidateToken(token)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = user
	c.ServeJSON()
}

func (c *AuthController) GetUserCount() {

	user, err := c.Service.GetUsersCount()
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = user
	c.ServeJSON()
}
