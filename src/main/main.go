package main

import (
	database "go_mysql/src/db"
	handlers "go_mysql/src/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Print(err)
		return
	}
	defer db.Close()

	handlers.Db = db

	router := mux.NewRouter()

	router.HandleFunc("/dogs", handlers.GetDogs).Methods("GET")
	router.HandleFunc("/dog", handlers.PostDog).Methods("POST")
	router.HandleFunc("/dogs/{id}", handlers.GetDogById).Methods("GET")
	router.HandleFunc("/dogs/{id}", handlers.UpdateDog).Methods("PUT")
	router.HandleFunc("/dogs/{id}", handlers.DeleteDog).Methods("DELETE")

	router.HandleFunc("/cats", handlers.GetCats).Methods("GET")
	router.HandleFunc("/cat", handlers.PostCat).Methods("POST")
	router.HandleFunc("/cats/{id}", handlers.GetCatById).Methods("GET")
	router.HandleFunc("/cats/{id}", handlers.UpdateCat).Methods("PUT")
	router.HandleFunc("/cats/{id}", handlers.DeleteCat).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}
