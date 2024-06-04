package main

import (
	"log"
	"net-cat/models"
	"net-cat/pkg/bot"
	"net-cat/pkg/config"
	"net-cat/pkg/handler"
	"net-cat/pkg/server"
	"os"
	"strconv"
)

func main() {
	number := getNumberLocalhost()
	if number == "" {
		log.Fatal("please write in the right format\n Example: go run . 8080")
	}

	key, _ := handler.GenerateKey()
	app := config.AppConfig{
		HostNumber: number,
		ChatRoom:   make(map[string]*models.ChatRoom),
		ChatsName:  make([]string, 0, 3),
		Bot:        *bot.NewBot("Chatbot"),
		Key:        key,
	}

	server.NewServer(&app)
	server.StartServer()
}

func getNumberLocalhost() string {
	args := os.Args[1:]
	if len(args) == 0 {
		return "8080"
	} else if len(args) == 1 {
		if len(args[0]) != 4 {
			return ""
		}
		_, err := strconv.Atoi(args[0])
		if err != nil {
			return ""
		}
		return args[0]
	}
	return ""
}
