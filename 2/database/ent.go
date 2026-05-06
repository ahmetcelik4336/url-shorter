package database

import (
	"2/ent"
	"context"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var Client *ent.Client

/*
func InitMysql() *ent.Client {

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
*/
func Init() *ent.Client {

	driver := os.Getenv("DATABASE_DRIVER")
	dsn := os.Getenv("DATABASE_URL")

	client, err := ent.Open(driver, dsn)
	if err != nil {
		log.Fatalf("PostgreSQL bağlantısı kurulamadı: %v", err)
	}

	// 3. Tabloları oluştur (Auto Migration)
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Schema oluşturulamadı: %v", err)
	}

	return client
}
