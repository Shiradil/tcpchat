package handler

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net-cat/models"
	"strings"
	"sync"
	"time"
)

var encryptionKey []byte

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

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
			cryptedMessage, err := Encrypt(app.Key, message)
			if err != nil {
				fmt.Println(err)
			}
			acceptMessage(cryptedMessage, user)
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

func acceptMessage(crypteMessage string, user models.User) {
	msg, _ := Decrypt(app.Key, crypteMessage)
	fmt.Fprintf(user.Conn, "\n%s\n", msg)
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(user.Conn, "[%s][%s]:", now, user.Name)
}

// Encrypt function
func Encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)

	return fmt.Sprintf("%x:%x", iv, ciphertext[aes.BlockSize:]), nil
}

// Decrypt function
func Decrypt(key []byte, cryptoText string) (string, error) {
	parts := strings.SplitN(cryptoText, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid encrypted text format")
	}

	iv, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
