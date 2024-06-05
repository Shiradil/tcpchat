package bot

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Bot struct {
	Name string
}

var botCommands = map[string]func(net.Conn, *sync.Mutex, string){
	"time": sendTime,
	"joke": sendJoke,
	"roll": sendRandomNmbr,
}

func NewBot(name string) *Bot {
	return &Bot{Name: name}
}

func (b *Bot) HandleCommand(conn net.Conn, command string, mu *sync.Mutex, user string) {
	if cmd, exists := botCommands[command]; exists {
		cmd(conn, mu, user)
	} else {
		fmt.Fprintf(conn, "Bot: Unknown command\n")
	}
}

func sendTime(conn net.Conn, mu *sync.Mutex, user string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	mu.Lock()
	fmt.Fprintf(conn, "Bot: Current server time is %s\n", now)
	mu.Unlock()
}

func sendJoke(conn net.Conn, mu *sync.Mutex, user string) {
	joke := "Why do programmers prefer dark mode? Because light attracts bugs!\n"
	mu.Lock()
	fmt.Fprintf(conn, "Bot: %s\n", joke)
	mu.Unlock()
}

func sendRandomNmbr(conn net.Conn, mu *sync.Mutex, user string) {
	nmbr := rand.Intn(100)
	mu.Lock()
	fmt.Println("privet")
	fmt.Fprintf(conn, "%s has rolled the number [1-100]: %d\n", user, nmbr)
	mu.Unlock()
}
