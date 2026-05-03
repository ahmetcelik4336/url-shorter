package main

import (
	_ "client_/routers"
	"shared/validator"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	validator.Init()
	beego.Run()
}
