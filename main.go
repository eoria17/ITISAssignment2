package main

import (
	"fmt"
	"net/http"

	"github.com/tiuriandy/ITISAssignment2/config"
	"github.com/tiuriandy/ITISAssignment2/storage"

	"github.com/gorilla/mux"

	pos "github.com/tiuriandy/ITISAssignment2/controller"
)

func main() {
	storage_, err := storage.Open(config.DB_HOST, config.DB_NAME, config.DB_USER, config.DB_PASSWORD)
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

	fmt.Println("Listening to port 80")

	http.ListenAndServe(":80", router)

}
