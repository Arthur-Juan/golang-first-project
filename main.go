package main

import (
	"crud/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/users", server.CreateUser).Methods("POST")
	router.HandleFunc("/users", server.SearchUsers).Methods("GET")
	router.HandleFunc("/users/{id}", server.SearchUser).Methods("GET")
	router.HandleFunc("/users/{id}", server.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", server.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))

}
