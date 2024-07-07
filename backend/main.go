package main

import (
	dbutils "backend/chat/dbutils"
	"backend/chat/pkg/websocket/websocket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// This function is called whenever the websocket endpoint is hit.
// it accepts the ws pool, and adds a new client.
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

// We use setup routes because in development it is necessary to use cors.
// Due to the frontend and backend running on different.
func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

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

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		username := r.URL.Query().Get("username")
		token := r.URL.Query().Get("token")
		fmt.Println("authenitcating user for ws", username, token)
		if dbutils.FindToken(token, username) {
			serveWs(pool, w, r, username, token)
		} else {

			return
		}
	})
	//Serve the static page
	fs := http.FileServer(http.Dir("build"))
	http.Handle("/", fs)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "build/index.html")
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "build/index.html")
	})

	err := http.ListenAndServe(":80", corsHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux))
}

// This route gets the queries the message history.
// The SQl calls need to be fixed with a join statement
// In order to run faster.
func MessageHist(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting history")
	username := r.URL.Query().Get("username")
	token := r.URL.Query().Get("token")

	if dbutils.FindToken(token, username) {
		fmt.Println("authenitcating user for message history", username, token)
		messageHistory := dbutils.GetMessageHistory(10, username)
		length := len(messageHistory)

		messageHistoryJSON, err := json.Marshal(messageHistory)
		if length == 0 {
			//messageHistory = []string{}
			messageHistoryJSON, err = json.Marshal([]string{})
		}

		if err != nil {
			// Handle error
			log.Println("Error marshaling message history:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//fmt.Println(messageHistory)
		_, err = w.Write(messageHistoryJSON)
		if err != nil {
			// Handle error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}
}

// On Logout, the user auth token is deleted from the stack
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

// Checks if a given account is in the database
func checkAcc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	username := requestBody.Username

	exists := dbutils.UsernameExists(username)

	response := map[string]bool{"userExists": exists}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	username := login.Username
	password := login.Password
	res := dbutils.AuthUser(username, password)
	var jsonResponse map[string]string

	if res {
		token := dbutils.Make_and_store_token(username)
		jsonResponse = map[string]string{"message": "success",
			"token": token}
	} else {
		jsonResponse = map[string]string{"message": "failure"}
	}
	responseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		// TODO Handle error
	}

}

// This will accept a json of username and password, and try to sign them up. Duplicate users not allowed
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

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

// Alex's Chat App
func main() {
	fmt.Println("Alex's Chat App")
	dbutils.Dbsetup()
	dbutils.Create_db_if_not_exists()
	setupRoutes()
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
