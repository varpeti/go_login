package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/CloudyKit/jet/v6"
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

func Println(a ...any) {
	pc, file, _, ok := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fileName := filepath.Base(file)
	if ok {
		prefix := fmt.Sprintf("%s%s%s#%s%s%s", magenta, fileName, def, magenta, funcName, def)
		a = append([]any{prefix}, a...)
	}
	log.Println(a...)
}

func Fatal(a ...any) {
	pc, file, _, ok := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fileName := filepath.Base(file)
	if ok {
		prefix := fmt.Sprintf("%s%s%s#%s%s%s", magenta, fileName, def, magenta, funcName, def)
		a = append([]any{prefix}, a...)
	}
	log.Fatal(a...)
}

func InitDb() *gorm.DB {
	var DB *gorm.DB
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		Fatal("failed to load .env:", err)
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		Fatal("failed to open postgres DB:", err)
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		Fatal("failed to migrate DB:", err)
	}

	return DB
}

type Req struct {
	Data        []byte
	MessageType []string
}
type Res []byte

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NextMessagType(req *Req) (string, error) {
	var ret string

	var meta struct {
		Headers struct {
			MessageType string `json:"HX-Trigger-Name"`
		} `json:"HEADERS"`
	}

	if req.MessageType == nil {
		err := json.Unmarshal(req.Data, &meta)
		if err != nil {
			return "", fmt.Errorf("failed to get MessageType from message: %w", err)
		}

		req.MessageType = strings.Split(meta.Headers.MessageType, "#")
	}

	if len(req.MessageType) < 1 {
		return "", fmt.Errorf("message_type is invalid (len < 1): %v", req.MessageType)
	}

	ret, req.MessageType = req.MessageType[0], req.MessageType[1:]
	return ret, nil
}

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.DevelopmentMode(true), // remove or set false in production
)

func Jet(template_name string, variables jet.VarMap) Res {
	view, err := views.GetTemplate(template_name)
	if err != nil {
		Println("Unexpected error when loading `", template_name, "`:", err)
		return nil
	}

	var buf bytes.Buffer
	err = view.Execute(&buf, variables, nil)
	if err != nil {
		Println("Unexpected error when rendering `", template_name, "`:", err)
		return nil
	}

	return buf.Bytes()
}

func Jeti(template_name string, template string, variables jet.VarMap) Res {
	view, err := views.Parse(template_name, template)
	if err != nil {
		Println("Unexpected error when parsing `", template_name, "`:", err)
		return nil
	}

	var buf bytes.Buffer
	err = view.Execute(&buf, variables, nil)
	if err != nil {
		Println("Unexpected error when rendering `", template_name, "`:", err)
		return nil
	}

	return buf.Bytes()
}
