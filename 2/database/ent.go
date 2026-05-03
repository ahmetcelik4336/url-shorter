package database

import (
	"2/ent"
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var Client *ent.Client

func Init() *ent.Client {

	var err error

	Client, err = ent.Open("mysql", "root:Ahmet.4336@tcp(127.0.0.1:3306)/beego2?parseTime=True")
	if err != nil {
		log.Fatal(err)
	}

	if err := Client.Schema.Create(context.Background()); err != nil {
		log.Fatal(err)
	}

	return Client
}
