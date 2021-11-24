package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Animal struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    string `json:"age"`
	Color  string `json:"color"`
	Gender string `json:"gender"`
	Breed  string `json:"breed"`
	Weight string `json:"weight"`
}

var Db *sql.DB

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func internalServiceError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	log.Printf("Error: %s", err)
}

func createUpdateString(patchMap map[string]string, table string, id string) string {
	setStr := fmt.Sprintf("UPDATE %s SET", table)
	for k, v := range patchMap {
		setStr += fmt.Sprintf(" %s = '%s',", k, v)
	}
	setStr = setStr[:len(setStr)-1]
	setStr += fmt.Sprintf(" WHERE id = %s", id)
	return setStr
}
