package main

import (
	_ "client_/routers"
	"shared/helpers"
	"shared/validator"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	validator.Init()
	helpers.Init()
	beego.Run()
}
