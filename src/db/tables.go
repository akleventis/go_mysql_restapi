package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func createDogTable(db *sql.DB) error {
	query := tableQueryString("dogs")
	res, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating dogs table: %s", err)
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %s", err)
		return err
	}
	return nil
}

func createCatTable(db *sql.DB) error {
	query := tableQueryString("cats")
	res, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating cats table: %s", err)
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %s", err)
		return err
	}
	return nil
}

func tableQueryString(tableName string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(id int primary key auto_increment not null, name varchar(20), age int, color varchar(20), gender varchar(20), breed varchar(20), weight int)`, tableName)
}

func createTables(db *sql.DB) error {
	err := createDogTable(db)
	if err != nil {
		log.Print("Error creating dog table:", err)
		return err
	}
	err = createCatTable(db)
	if err != nil {
		log.Print("Error creating cat table:", err)
		return err
	}
	return nil
}
