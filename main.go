package main

import (
	"fmt"
	"importxcel/db"
	"importxcel/service"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.DbConnect()

	r := mux.NewRouter()
	r.HandleFunc("/import", service.Upload).Methods("POST")
	r.HandleFunc("/employee/{id}", service.GetEmployee).Methods("GET")
	r.HandleFunc("/employee/{id}", service.DeleteEmployee).Methods("DELETE")
	r.HandleFunc("/employees/{id}", service.UpdateEmployee).Methods("PUT")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

}
