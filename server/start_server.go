package server

import (
	"fmt"
	"log"
	"net"
	"net-cat/config"
	"net-cat/handler"
	"sync"
)

var (
	app *config.AppConfig
)

func NewServer(a *config.AppConfig) {
	app = a
}

func StartServer() {
	var mu sync.Mutex
	li, err := net.Listen("tcp", ":"+app.HostNumber)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Listening on the port :%s\n", app.HostNumber)

	handler.NewHanlder(app)

	defer li.Close()
	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}

		go handler.UserHandler(conn, &mu)
	}
}
