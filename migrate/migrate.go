package main

import (
	"github.com/ITISAssignment2/config"
	"github.com/ITISAssignment2/model"
	"github.com/ITISAssignment2/storage"
)

func main() {
	storage_, err := storage.Open(config.DB_HOST, config.DB_NAME, config.DB_USER, config.DB_PASSWORD, config.DB_PORT)
	if err != nil {
		panic(err)
	}
	defer storage_.DB.Close()

	db := storage_.DB

	db.AutoMigrate(&model.Menu{})
	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&model.OrderLine{})
}
