package main

import (
	"github.com/gorilla/mux"
	"location-history/api"
	"log"
	"net/http"
)

const HistoryServerListenAddr = ":8080"

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/location/{order_id}", api.AddLocationData).Methods("POST")
	router.HandleFunc("/location/{order_id}", api.GetLocationData).Methods("GET")
	router.HandleFunc("/location/{order_id}", api.DeleteLocationData).Methods("DELETE")

	server := &http.Server{
		Addr:    HistoryServerListenAddr,
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting the server %v\n", err)
	}
}
