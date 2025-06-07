package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	Req []byte
	Res []byte
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
		req_type, req, err := conn.ReadMessage()
		if err != nil {
			// Connection broken
			break
		}

		var res Res
		switch req_type {
		case websocket.BinaryMessage:
			res, err = handle_binary_message(req)
		case websocket.TextMessage:
			res, err = handele_text_message(req)
		case websocket.CloseMessage:
			// TODO
		case websocket.CloseMessageTooBig:
			// TODO
		case websocket.PingMessage:
			conn.WriteMessage(websocket.PongMessage, []byte{})
		case websocket.PongMessage:
			conn.WriteMessage(websocket.PingMessage, []byte{})
		default:
			err = fmt.Errorf("websocket request type is invalid: %d", req_type)
		}

		if err != nil {
			log.Println("Error: ", err)
		}

		if res != nil {
			conn.WriteMessage(websocket.TextMessage, res)
		}
	}
}

func handele_text_message(message []byte) (Res, error) {
	var meta struct {
		Headers struct {
			MessageType string `json:"HX-Trigger-Name"`
		} `json:"HEADERS"`
	}
	err := json.Unmarshal(message, &meta)
	if err != nil {
		err = fmt.Errorf("failed to get MessageType from message: %w", err)
		return nil, err
	}

	var res Res
	switch meta.Headers.MessageType {
	case "login_with_password":
		res, err = Login_with_password(message)
	default:
		err = fmt.Errorf("invalid MessageType: %s", meta.Headers.MessageType)
	}
	return res, err
}

func handle_binary_message(message []byte) (Res, error) {
	// TODO
	return nil, fmt.Errorf("unimplemented")
}
