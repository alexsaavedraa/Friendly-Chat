package main

import (
	dbutils "backend/chat/dbutils"
	"backend/chat/pkg/websocket/websocket"
	"fmt"
	"io/ioutil"
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

	// Define the handler for the /auth route
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
}

// Handler function for the /auth route
func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Print the request body
	fmt.Println("Request body:", string(body))

	// You can add your authentication logic here

	// Send a response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Authentication successful")
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	dbutils.Dbsetup()
	dbutils.Create_db_if_not_exists()
	setupRoutes()
	http.ListenAndServe(":8080", nil)

}
