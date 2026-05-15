package database

import (
	"api/ent"
	"context"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var Client *ent.Client

func Init() (*ent.Client, string) {

	driver := os.Getenv("DATABASE_DRIVER")
	//dsn := os.Getenv("DATABASE_URL_MYSQL")
	dsn := os.Getenv("DATABASE_URL")
	client, err := ent.Open(driver, dsn)
	if err != nil {
		log.Fatalf("PostgreSQL bağlantısı kurulamadı: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Schema oluşturulamadı: %v", err)
	}

	return client, driver
}
