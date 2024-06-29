# TCPChat

## Created by
1. Shirbayev Adilzhan
2. Abay Aliyev 

## Description
(Additional functionality has been addded)
This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

## Usage
1. Run the cmd directory ```go run ./cmd```
2. Create new terminal, and connect to tcp server using nc. Ex. ```nc localhost 8080```
3. Write your name
4. Connect to the chat using ```/join [CHAT NAME]```

## Commands
1. ```/help``` - to get information of commands
2. ```/create [CHAT NAME]``` - create new chat
3. ```/join [CHAT NAME]``` - join to created chats
4. ```/list``` - list of all chats

Bot commands(only in chat):
1. ```/bot roll``` - roll a random number between 1-100
2. ```/bot joke``` - static joke (not useful)
3. ```/bot time``` - to get current server time