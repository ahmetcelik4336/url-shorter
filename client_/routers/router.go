package routers

import (
	"client_/controllers"
	"client_/middleware"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/auth/login", &controllers.AuthController{}, "get:Login")
	beego.Router("/auth/login", &controllers.AuthController{}, "post:LoginHandler")
	beego.Router("/auth/register", &controllers.AuthController{}, "get:Register")
	beego.Router("/auth/register", &controllers.AuthController{}, "post:RegisterHandler")
	beego.Router("/panel", &controllers.PanelController{}, "get:Panel")

	beego.InsertFilter("/panel/*", beego.BeforeRouter, func(ctx *context.Context) {
		middleware.PanelFilter(ctx)
	})
}
