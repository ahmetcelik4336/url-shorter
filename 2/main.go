package main

import (
	"2/container"
	"2/database"
	"2/routers"
	_ "2/routers"
	"shared/validator"

	"github.com/beego/beego/v2/server/web"
	"github.com/redis/go-redis/v9"

	_ "2/docs"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"

	httpSwagger "github.com/swaggo/http-swagger"

	"2/middleware"
)

var appContainer *container.Container

// @title           Url Shortener API
// @version         1.0
// @description     Bu bir URL kısaltma servisidir.
// @host            localhost:8080
// @BasePath        /v1/
func main() {
	validator.Init()
	redisIsActive, _ := beego.AppConfig.Bool("redis_active")
	if redisIsActive {
		opt, err := redis.ParseURL("rediss://default:gQAAAAAAAbIoAAIgcDFmYzNjYWQ2MzZiMjE0MTFmODY5MmMzNTRhN2I0OWMxZg@learning-primate-111144.upstash.io:6379")
		if err != nil {
			panic(err)
		}

		client := redis.NewClient(opt)

		middleware.InitRateLimiter(client)
	}

	_ = godotenv.Load()

	db, driver := database.Init()

	appContainer = container.NewContainer(db, driver)

	routers.Init(appContainer)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.BConfig.WebConfig.DirectoryIndex = false

	// Beego'nun çalıştığı port üzerinden Swagger'ı açmak için:
	beego.Handler("/swagger/*", httpSwagger.WrapHandler)
	//port := os.Getenv("PORT")
	web.Run()
}
