package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		// Decode the received base64 username and password
		var data struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err = conn.ReadJSON(&data)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			break
		}

		log.Println("Data:", data)

		// hash, err := argon2id.CreateHash(data.Password, argon2id.DefaultParams)
		// if err != nil {
		// 	log.Println("Error CreateHash", err)
		// }
		//
		// log.Println("Hash:", hash)
		//
		// err = conn.WriteMessage(websocket.TextMessage, []byte(hash))
		// if err != nil {
		// 	log.Println("Error conn.WriteMessage", err)
		// }

		// password
		const hash = "$argon2id$v=19$m=65536,t=1,p=12$OARECWZoZ+X7OY5rMHzFWg$9ZUSA8gMIslrdWnTgAdNAHDRZ3rUahlzAsHrM0s2jlk"
		match, err := argon2id.ComparePasswordAndHash(data.Password, hash)
		if err != nil {
			log.Println("Error argon2id.ComparePasswordAndHash", err)
		}

		if match {
			log.Println("Login successful")
			// Send a confirmation message back to client if needed
			err = conn.WriteMessage(websocket.TextMessage, []byte("Welcome!"))
		} else {
			log.Println("Login failed")
			err = conn.WriteMessage(websocket.TextMessage, []byte("Invalid credentials"))
		}
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
