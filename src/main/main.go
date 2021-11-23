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

	// router.HandleFunc("/dogs", getDogs).Methods("GET")
	// router.HandleFunc("/dogs", createDog).Methods("POST")
	// router.HandleFunc("/dogs/{id}", getDog).Methods("GET")
	// router.HandleFunc("/dogs/{id}", updateDog).Methods("PUT")
	// router.HandleFunc("/dogs/{id}", deleteDog).Methods("DELETE")

	router.HandleFunc("/cats", handlers.GetCats).Methods("GET")
	router.HandleFunc("/cats", handlers.PostCat).Methods("POST")
	router.HandleFunc("/cats/{id}", handlers.GetCatById).Methods("GET")
	// router.HandleFunc("/cats/{id}", updateCat).Methods("PUT")
	// router.HandleFunc("/cats/{id}", deleteCat).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}
