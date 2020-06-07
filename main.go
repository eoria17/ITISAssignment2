package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tiuriandy/ITISAssignment2/config"
	"github.com/tiuriandy/ITISAssignment2/storage"

	"github.com/gorilla/mux"

	pos "github.com/tiuriandy/ITISAssignment2/controller"
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

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/"))))
	http.Handle("/assets/", router)

	fmt.Println("Currently Listening to port 8080..")

	log.Println(http.ListenAndServe(":8080", router))

}
