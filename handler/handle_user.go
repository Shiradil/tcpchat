package handler

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/config"
	"net-cat/models"
	"strings"
	"sync"
)

func init() {
	getLogo()
}

var app *config.AppConfig

func NewHanlder(a *config.AppConfig) {
	app = a
}

func UserHandler(conn net.Conn, mu *sync.Mutex) {
	scanner := bufio.NewScanner(conn)
	fmt.Fprintf(conn, "%s", logo)

	name := getName(conn, mu)
	user := models.User{
		Name: name,
		Conn: conn,
	}

	mu.Lock()
	app.Users = append(app.Users, user)
	mu.Unlock()
	fmt.Fprintf(conn, "You have successfully joined\n")

	defer func() {
		conn.Close()
		deleteUser(name, mu)
	}()

	for scanner.Scan() {
		inputMessage := scanner.Text()

		fmt.Println(inputMessage)
		if checkMessage(inputMessage) {
			if inputMessage[0] == 47 {
				inputMessageSlice := strings.Split(inputMessage, " ")
				switch inputMessageSlice[0][1:] {
				case "create":
					if len(inputMessageSlice) != 2 {
						fmt.Fprintf(conn, "Please provide a name for the chat \n")
						continue
					}
					var muChat sync.Mutex
					chat := models.ChatRoom{
						Name:        inputMessageSlice[1],
						Mu:          &muChat,
						BannedUsers: make(map[string]bool),
					}
					mu.Lock()
					app.ChatRoom[chat.Name] = &chat
					app.ChatsName = append(app.ChatsName, chat.Name)
					mu.Unlock()
				case "list":
					mu.Lock()
					ChatsName := make([]string, len(app.ChatsName))
					copy(ChatsName, app.ChatsName)
					mu.Unlock()
					result := ""
					for i, name := range ChatsName {
						result += fmt.Sprintf("%d: %s\n", (i + 1), name)
					}
					fmt.Fprintf(conn, result)
				case "join":
					if len(inputMessageSlice) != 2 {
						fmt.Fprintf(conn, "Please provide a name of the chat \n")
						continue
					}
					var muChat *sync.Mutex
					mu.Lock()
					chat, ok := app.ChatRoom[inputMessageSlice[1]]
					if !ok {
						mu.Unlock()
						fmt.Fprintf(conn, "Chat doesn't exist\n")
						continue
					}
					muChat = chat.Mu
					mu.Unlock()

					UserChatHandler(muChat, chat, user)
				case "help":
					if name == "admin" {
						fmt.Fprintf(conn, "Admin commands:\n/ban [user_name] - Ban a user\n/kick [user_name] - Kick a user\n/stat - statistics")
					} else {
						fmt.Fprintf(conn, "Commands:\n/create [chat_name] - Create a new chat\n/list - List all chats\n/join [chat_name] - Join an existing chat\n/leave - Leave the current chat\n")
					}
				case "leave":
					fmt.Fprintf(conn, "You have not joined any chat\n")
				case "ban":
					if name == "admin" {
						if len(inputMessageSlice) != 2 {
							fmt.Fprintf(conn, "Please provide a name of the user to ban\n")
							continue
						}
						mu.Lock()
						banUser(inputMessageSlice[1], mu)
						mu.Unlock()
						fmt.Fprintf(conn, "User %s has been banned\n", inputMessageSlice[1])
					} else {
						fmt.Fprintf(conn, "You are not admin")
					}
				case "kick":
					if name == "admin" {
						if len(inputMessageSlice) != 2 {
							fmt.Fprintf(conn, "Please provide a name of the user to kick\n")
							continue
						}
						mu.Lock()
						kickUser(inputMessageSlice[1], mu)
						mu.Unlock()
						fmt.Fprintf(conn, "User %s has been kicked\n", inputMessageSlice[1])
					} else {
						fmt.Fprintf(conn, "You are not admin")
					}
				case "stat":
					if name == "admin" {
						mu.Lock()
						ChatsName := make([]string, len(app.ChatsName))
						copy(ChatsName, app.ChatsName)
						users := make([]string, len(app.Users))
						names := make([]string, len(app.Users))
						for _, user := range app.Users {
							names = append(names, user.Name)
						}
						copy(users, names)
						mu.Unlock()
						result := ""
						for i, name := range ChatsName {
							result += fmt.Sprintf("%d: %s\n", (i + 1), name)
						}
						resultu := ""
						fmt.Println(names)
						for i, name := range names {
							if i >= 2 {
								resultu += fmt.Sprintf("%d: %s\n", (i - 1), name)
							}
						}

						fmt.Fprintf(conn, "Chats:\n")
						fmt.Fprintf(conn, result)
						fmt.Fprintf(conn, "Users:\n")
						fmt.Fprintf(conn, resultu)
					} else {
						fmt.Fprintf(conn, "You are not admin")
					}
				default:
					fmt.Fprintf(conn, "Command doesn't exist\n")
				}
			}
		}
	}
}

func getName(conn net.Conn, mu *sync.Mutex) string {
	scanner := bufio.NewScanner(conn)
	var name string
	fmt.Fprint(conn, "[ENTER YOUR NAME]:")
	for scanner.Scan() {
		name = scanner.Text()
		if !checkMessage(name) {
			fmt.Fprintf(conn, "Wrong input!!! Enter your name again\n[ENTER YOUR NAME]:")
			continue
		} else if !nameIsBusy(name, mu) {
			fmt.Println()
			fmt.Fprintf(conn, "Name is busy!!! Enter your name again\n[ENTER YOUR NAME]:")
			continue
		}

		break
	}
	return name
}

func deleteUser(name string, mu *sync.Mutex) {
	mu.Lock()
	for i, user := range app.Users {
		if name == user.Name {
			app.Users = append(app.Users[:i], app.Users[i+1:]...)
		}
	}
	mu.Unlock()
}

func nameIsBusy(name string, mu *sync.Mutex) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range app.Users {
		if name == user.Name {
			return false
		}
	}
	return true
}

func checkMessage(str string) bool {
	if str == "" {
		return false
	}
	if str[0] == 27 {
		return false
	}

	for _, ch := range str {
		if ch != ' ' {
			return true
		}
	}

	return false
}

func banUser(name string, mu *sync.Mutex) {
	for _, chatRoom := range app.ChatRoom {
		if _, exists := chatRoom.BannedUsers[name]; !exists {
			chatRoom.BannedUsers[name] = true
			for i, user := range chatRoom.Users {
				if user.Name == name {
					chatRoom.Users = append(chatRoom.Users[:i], chatRoom.Users[i+1:]...)
					user.Conn.Close()
					break
				}
			}
		}
	}
}

func kickUser(name string, mu *sync.Mutex) {
	fmt.Println("not kicked")
	for _, chatRoom := range app.ChatRoom {
		for i, user := range chatRoom.Users {
			if user.Name == name {
				chatRoom.Users = append(chatRoom.Users[:i], chatRoom.Users[i+1:]...)
				user.Conn.Close()
				fmt.Println("kicked")
				break
			}
		}
	}
}
