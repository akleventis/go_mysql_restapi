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

	handler := handlers.App{
		DB: db,
	}

	router := mux.NewRouter()

	// router.HandleFunc("/dogs", handler.GetDogs).Methods("GET")
	// router.HandleFunc("/dog", handler.PostDog).Methods("POST")
	// router.HandleFunc("/dogs/{id}", handler.GetDogById).Methods("GET")
	// router.HandleFunc("/dogs/{id}", handler.UpdateDog).Methods("PATCH")
	// router.HandleFunc("/dogs/{id}", handler.DeleteDog).Methods("DELETE")

	router.HandleFunc("/cats", handler.GetCats).Methods("GET")
	router.HandleFunc("/cat", handler.PostCat).Methods("POST")
	router.HandleFunc("/cats/{id}", handler.GetCatById).Methods("GET")
	router.HandleFunc("/cats/{id}", handler.UpdateCat).Methods("PATCH")
	router.HandleFunc("/cats/{id}", handler.DeleteCat).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}
