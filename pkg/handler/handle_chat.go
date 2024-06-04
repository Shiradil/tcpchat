package handler

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/models"
	"strings"
	"sync"
	"time"
)

func UserChatHandler(mu *sync.Mutex, chatRoom *models.ChatRoom, user models.User) {
	conn := user.Conn

	mu.Lock()
	if chatRoom.BannedUsers[user.Name] {
		fmt.Fprintf(conn, "You are banned from this chat room\n")
		conn.Close()
		mu.Unlock()
		return
	}
	chatRoom.Users = append(chatRoom.Users, user)
	mu.Unlock()
	SendMessage(green+user.Name+" has joined the chat"+end, conn, mu)

	now := time.Now().Format("TCP SERVER 15:04")
	fmt.Fprintf(conn, "[%s][%s]:", now, user.Name)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		inputMessage := scanner.Text()

		fmt.Println(inputMessage)
		if checkMessage(inputMessage) {
			if inputMessage[0] == 47 {
				inputMessageSlice := strings.Split(inputMessage, " ")
				switch inputMessageSlice[0][1:] {
				case "leave":
					SendMessage(green+user.Name+" has left the chat"+end, conn, mu)
					LeaveChat(chatRoom, user, mu)
					return
				case "bot":
					if len(inputMessageSlice) < 2 {
						fmt.Fprintf(conn, "Bot: Please provide a command\n")
						continue
					}
					command := inputMessageSlice[1]
					app.Bot.HandleCommand(conn, command, mu, user.Name)
				default:
					fmt.Fprintf(conn, "Command doesn't exist\n")
				}
			} else {
				SendMessage(fmt.Sprintf("[%s][%s]: %s", now, user.Name, inputMessage), conn, mu)
			}
		}
		fmt.Fprintf(conn, "[%s][%s]:", now, user.Name)
	}

	defer func() {
		conn.Close()
		SendMessage(red+user.Name+" has left the chat"+end, conn, mu)
		LeaveChat(chatRoom, user, mu)
	}()
}

func SendMessage(message string, sender net.Conn, mu *sync.Mutex) {
	var users []models.User
	mu.Lock()
	for _, chatRoom := range app.ChatRoom {
		users = append(users, chatRoom.Users...)
	}
	mu.Unlock()

	for _, user := range users {
		if user.Conn != sender {
			fmt.Fprintf(user.Conn, "\n%s\n", message)
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(user.Conn, "[%s][%s]:", now, user.Name)
		}
	}
}

func LeaveChat(chatRoom *models.ChatRoom, user models.User, mu *sync.Mutex) {
	mu.Lock()
	for i, u := range chatRoom.Users {
		if u.Name == user.Name {
			chatRoom.Users = append(chatRoom.Users[:i], chatRoom.Users[i+1:]...)
			break
		}
	}
	mu.Unlock()
}
