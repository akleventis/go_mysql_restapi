package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *App) GetCats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var animals []Animal

	res, err := app.DB.Query("SELECT id, name, COALESCE(age, '') as age, COALESCE(color, '') as color, COALESCE(gender, '') as gender, COALESCE(breed, '') as breed, COALESCE(weight, '') as weight FROM cats")

	if err != nil {
		internalServiceError(w, err)
		return
	}
	defer res.Close()
	for res.Next() {
		var animal Animal
		// copy columns into struct, empty string if null value in DB (coalesce)
		err := res.Scan(&animal.ID, &animal.Name, &animal.Age, &animal.Color, &animal.Gender, &animal.Breed, &animal.Weight)
		if err != nil {
			internalServiceError(w, err)
			return
		}
		animals = append(animals, animal)
	}
	json.NewEncoder(w).Encode(animals)
}

func (app *App) GetCatById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	res, err := app.DB.Query("SELECT id, name, COALESCE(age, '') as age, COALESCE(color, '') as color, COALESCE(gender, '') as gender, COALESCE(breed, '') as breed, COALESCE(weight, '') as weight FROM cats WHERE id=?", params["id"])
	if err != nil {
		internalServiceError(w, err)
		return
	}
	defer res.Close()

	var animal Animal
	for res.Next() {
		// copy columns into struct, empty string if null value in DB (coalesce)
		err := res.Scan(&animal.ID, &animal.Name, &animal.Age, &animal.Color, &animal.Gender, &animal.Breed, &animal.Weight)
		if err != nil {
			internalServiceError(w, err)
			return
		}
	}
	json.NewEncoder(w).Encode(animal)
}

func (app *App) PostCat(w http.ResponseWriter, r *http.Request) {
	statement, err := app.DB.Prepare("INSERT INTO cats(name, age, color, gender, breed, weight) VALUES(?,?,?,?,?,?)")
	if err != nil {
		internalServiceError(w, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalServiceError(w, err)
		return
	}
	pet := make(map[string]string)
	json.Unmarshal(body, &pet)

	// must include a name, all other fields nullable => 422
	if pet["name"] == "" {
		http.Error(w, "You must provide a valid name field", http.StatusUnprocessableEntity)
		return
	}
	// if key not specified, input null into database (newNullString => helper func in utils.go)
	_, err = statement.Exec(pet["name"], newNullString(pet["age"]), newNullString(pet["color"]), newNullString(pet["gender"]), newNullString(pet["breed"]), newNullString(pet["weight"]))
	if err != nil {
		internalServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

// PATCH => not 'PUT'ting whole resource, only updating specified fields
func (app *App) UpdateCat(w http.ResponseWriter, r *http.Request) {
	// grab url params ("id")
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalServiceError(w, err)
		return
	}
	// create map of patch body request
	petUpdate := make(map[string]string)
	json.Unmarshal(body, &petUpdate)
	// build string from arbitrarily set params in request body (createUpdateString => helper func in utils.go)
	query := createUpdateString(petUpdate, "cats", params["id"])

	statement, err := app.DB.Prepare(query)
	if err != nil {
		internalServiceError(w, err)
		return
	}
	result, err := statement.Exec()
	if err != nil {
		internalServiceError(w, err)
		return
	}
	// Check if any rows were affected (wrong id)
	affected, _ := result.RowsAffected()
	if affected == 0 {
		resourceNotFound(w)
		return
	}
	fmt.Fprintf(w, "Cat with ID %s was updated", params["id"])
}

func (app *App) DeleteCat(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	statement, err := app.DB.Prepare("DELETE FROM cats WHERE id=?")
	if err != nil {
		internalServiceError(w, err)
		return
	}
	result, err := statement.Exec(params["id"])
	if err != nil {
		internalServiceError(w, err)
		return
	}
	// Check if any rows were affected (wrong id)
	affected, _ := result.RowsAffected()
	if affected == 0 {
		resourceNotFound(w)
		return
	}
	fmt.Fprintf(w, "Cat with ID %s was deleted", params["id"])
}
