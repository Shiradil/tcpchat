package config

import (
	"net-cat/bot"
	"net-cat/models"
)

type AppConfig struct {
	HostNumber string
	Users      []models.User
	ChatRoom   map[string]*models.ChatRoom
	ChatsName  []string
	AdminPanel models.AdminPanel
	Bot        bot.Bot
	Key        []byte
}
