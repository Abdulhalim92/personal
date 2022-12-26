package db

import (
	"fmt"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moneytracker/pkg/logging"
)

func GetDbConnection() (*gorm.DB, error) {
	log := logging.GetLogger()

	// todo надо хранить в конфигах данные
	host := "localhost"
	port := "5432"
	user := "humo"
	password := "pass"
	dbname := "accounting_db"

	connString := fmt.Sprintf("host=%s user =%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Dushanbe",
		host, user, password, dbname, port)
	conn, err := gorm.Open(postgresDriver.Open(connString))
	if err != nil {
		log.Printf("%s GoPostgresConnection -> Open error", err.Error())
		return nil, err
	}

	log.Println("Connection success", host)

	return conn, nil
}
