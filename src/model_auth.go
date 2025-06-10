package main

import (
	"context"
	"encoding/json"
	"go_login/templates"

	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email   string `gorm:"unique"`
	Pw_hash string
}

func Auth_handler(req Req) (Res, error) {
	var err error
	var res Res

	var message_type string
	message_type, err = NextMessagType(&req)
	if err != nil {
		return res, err
	}
	switch message_type {
	case "login_page":
		res, err = login_page(req)
	case "login_with_password":
		res, err = login_with_password(req)
	// case "register_page":
	// 	res, err = register_page(req)
	case "register":
		res, err = register(req)
	default:
		err = MyErrorf("invalid MessageType: %s", message_type)
	}

	return res, err
}

func login_page(_ Req) (Res, error) {
	var err error
	var res Res

	err = templates.Login_page().Render(context.Background(), &res)
	if err != nil {
		err = MyErrorf("failed to render template: %w", err)
		return Res{}, err
	}

	return res, err
}

func login_with_password(req Req) (Res, error) {
	var err error
	var res Res
	// res = Res(`<div hx-swap-oob="innerHTML:#status_login">Invalid Email or Password!</div>`)

	var data struct {
		Email    string
		Password string
	}

	err = json.Unmarshal(req.Data, &data)
	if err != nil {
		err = MyErrorf("failed to parse request: %w", err)
		return Res{}, err
	}

	var users []User
	result := DB.Where("email = ?", data.Email).Find(&users)
	if result.Error != nil {
		err = MyErrorf("db error: failed to get user by email:  %w", result.Error)
		return res, err
	}

	if len(users) != 1 {
		err = MyErrorf("user not found")
		return res, err
	}

	user := users[0]

	match, err := argon2id.ComparePasswordAndHash(data.Password, user.Pw_hash)
	if err != nil {
		err = MyErrorf("failed to compare password with hash: %w", err)
		return Res{}, err
	}

	if match {
		// res = Res(`<div hx-swap-oob="innerHTML:#status_login">Logging in...</div>`)
	}
	return res, err
}

func register(req Req) (Res, error) {
	var res Res
	var err error

	var data struct {
		Email    string
		Password string
	}

	err = json.Unmarshal(req.Data, &data)
	if err != nil {
		err = MyErrorf("failed to parse request: %w", err)
		return Res{}, err
	}

	pw_hash, err := argon2id.CreateHash(data.Password, argon2id.DefaultParams)
	if err != nil {
		err = MyErrorf("failed to createHash: %w", err)
		return Res{}, err
	}

	new_user := User{
		Email:   data.Email,
		Pw_hash: pw_hash,
	}

	result := DB.Create(&new_user)
	if result.Error != nil {
		err = MyErrorf("failed to Create user: %w", result.Error)
		// res = Res(`<div hx-swap-oob="innerHTML:#status_register">Invalid Email or Password</div>`)
		return res, err
	}

	// res = Res(`<div hx-swap-oob="innerHTML:#status_register">Registered!</div>`)

	return res, nil
}
