package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "gosql"
	password = "gosql"
	hostname = "127.0.0.1:3306"
	dbname   = "pets"
)

func sqlString(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("mysql", sqlString(""))
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}

	res, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		log.Printf("Error while creating database: %s", err)
		return nil, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error while fetching rows: %s", err)
	}

	db.Close()
	db, err = sql.Open("mysql", sqlString(dbname))
	if err != nil {
		log.Printf("Error opening db: %s", err)
		return nil, err
	}
	log.Printf("Successfully connected to %s db\n", dbname)
	return db, nil
}

func ConnectDatabase() (*sql.DB, error) {
	db, err := dbConnect()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return nil, err
	}

	// create tables
	err = createTables(db)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Printf("Successfully updated tables")
	return db, nil
}
