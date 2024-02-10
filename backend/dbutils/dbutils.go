package dbutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func Function2() {
	fmt.Println("This is function 2 from file2.go")
}

var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "testdb"
)

func Dbsetup() {
	updateIfSet := func(envVar, value string, defaultValue *string) {
		envValue := os.Getenv(envVar)
		if envValue != "" {
			value = envValue
			fmt.Printf("%s is set to: %s\n", envVar, envValue)
		} else {
			if defaultValue != nil {
				value = *defaultValue
			}
			fmt.Printf("%s is not set in your environment. Using default value: %s\n", envVar, value)
		}
	}

	updateIfSet("DB_HOST", host, &host)
	updateIfSet("DB_PORT", strconv.Itoa(port), nil)
	updateIfSet("DB_USERNAME", user, &user)
	updateIfSet("DB_PASSWORD", password, &password)
	updateIfSet("DB_NAME", dbname, &dbname)
}

func Create_db_if_not_exists() {

	fmt.Println("beginning creation of DB tables.go")
	// Connect to PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Create the table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR,
		user_id VARCHAR,
		is_active BOOLEAN,
		created_at TIMESTAMP,
		last_online TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS passwords (
		password_id SERIAL PRIMARY KEY,
		password_hash VARCHAR,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		msg_id VARCHAR,
		body TEXT,
		user_id INTEGER,
		created_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS votes (
		id SERIAL PRIMARY KEY,
		user_id INTEGER,
		msg_id INTEGER,
		vote_status VARCHAR,
		created_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (msg_id) REFERENCES messages(id)
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	fmt.Println("Table my_table created successfully!")
}
