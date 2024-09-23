package db

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func SetupDB() *sql.DB {
	fmt.Println("Setting up DB...")

	DB, err := sql.Open("sqlite3", "./database.DB")

	if err != nil {
		log.Fatal(err)
	}

	initDB()
	return DB
}

func initDB() {

	sqlStmt := `
CREATE TABLE IF NOT EXISTS meal_type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS  meal (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id INTEGER NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (type_id) REFERENCES meal_type (id)
);

CREATE TABLE IF NOT EXISTS donation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    meal_id INTEGER NOT NULL,
    donor_name TEXT NOT NULL,
    claimed BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (meal_id) REFERENCES meal (id)
);

CREATE TABLE IF NOT EXISTS claim (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    donation_id INTEGER,
    claimer_name TEXT NOT NULL,
    fulfilled BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (donation_id) REFERENCES donation (id)
);
	`
	_, err := DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
