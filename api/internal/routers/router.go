package routers

import (
	"api/controllers"
	"api/internal/container"
	"api/internal/middleware"

	_ "api/docs"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func Init(c *container.Container) {

	beego.SetStaticPath("/swagger", "swagger")

	beego.Router("/v1/url/:shortcode/?:password", &controllers.RedirectController{
		Service:    c.UrlService,
		LogService: c.LogService,
	}, "post:Redirect")

	beego.Router("/v1/create/:type", &controllers.UrlController{
		Service: c.UrlService,
	}, "post:CreateSite")

	beego.Router("/v1/user/register", &controllers.AuthController{
		Service: c.UserService,
	}, "post:Register")

	beego.Router("/v1/user/login", &controllers.AuthController{
		Service: c.UserService,
	}, "post:Login")

	beego.Router("/v1/user/profile", &controllers.AuthController{
		Service: c.UserService,
	}, "get:Profile")

	beego.Router("/v1/panel/create", &controllers.UrlController{
		Service: c.UrlService,
	}, "post:Create")

	beego.Router("/v1/panel/history", &controllers.UrlController{
		Service: c.UrlService,
	}, "get:History")

	beego.Router("/v1/panel/historybyid/?:id", &controllers.UrlController{
		Service: c.UrlService,
	}, "get:HistoryById")

	beego.Router("/v1/panel/update", &controllers.UrlController{
		Service: c.UrlService,
	}, "put:Update")

	beego.Router("/v1/panel/delete/:id", &controllers.UrlController{
		Service: c.UrlService,
	}, "delete:Delete")

	beego.Router("/v1/panel/bulkcreate", &controllers.UrlController{
		Service: c.UrlService,
	}, "post:Bulkcreate")

	beego.Router("/v1/panel/GetURLStats", &controllers.AnalysisController{
		Service: c.AnalysisService,
	}, "post:GetURLStats")

	beego.Router("/v1/panel/GetURLCount", &controllers.AnalysisController{
		Service: c.AnalysisService,
	}, "post:GetURLCount")

	beego.Router("/v1/user/GetUserCount", &controllers.AuthController{
		Service: c.UserService,
	}, "get:GetUserCount")

	beego.Router("/v1/user/GetPerDayClick", &controllers.AnalysisController{
		Service:    c.AnalysisService,
		LogService: c.LogService,
	}, "get:GetPerDayClick")

	beego.Router("/v1/user/GetURLCountTotal", &controllers.AnalysisController{
		Service:    c.AnalysisService,
		LogService: c.LogService,
	}, "get:GetURLCountTotal")

	beego.Router("/v1/user/ValidateToken", &controllers.AuthController{
		Service: c.UserService,
	}, "post:ValidateToken")

	beego.Router("/v1/panel/setting/save/:type", &controllers.SettingController{
		Service: c.SettingService,
	}, "post:Save")

	beego.Router("/v1/setting/get/general", &controllers.SettingController{
		Service: c.SettingService,
	}, "get:GetGeneralSettings")

	beego.Router("/v1/panel/TotalReading", &controllers.AnalysisController{
		Service:    c.AnalysisService,
		LogService: c.LogService,
	}, "post:TotalReading")

	beego.Router("/v1/panel/urlTrackAnalysis", &controllers.AnalysisController{
		Service:    c.AnalysisService,
		LogService: c.LogService,
	}, "post:UrlTrackAnalysis")

	beego.Router("/v1/panel/LogDatatable", &controllers.AnalysisController{
		Service:    c.AnalysisService,
		LogService: c.LogService,
	}, "post:LogDatatable")

	// =========================
	// JWT MIDDLEWARE
	// =========================
	beego.InsertFilter("v1/panel/*", beego.BeforeRouter, func(ctx *context.Context) {
		middleware.JWTMiddleware(ctx)
	})

	beego.InsertFilter("v1/user/profile/*", beego.BeforeRouter, func(ctx *context.Context) {
		middleware.JWTMiddleware(ctx)
	})

	beego.InsertFilter("v1/*", beego.BeforeRouter, func(ctx *context.Context) {
		middleware.ApiKeyMiddleware(ctx)
	})

	/*beego.InsertFilter("/v1/url/*", beego.BeforeRouter, func(ctx *context.Context) {
		redisIsActive, _ := beego.AppConfig.Bool("redis_active")
		if redisIsActive && !middleware.CheckRateLimit(ctx, "create_", 5) {
			return
		}
	})

	beego.InsertFilter("/v1/panel/bulkcreate/*", beego.BeforeRouter, func(ctx *context.Context) {
		redisIsActive, _ := beego.AppConfig.Bool("redis_active")
		if redisIsActive && !middleware.CheckRateLimit(ctx, "bulkcreate_", 3) {
			return
		}
	})*/

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:8081"}, // Client adresin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

}
