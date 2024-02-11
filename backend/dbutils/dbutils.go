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

func UsernameExists(username string) bool {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()
	var exists bool
	checkUserExistsQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	// Execute the query and scan the result into 'exists' variable
	var scanErr error
	fmt.Println("checking if user exists: ", username)
	err = db.QueryRow(checkUserExistsQuery, username).Scan(&exists)
	if err != nil {
		scanErr = err
		log.Fatal("Error connecting to the database: ", scanErr)
	}

	return exists
}

func AuthUser(usr, pass string) {
	fmt.Println(UsernameExists(usr))
}

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

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the Postgres: ", err)
	}
	defer db.Close()

	// Create the database if it doesn't exist
	// Check if the database exists before attempting to create it
	checkDatabaseQuery := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'", dbname)
	rows, err := db.Query(checkDatabaseQuery)
	if err != nil {
		log.Fatal("Error checking if database exists: ", err)
	}
	defer rows.Close()

	var exists bool
	for rows.Next() {
		exists = true
	}

	if !exists {
		// If the database does not exist, create it
		createDatabaseQuery := fmt.Sprintf("CREATE DATABASE %s", dbname)
		_, err = db.Exec(createDatabaseQuery)
		if err != nil {
			log.Fatal("Error creating database: ", err)
		}
	} else {
		fmt.Println("Database already existed")
	}

	// Connect to PostgreSQL
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", connStr)
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

	fmt.Println("Database and table created successfully! \n Inserting dummy User")

	insertStatement := `
	INSERT INTO users (username, user_id, is_active, created_at, last_online) 
	VALUES 
		($1, $2, $3, NOW(), NOW()),
		($4, $5, $6, NOW(), NOW());
`

	// Execute the SQL statement

	a := UsernameExists("Alex")
	if !a {
		_, err = db.Exec(insertStatement, "Prachi", "123456", true, "Alex", "789012", false)
		if err != nil {
			log.Fatal("Error creating table: ", err)
		}

	} else {
		fmt.Println("Dummy already exists")
	}
	//AuthUser("Prachi", "secretdata")

}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionToken struct {
	Token string `json:"token"`
}
