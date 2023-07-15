package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thallesrangel/crud_go/logic"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user", logic.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user", logic.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", logic.GetUserById).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", logic.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/user/{id}", logic.DeleteUser).Methods(http.MethodDelete)

	fmt.Println("Listen")
	log.Fatal(http.ListenAndServe(":5000", router))
}
