package config

import (
	"net-cat/models"
	"net-cat/pkg/bot"
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
