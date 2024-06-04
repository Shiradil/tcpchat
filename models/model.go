package models

import (
	"net"
	"sync"
)

type User struct {
	Name string
	Conn net.Conn
}

type ChatRoom struct {
	Name         string
	Users        []User
	Host         User
	NumberOfuser int
	History      string
	Mu           *sync.Mutex
	BannedUsers  map[string]bool
}

type AdminPanel struct {
	Mu             *sync.Mutex
	ConnectedUsers []User
}
