package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
)

func Login_with_password(req Req) (Res, error) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(req, &data)
	if err != nil {
		log.Println("Failed to parse request")
	}

	// password
	const hash = "$argon2id$v=19$m=65536,t=1,p=12$OARECWZoZ+X7OY5rMHzFWg$9ZUSA8gMIslrdWnTgAdNAHDRZ3rUahlzAsHrM0s2jlk"
	match, err := argon2id.ComparePasswordAndHash(data.Password, hash)
	if err != nil {
		err = fmt.Errorf("failed to compare password with hash: %w", err)
		return nil, err
	}

	if match {
		return Res(`<script hx-swap-oob="beforeend:#ws">window.location.href = "/test" </script>`), nil
	} else {
		return Res(`<div hx-swap-oob="innerHTML:#status">Invalid Username or Password!</div>`), nil
	}
}

func register(req Req) (Res, error) {
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
	return nil, fmt.Errorf("unimplemented")
}
