package main

import (
	"encoding/json"

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
		res = login_page(req)
	case "login_with_password":
		res = login_with_password(req)
	case "register_page":
		res = register_page(req)
	case "register":
		res = register(req)
	default:
		Println("invalid MessageType:", message_type)
		return nil
	}
	return res
}

func login_page(_ Req) Res {
	res := Res(`
<div hx-swap-oob="innerHTML:#ws">
	<div class="centerform">
		<h2>Login</h2>
		<form name="auth#login_with_password" ws-send>
			<label for="username">Username:</label><br>
			<input type="text" name="Email"><br><br>
			<label for="password">Password:</label><br>
			<input type="password" name="Password"><br><br>
			<button type="submit">Login</button><br><br>
			<div id="status_login"></div>
		</form>
	</div>
</div>`)
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

	invalid_email_or_password := Res(`<div hx-swap-oob="innerHTML:#status_login">Invalid Email or Password!</div>`)

	if len(users) != 1 {
		return invalid_email_or_password
	}

	user := users[0]
	match, err := argon2id.ComparePasswordAndHash(data.Password, user.Pw_hash)
	if err != nil {
		Println("failed to compare password with hash: %w", err)
		return nil
	}

	if !match {
		return invalid_email_or_password
	}

	return Home_page(req)
}

func register_page(_ Req) Res {
	res := Res(`
		<div hx-swap-oob="innerHTML:#ws">
			<div class="centerform">
				<h2>Register</h2>
				<form name="auth#register" ws-send>
					<label for="username">Username:</label><br>
					<input type="text" name="Email"><br><br>
					<label for="password">Password:</label><br>
					<input type="password" name="Password"><br><br>
					<button type="submit">Register</button><br><br>
					<div id="status_register"></div>
				</form>
			</div>
		</div>`)
	return res
}

func register(req Req) Res {
	var data struct {
		Email    string
		Password string
	}

	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		Println("failed to parse request: %w", err)
		return nil
	}

	pw_hash, err := argon2id.CreateHash(data.Password, argon2id.DefaultParams)
	if err != nil {
		Println("failed to createHash: %w", err)
		return nil
	}

	new_user := User{
		Email:   data.Email,
		Pw_hash: pw_hash,
	}

	result := DB.Create(&new_user)
	if result.Error != nil {
		res := Res(`<div hx-swap-oob="innerHTML:#status_register">User already exists</div>`)
		return res
	}

	return login_page(req)
}
