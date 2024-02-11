package main

import (
	dbutils "backend/chat/dbutils"
	"backend/chat/pkg/websocket/websocket"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
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

	// Define the handler for the /ws route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Print authentication token to console
		authToken := r.Header.Get("Sec-Websocket-Protocol")
		if authToken != "login" {
			fmt.Println("User authentication token:", authToken)
			serveWs(pool, w, r)
		} else if authToken == "login" {
			fmt.Println("User authentication token: Needs token", authToken)
			// http.Redirect(w, r, "https://freshman.tech", http.StatusTemporaryRedirect)
			return
		}
	})
	http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux))
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

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Print the request body
	fmt.Println(w, "Request body:", string(body))

	// You can add your authentication logic here

	// Send a response
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	dbutils.Dbsetup()
	dbutils.Create_db_if_not_exists()
	setupRoutes()
	http.ListenAndServe(":8080", nil)

}
