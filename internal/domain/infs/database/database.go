package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {

	connectDB := `user=postgres password=tr12 dbname=rapi sslmode=disable`

	var err error

	DB, err = sql.Open("postgres", connectDB)

	if err != nil {
		return err
	}

	err = DB.Ping()

	if err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}
	fmt.Println("Database connection established successfully!")

	return nil
}
