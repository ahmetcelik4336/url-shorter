package routers

import (
	"client_/controllers"
	"client_/middleware"
	"os"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	ns := beego.NewNamespace("/:lang",

		beego.NSBefore(func(ctx *context.Context) {
			lang := ctx.Input.Param(":lang")
			if lang == "" {
				lang = os.Getenv("DEFAULTLANG")
			}
			ctx.Input.SetData("lang", lang)
		}),

		beego.NSNamespace("/auth",
			beego.NSRouter("/login", &controllers.AuthController{}, "get:Login;post:LoginHandler"),
			beego.NSRouter("/register", &controllers.AuthController{}, "get:Register;post:RegisterHandler"),
		),

		beego.NSRouter("/", &controllers.MainController{}),

		beego.NSNamespace("/panel",
			beego.NSRouter("/", &controllers.PanelController{}, "get:Panel"),
		),
	)

	beego.AddNamespace(ns)

	beego.InsertFilter("/:lang/panel/*", beego.BeforeRouter, middleware.PanelFilter)
	beego.InsertFilter("/:lang/panel", beego.BeforeRouter, middleware.PanelFilter)

	beego.Get("/", func(ctx *context.Context) {
		lang := ctx.Input.Param(":lang")
		if lang == "" {
			lang = os.Getenv("DEFAULTLANG")
		}
		ctx.Redirect(302, "/"+lang)
	})
}
