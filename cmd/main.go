package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"moneytracker/config"
	"moneytracker/internal/db"
	"moneytracker/internal/repository"
	"moneytracker/internal/server"
	"moneytracker/internal/services"
	"moneytracker/pkg/logging"
	"net"
	"net/http"
	"os"
)

func recoverAll() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}

func main() {
	log := logging.GetLogger()

	defer recoverAll()

	err := execute()
	if err != nil {
		log.Println(err)
		log.Println(err)
		os.Exit(1)
	}
}

func execute() error {

	// регистрация роутеров
	router := mux.NewRouter()

	// подключение к БД
	connection, err := db.GetDbConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// переменная для работы с БД
	newRepository := repository.NewRepository(connection)

	// обращение к сервису
	service := services.NewService(newRepository)

	// обращение к серверу
	newServer := server.NewServer(router, service)

	// получение конфигураций
	getConfig, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// запуск роутеров
	newServer.Init()

	// получение адреса для работы с сервером
	address := net.JoinHostPort(getConfig.Host, getConfig.Port)

	// структура сервер
	srv := http.Server{
		Addr:    address,
		Handler: newServer,
	}

	// обслуживание сервера
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
