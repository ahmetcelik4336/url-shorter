package main

import (
	_ "client_/routers"
	"log"
	"os"
	"shared/helpers"
	"shared/utils"
	"shared/validator"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := utils.InitUtils(); err != nil {
		log.Fatal(err)
	}

	if err := utils.InitJWT(); err != nil {
		log.Fatal(err)
	}

	validator.Init()
	helpers.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	beego.Run(":" + port)
}
