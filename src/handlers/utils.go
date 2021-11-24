package handlers

import "database/sql"

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
