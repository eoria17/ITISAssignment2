package main

import (
	"fmt"
	"net/http"

	"github.com/ITISAssignment2/config"
	"github.com/ITISAssignment2/storage"

	"github.com/gorilla/mux"

	pos "github.com/ITISAssignment2/controller"
)

func main() {
	storage_, err := storage.Open(config.DB_HOST, config.DB_NAME, config.DB_USER, config.DB_PASSWORD, config.DB_PORT)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	posEngine := pos.PosEngine{
		Storage: storage_,
	}
	posEngine.Route(router)

	fmt.Println("Listening to port 8080...")

	http.ListenAndServe(":80", router)

}
