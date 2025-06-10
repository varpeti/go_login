package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	var err error
	DB, err = InitDb()
	if err != nil {
		log.Fatal("DB init failed:", err)
	}

	http.HandleFunc("/ws", handle_web_socket)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func handle_web_socket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket", err)
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
			err = MyErrorf("websocket request type is invalid: %d", req_type)
		}

		if err != nil {
			log.Println("Error:", err)
		}

		if res.Len() > 0 {
			conn.WriteMessage(websocket.TextMessage, res.Bytes())
		}
	}

	log.Println("User disconnected:", conn.RemoteAddr())
}

func handele_text_message(req Req) (Res, error) {
	var res Res
	var err error

	var message_type string
	message_type, err = NextMessagType(&req)
	if err != nil {
		return res, err
	}

	switch message_type {
	case "auth":
		res, err = Auth_handler(req)
	default:
		err = MyErrorf("invalid MessageType: %s", message_type)
	}
	return res, err
}

func handle_binary_message(req Req) (Res, error) {
	// TODO
	return Res{}, MyErrorf("unimplemented")
}
