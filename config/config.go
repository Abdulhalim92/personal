package config

import (
	"encoding/json"
	"log"
	"moneytracker/internal/models"
	"os"
)

func GetConfig() (*models.Config, error) {
	// чтение и десериализация конфигов
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Println(err) // todo
		return nil, err
	}

	var config models.Config

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Println(err) // todo
		return nil, err
	}

	return &config, nil
}
