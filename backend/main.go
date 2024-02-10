package main

import (
	dbutils "backend/chat/dbutils"
	"backend/chat/pkg/websocket/websocket"
	"fmt"
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

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Print authentication token to console
		authToken := r.Header.Get("Sec-Websocket-Protocol")
		if authToken != "login" {
			fmt.Println("User authentication token:", authToken)
			serveWs(pool, w, r)
		} else if authToken == "login" {
			fmt.Println("User authentication token: Needs token", authToken)
			http.Redirect(w, r, "https://freshman.tech", http.StatusTemporaryRedirect)
			return
		}

	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	dbutils.Dbsetup()
	dbutils.Create_db_if_not_exists()
	setupRoutes()
	http.ListenAndServe(":8080", nil)

}
