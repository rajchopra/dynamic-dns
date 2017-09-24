package main

import (
	"fmt"
	"log"
	"net/http"

	"dynamic-dns/api"

	"github.com/gorilla/mux"
)

func main() {
	var configPath = "/etc/dynamicdns.json"
	api.Init(configPath)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/add", api.RequestHandlerForAdd).Methods("GET")
	router.HandleFunc("/update", api.RequestHandlerForUpdate).Methods("GET")
	router.HandleFunc("/delete", api.RequestHandlerForDelete).Methods("GET")

	log.Println(fmt.Sprintf("Starting dynamic-dns updator service on 0.0.0.0:8080..."))
	log.Fatal(http.ListenAndServe(":8080", router))
}
