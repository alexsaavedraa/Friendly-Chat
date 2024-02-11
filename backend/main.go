package main

import (
	dbutils "backend/chat/dbutils"
	"backend/chat/pkg/websocket/websocket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request, username string, token string) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn:     conn,
		Pool:     pool,
		Username: username,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	// allow cors
	// Define a CORS middleware handler
	corsHandler := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				return
			}

			// Call the next handler
			handler.ServeHTTP(w, r)
		})
	}

	// Define the handler for the /auth route
	http.HandleFunc("/check-account", checkAcc)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/history", MessageHist)
	http.HandleFunc("/logout", logout)

	// Define the handler for the /ws route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Print authentication token to console
		username := r.URL.Query().Get("username")
		token := r.URL.Query().Get("token")

		fmt.Println("authenitcating user for ws", username, token)
		if dbutils.FindToken(token, username) {
			serveWs(pool, w, r, username, token)
		} else {

			return
		}
	})
	http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux))
}
func MessageHist(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	token := r.URL.Query().Get("token")

	fmt.Println("authenitcating user for message history", username, token)
	if dbutils.FindToken(token, username) {
		messageHistory := dbutils.GetMessageHistory(10)
		messageHistoryJSON, err := json.Marshal(messageHistory)
		if err != nil {
			// Handle error
			log.Println("Error marshaling message history:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(messageHistoryJSON)
		if err != nil {
			// Handle error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	token := r.URL.Query().Get("token")

	fmt.Println("authenitcating user logout", username, token)
	if dbutils.FindToken(token, username) {
		if dbutils.RemoveToken(token, username) {
			fmt.Println("successfully removedd user token: ", username)
		}
	}

}

// Handler function for the /auth route
func checkAcc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking user exists")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	var requestBody struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the request body

	// Extract the username from the struct
	username := requestBody.Username

	exists := dbutils.UsernameExists(username)

	response := map[string]bool{"userExists": exists}

	// Encode the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)

	// Send response based on username existence

}

func authHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("authenticating user")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the request body
	// Print the request body
	username := login.Username
	password := login.Password
	fmt.Println(username, password)
	res := dbutils.AuthUser(username, password)
	fmt.Println("validated user", username)

	var jsonResponse map[string]string

	if res {
		token := dbutils.Make_and_store_token(username)
		jsonResponse = map[string]string{"message": "success",
			"token": token}
	} else {
		jsonResponse = map[string]string{"message": "failure"}
	}

	// Marshal the JSON response only once
	responseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		// Handle error
	}

}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("authenticating user")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the request body
	// Print the request body
	username := login.Username
	password := login.Password

	if dbutils.UsernameExists(username) {
		fmt.Println("User already exists")

	} else {
		dbutils.InsertUser(username, password)
	}

	fmt.Println(username, password)

}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	dbutils.Dbsetup()
	dbutils.Create_db_if_not_exists()
	setupRoutes()
	http.ListenAndServe(":8080", nil)

}
