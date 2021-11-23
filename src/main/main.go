package main

import (
	database "go_mysql/src/db"
	"log"
)

func main() {
    err := database.ConnectDatabase()
    if err != nil {
        log.Print(err)
        return
    }

}

