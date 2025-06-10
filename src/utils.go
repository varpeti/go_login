package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	escape  = "\x1b["
	magenta = escape + "35m"
	def     = escape + "0m"
)

func MyErrorf(format string, a ...any) error {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf(format, a...)
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := filepath.Base(file)
	prefix := fmt.Sprintf("%s%s%s#%s%s%s", magenta, fileName, def, magenta, funcName, def)

	a = append([]any{prefix}, a...)

	return fmt.Errorf("%s "+format, a...)
}

func InitDb() (*gorm.DB, error) {
	var DB *gorm.DB
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		err = MyErrorf("failed to load .env: %w", err)
		return nil, err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		err = MyErrorf("failed to open postgres DB: %w", err)
		return nil, err
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		err = MyErrorf("failed to migrate DB: %w", err)
		return nil, err
	}

	return DB, nil
}

type Req struct {
	Data        []byte
	MessageType []string
}
type Res struct {
	bytes.Buffer
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NextMessagType(req *Req) (string, error) {
	var ret string
	var err error

	var meta struct {
		Headers struct {
			MessageType string `json:"HX-Trigger-Name"`
		} `json:"HEADERS"`
	}

	if req.MessageType == nil {
		err = json.Unmarshal(req.Data, &meta)
		if err != nil {
			err = MyErrorf("failed to get MessageType from message: %w", err)
			return ret, err
		}

		req.MessageType = strings.Split(meta.Headers.MessageType, "#")
	}

	if len(req.MessageType) < 1 {
		err = MyErrorf("message_type is invalid (len < 1): %v", req.MessageType)
		return ret, err
	}

	ret, req.MessageType = req.MessageType[0], req.MessageType[1:]

	return ret, err
}
