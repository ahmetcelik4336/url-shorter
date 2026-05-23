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
			beego.NSRouter("/logout", &controllers.AuthController{}, "get:LogOut"),
		),

		beego.NSRouter("/", &controllers.MainController{}),

		beego.NSNamespace("/panel",
			beego.NSRouter("/", &controllers.PanelController{}, "get:Panel"),
			beego.NSRouter("/urls/?:id", &controllers.PanelController{}, "get:Url;delete:UrlDelete"),
			beego.NSRouter("/urls/add/?:id", &controllers.PanelController{}, "get:UrlAdd"),
			beego.NSRouter("/urls/save/?:id", &controllers.PanelController{}, "post:UrlSaveHandler"),
			beego.NSRouter("/urls/batch", &controllers.PanelController{}, "post:BatchExcel"),
			beego.NSRouter("/urls/download", &controllers.PanelController{}, "get:DownloadTemplate"),
			beego.NSRouter("/urls/qr/:id", &controllers.PanelController{}, "get:Qr"),
			beego.NSRouter("/urls/qrcreate/:id", &controllers.PanelController{}, "get:QrCreate"),
			beego.NSRouter("/urls/logoUpload", &controllers.PanelController{}, "post:UploadLogo"),
		),
	)
	beego.AddNamespace(ns)

	beego.Router("/:type/:shortcode/?:password", &controllers.RedirectController{}, "get:Redirect")

	beego.InsertFilter("/:lang/panel/*", beego.BeforeRouter, middleware.PanelFilter)
	beego.InsertFilter("/:lang/panel", beego.BeforeRouter, middleware.PanelFilter)
	beego.InsertFilter("/:lang/auth/logout", beego.BeforeRouter, middleware.PanelFilter)

	beego.Get("/", func(ctx *context.Context) {
		lang := ctx.Input.Param(":lang")
		if lang == "" {
			lang = os.Getenv("DEFAULTLANG")
		}
		ctx.Redirect(302, "/"+lang)
	})

	beego.Get("/captcha/*.png", func(ctx *context.Context) {
		controllers.Cpt.Handler(ctx)
	})

	beego.Get("/captcha/refresh", func(ctx *context.Context) {
		controllers.Cpt.Handler(ctx)
	})

}
