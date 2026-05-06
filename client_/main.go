package main

import (
	_ "client_/routers"
	"os"
	"shared/helpers"
	"shared/validator"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	validator.Init()
	helpers.Init()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	beego.Run(":" + port)
}
