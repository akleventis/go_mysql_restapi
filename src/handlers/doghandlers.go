package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func GetDogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var animals []Animal

	res, err := Db.Query("SELECT id, name, COALESCE(age, '') as age, COALESCE(color, '') as color, COALESCE(gender, '') as gender, COALESCE(breed, '') as breed, COALESCE(weight, '') as weight FROM dogs")
	if err != nil {
		panic(err.Error())
	}
	defer res.Close()

	for res.Next() {
		var animal Animal
		err := res.Scan(&animal.ID, &animal.Name, &animal.Age, &animal.Color, &animal.Gender, &animal.Breed, &animal.Weight)
		if err != nil {
			panic(err.Error())
		}
		animals = append(animals, animal)
	}
	json.NewEncoder(w).Encode(animals)
}

func GetDogById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	res, err := Db.Query("SELECT id, name, COALESCE(age, '') as age, COALESCE(color, '') as color, COALESCE(gender, '') as gender, COALESCE(breed, '') as breed, COALESCE(weight, '') as weight FROM dogs WHERE id=?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer res.Close()

	var animal Animal
	for res.Next() {
		err := res.Scan(&animal.ID, &animal.Name, &animal.Age, &animal.Color, &animal.Gender, &animal.Breed, &animal.Weight)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(animal)
}

func PostDog(w http.ResponseWriter, r *http.Request) {
	statement, err := Db.Prepare("INSERT INTO dogs(name, age, color, gender, breed, weight) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	pet := make(map[string]string)
	json.Unmarshal(body, &pet)

	// must include a name, all other fields nullable
	if pet["name"] == "" {
		json.NewEncoder(w).Encode("You provide a name")
		return
	}

	_, err = statement.Exec(pet["name"], newNullString(pet["age"]), newNullString(pet["color"]), newNullString(pet["gender"]), newNullString(pet["breed"]), newNullString(pet["weight"]))
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Dog post success!")
}

// Technically a PATCH
func UpdateDog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	petUpdate := make(map[string]string)
	json.Unmarshal(body, &petUpdate)

	setString := "UPDATE dogs SET"
	for k, v := range petUpdate {
		setString += fmt.Sprintf(" %s = '%s',", k, v)
	}
	setString = setString[:len(setString)-1]
	setString += fmt.Sprintf(" WHERE id = %s", params["id"])
	fmt.Println(setString)
	statement, err := Db.Prepare(setString)
	if err != nil {
		panic(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID %s was updated", params["id"])
}

func DeleteDog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	statement, err := Db.Prepare("DELETE FROM dogs WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	_, err = statement.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID %s was deleted", params["id"])
}
