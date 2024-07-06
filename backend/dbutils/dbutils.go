package dbutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// the environmental variables used to access database.
var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "testdb"
)

// Token Stack store currently authenticated users
var tokenStack [][]string

// Given a username, returns if the username exists in the database
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
	var scanErr error
	err = db.QueryRow(checkUserExistsQuery, username).Scan(&exists)
	if err != nil {
		scanErr = err
		log.Fatal("Error connecting to the database: ", scanErr)
	}
	return exists
}

// Authenticates a username and password, returns bool
func AuthUser(username, inpassword string) bool {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()
	if !UsernameExists(username) {

	}
	var userID string
	row := db.QueryRow("SELECT user_id FROM users WHERE username = $1", username)
	err = row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			fmt.Println(err)
		}
		fmt.Println(err)
	}
	var hashedPasswordFromDB string
	// Query to fetch hashed password based on user_id
	row = db.QueryRow("SELECT password_hash FROM passwords WHERE user_id = $1", userID)
	err = row.Scan(&hashedPasswordFromDB)
	if err != nil {
		if err == sql.ErrNoRows {
			// Password not found
			return false
		}
		return false
	}

	// Compare fetched hashed password with provided hashed password
	return comparePasswords(inpassword, hashedPasswordFromDB)

}

// Compare hashed password from DB with provided password using BCRYPT.
func comparePasswords(providedPassword string, hashedPasswordFromDB string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordFromDB), []byte(providedPassword))
	fmt.Println(err)
	return err == nil
}

// Hashes a password with BCRYPT.
func hashpassword(unhashed_password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(unhashed_password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error connecting to the Postgres: ", err)
	}
	res := string(hashed)

	return res
}

// for a given username, make and store the auth token.
func Make_and_store_token(username string) string {
	token := uuid.New().String()
	for i, pair := range tokenStack {
		if pair[1] == username {
			// Remove the existing username-token pair.
			tokenStack = append(tokenStack[:i], tokenStack[i+1:]...)
			break // Exit loop once removed
		}
	}
	// Add the new username-token pair.
	tokenStack = append(tokenStack, []string{token, username})
	fmt.Println(tokenStack)
	return token
}

// returns if the token username combo is valid
func FindToken(token string, username string) bool {
	fmt.Println(token, username, "\n token stak is", tokenStack)
	for _, pair := range tokenStack {
		if pair[1] == username && pair[0] == token {
			return true
		}
	}
	return false
}

// removes a username token pair.
func RemoveToken(token string, username string) bool {
	for i, pair := range tokenStack {
		if pair[1] == username && pair[0] == token {
			tokenStack = append(tokenStack[:i], tokenStack[i+1:]...)
			return true
		}
	}
	return false
}

// updates the votes table with a new vote
func UpdateVotes(messageID, username, timestamp, status string) int {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()
	stmt := `INSERT INTO votes (msg_id, username, vote_status, created_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (msg_id, username) DO UPDATE
	SET vote_status = $3, created_at = $4`

	_, err = db.Exec(stmt, messageID, username, status, timestamp)
	if err != nil {
		log.Fatal("Error connecting adding vote to database: ", err)
	}
	return countvotes(messageID)
}

// for a given username and message, returns what the user has voted for.
func uservotes(username, messageID string) string {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	var voteStatus string

	query := `SELECT vote_status FROM votes WHERE msg_id = $1 AND username = $2;`

	rows, err := db.Query(query, messageID, username)
	if err != nil {
		log.Fatal("Error executing query: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&voteStatus); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating over rows: ", err)
	}

	return voteStatus

}

// counts the votes for a given message.
func countvotes(messageID string) int {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	var voteStatuses []string

	// Iterate over the rows and append vote_status values to the slice
	query := `SELECT vote_status FROM votes WHERE msg_id = $1`

	// Execute the SQL statement
	rows, err := db.Query(query, messageID)
	if err != nil {
		log.Fatal("error executing SQL statement:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var voteStatus string
		if err := rows.Scan(&voteStatus); err != nil {
			log.Fatal("Error connecting to the database: ", err)
		}
		voteStatuses = append(voteStatuses, voteStatus)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("error iterating over rows:", err)
	}
	//fmt.Println(voteStatuses)
	upCount := 0
	downCount := 0
	for _, str := range voteStatuses {
		switch str {
		case "up":
			upCount++
		case "down":
			downCount++
		}
	}
	return upCount - downCount

}

// adds a message to the message table.
func AddMessage(body, category, timestamp, username string) string {
	//fmt.Println(" adding message: ", body, category, timestamp, username)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	var userID string
	row := db.QueryRow("SELECT user_id FROM users WHERE username = $1", username)
	err = row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			fmt.Println(err)
		}
		fmt.Println(err)
	}

	var id int64
	err = db.QueryRow("INSERT INTO messages (user_id, username, body, created_at) VALUES ($1, $2, $3, $4) RETURNING id", userID, username, body, timestamp).Scan(&id)
	if err != nil {
		log.Fatal("Error inserting user ope : ", err)
	}
	//fmt.Println("Inserted message with ID:", id)

	return fmt.Sprint(id)
}

type Message struct {
	ID        string `json:"MessageID"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"time"`
	Votes     string `json:"votes"`
	Category  string `json:"category"`
	Uservote  string `json:"user_vote"`
}

// gets x number of previous messages.
// TODO SQL calls need to be optimized.
func GetMessageHistory(number int, username string) []Message {

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	query := `
        SELECT id, body, user_id, username, created_at
        FROM messages
        ORDER BY id DESC
        LIMIT 10;
    `

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(rows)
	defer rows.Close()

	// Iterate through the result rows
	var res []Message
	for rows.Next() {
		var msg Message
		msg.Category = "message"
		err := rows.Scan(&msg.ID, &msg.Body, &msg.UserID, &msg.Username, &msg.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		msg.Votes = fmt.Sprint(countvotes(msg.ID))
		msg.Uservote = uservotes(username, msg.ID)
		res = append([]Message{msg}, res...)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return res

}

// Sets up the database on first run, inserting a dummy user.
func Dbsetup() {
	updateIfSet := func(envVar string, value *string, defaultValue *string) {
		envValue := os.Getenv(envVar)
		if envValue != "" {
			*value = envValue
			fmt.Printf("%s is set to: %s\n", envVar, envValue)
		} else {
			if defaultValue != nil {
				*value = *defaultValue
			}
			fmt.Printf("%s is not set in your environment. Using default value: %s\n", envVar, *value)
		}
	}

	updateIfSet("DB_HOST", &host, &host)
	updateIfSet("DB_USERNAME", &user, &user)
	updateIfSet("DB_PASSWORD", &password, &password)
	updateIfSet("DB_NAME", &dbname, &dbname)

}

// if the database doesnt exist, this creates one.
func Create_db_if_not_exists() {

	fmt.Println("beginning creation of DB tables.go")

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password)
	fmt.Println(host)
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
		id SERIAL,
		username VARCHAR,
		user_id VARCHAR PRIMARY KEY,
		is_active BOOLEAN,
		created_at TIMESTAMP,
		last_online TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS passwords (
		user_id VARCHAR PRIMARY KEY,
		password_hash VARCHAR,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		body TEXT,
		user_id VARCHAR,
		username VARCHAR,
		created_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);
	CREATE TABLE IF NOT EXISTS votes (
		id SERIAL PRIMARY KEY,
		username VARCHAR,
		msg_id INTEGER,
		vote_status VARCHAR,
		created_at TIMESTAMP,
		FOREIGN KEY (msg_id) REFERENCES messages(id),
		CONSTRAINT votes_msg_id_username_key UNIQUE (msg_id, username)
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	fmt.Println("Database and table created successfully! \n Inserting dummy User")

	a := UsernameExists("Alex")
	if !a {
		InsertUser("Alex", "password")
		InsertUser("Prachi", "QT")
		InsertUser("DemoUser", "password")
	} else {
		fmt.Println("Dummy already exists")
	}

}

// inserts a user, generates password hash, userID etc.
// Cannot insert a user with the same username.
func InsertUser(username string, inpassword string) {

	fmt.Println("inserting ", username, inpassword)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	userID := uuid.New().String()

	_, err = db.Exec("INSERT INTO users (username, user_id,is_active, created_at, last_online ) VALUES ($1, $2, $3,  NOW(), NOW())", username, userID, true)
	if err != nil {
		log.Fatal("Error inserting user ope : ", err)
	}

	// Get the ID of the inserted user

	hashed_password := hashpassword(inpassword)

	// Insert the hashed password into the passwords table
	_, err = db.Exec("INSERT INTO passwords (password_hash, user_id) VALUES ($1, $2)", hashed_password, userID)
	if err != nil {
		log.Fatal("Error inserting user hiii: ", err)
	}

}
