package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./internal/data/transactions.db")
	if err != nil {
		log.Fatal("Could not open the file", err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATETIME DEFAULT CURRENT_TIMESTAMP,
		amount FLOAT NOT NULL,
		category TEXT NOT NULL,
		description TEXT,
		payment_method TEXT NOT NULL
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
