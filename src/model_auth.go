package main

import (
	"encoding/json"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email   string `gorm:"unique"`
	Pw_hash string
}

func Auth_handler(req Req) Res {
	message_type, err := NextMessagType(&req)
	if err != nil {
		Println(err)
		return nil
	}

	var res Res
	switch message_type {
	case "login_page":
		res = Jet("login_page.html", nil)
	case "login_with_password":
		res = login_with_password(req)
	case "register_page":
		res = Jet("register_page.html", nil)
	case "register":
		res = register(req)
	default:
		Println("invalid MessageType:", message_type)
		return nil
	}
	return res
}

func login_with_password(req Req) Res {
	var data struct {
		Email    string
		Password string
	}

	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		Println("failed to parse request:", err)
		return nil
	}

	var users []User
	result := DB.Where("email = ?", data.Email).Find(&users)
	if result.Error != nil {
		Println("db error: failed to get user by email:", result.Error)
		return nil
	}

	invalid_email_or_password := Jeti("invalid_email_or_password", `<div hx-swap-oob="innerHTML:#status_login">Invalid Email or Password!</div>`, nil)

	if len(users) != 1 {
		return invalid_email_or_password
	}

	user := users[0]
	match, err := argon2id.ComparePasswordAndHash(data.Password, user.Pw_hash)
	if err != nil {
		Println("failed to compare password with hash:", err)
		return nil
	}

	if !match {
		return invalid_email_or_password
	}

	vars := make(jet.VarMap)
	vars.Set("email", user.Email)
	return Jet("home_page.html", vars)
}

func register(req Req) Res {
	var data struct {
		Email    string
		Password string
	}

	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		Println("failed to parse request:", err)
		return nil
	}

	pw_hash, err := argon2id.CreateHash(data.Password, argon2id.DefaultParams)
	if err != nil {
		Println("failed to createHash:", err)
		return nil
	}

	new_user := User{
		Email:   data.Email,
		Pw_hash: pw_hash,
	}

	result := DB.Create(&new_user)
	if result.Error != nil {
		res := Jeti("user_already_exists", `<div hx-swap-oob="innerHTML:#status_register">User already exists</div>`, nil)
		return res
	}

	return Jet("login_page.html", nil)
}
