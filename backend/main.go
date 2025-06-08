package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Req struct {
	Data        []byte
	MessageType []string
	DB          *gorm.DB
	// storage
}
type Res []byte

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var DB *gorm.DB

func main() {
	var err error
	DB, err = InitDb()
	if err != nil {
		log.Fatal("DB init failed:", err)
	}

	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
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

	log.Println("User connected:", conn.RemoteAddr())

L:
	for {
		req_type, data, err := conn.ReadMessage()
		if err != nil {
			// Connection broken
			break
		}

		var req Req
		req.Data = data
		req.DB = DB

		var res Res
		switch req_type {
		case websocket.BinaryMessage:
			res, err = handle_binary_message(req)
		case websocket.TextMessage:
			res, err = handele_text_message(req)
		case websocket.CloseMessage:
			// TODO
			break L
		case websocket.CloseMessageTooBig:
			// TODO
			break L
		case websocket.PingMessage:
			conn.WriteMessage(websocket.PongMessage, []byte{})
		case websocket.PongMessage:
			conn.WriteMessage(websocket.PingMessage, []byte{})
		default:
			err = fmt.Errorf("websocket request type is invalid: %d", req_type)
		}

		if err != nil {
			log.Println("Error:", err)
		}

		if res != nil {
			conn.WriteMessage(websocket.TextMessage, res)
		}
	}

	log.Println("User disconnected:", conn.RemoteAddr())
}

func handele_text_message(req Req) (Res, error) {
	var meta struct {
		Headers struct {
			MessageType string `json:"HX-Trigger-Name"`
		} `json:"HEADERS"`
	}
	err := json.Unmarshal(req.Data, &meta)
	if err != nil {
		err = fmt.Errorf("failed to get MessageType from message: %w", err)
		return nil, err
	}

	req.MessageType = strings.Split(meta.Headers.MessageType, "#")
	if len(req.MessageType) < 1 {
		err = fmt.Errorf("message_type is invalid (len < 1): %v", req.MessageType)
		return nil, err
	}

	var res Res
	switch req.MessageType[0] {
	case "auth":
		res, err = auth_handler(req)
	default:
		err = fmt.Errorf("invalid MessageType: %s", meta.Headers.MessageType)
	}
	return res, err
}

func handle_binary_message(req Req) (Res, error) {
	// TODO
	return nil, fmt.Errorf("unimplemented")
}
