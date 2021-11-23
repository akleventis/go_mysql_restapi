package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Animal struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Color  string `json:"color"`
	Gender string `json:"gender"`
	Breed  string `json:"breed"`
	Weight int    `json:"weight"`
}

var Db *sql.DB

func GetCats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var animals []Animal

	res, err := Db.Query("SELECT * from cats")
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

func PostCat(w http.ResponseWriter, r *http.Request) {
	statement, err := Db.Prepare("INSERT INTO cats(name, age, color, gender, breed, weight) VALUES(?,?,?,?,?,?)")
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
		json.NewEncoder(w).Encode("You must at least provide a name")
		return
	}

	_, err = statement.Exec(pet["name"], newNullString(pet["age"]), newNullString(pet["color"]), newNullString(pet["gender"]), newNullString(pet["breed"]), newNullString(pet["weight"]))
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Cat post success!")
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func GetCatById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// no more null error (:
	res, err := Db.Query("SELECT id, name, COALESCE(age, '') as age, COALESCE(color, '') as color, COALESCE(gender, '') as gender, COALESCE(breed, '') as breed, COALESCE(weight, '') as weight FROM cats WHERE id=?", params["id"])
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
