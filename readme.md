# TCP Chat

## Objectives

This project aims to recreate NetCat in a Server-Client Architecture, allowing multiple clients to engage in a group chat. The server listens on a specified port for incoming connections, while clients can connect to the server, transmit messages, and interact with each other in a chat environment.

## Usage

```bash
$ go run cmd/main.go
Listening on the port :8080

$ go run . cmd/main.go --port=2525
Listening on port 127.0.0.1:2525
```

### Client Interaction

#### Server Response to Client Connection

```bash 
$ nc $IP $port
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```

#### Accept Connection with Non-Empty Name

```bash
$ nc $IP $port
[2020-01-20 16:03:43][client.name]:hello
[2020-01-20 16:03:46][client.name]:How are you?
[2020-01-20 16:04:10][client.name]:
NewClient has joined our chat...
[2020-01-20 16:04:15][client.name]:
[2020-01-20 16:04:32][NewClient]:Hi everyone!
[2020-01-20 16:04:32][client.name]:
[2020-01-20 16:04:35][NewClient]:How are you?
[2020-01-20 16:04:35][client.name]:great, and you?
[2020-01-20 16:04:41][client.name]:
[2020-01-20 16:04:44][NewClient]:good!
[2020-01-20 16:04:44][client.name]:
[2020-01-20 16:04:50][NewClient]:alright, see ya!
[2020-01-20 16:04:50][client.name]:bye-bye!
[2020-01-20 16:04:57][client.name]:
NewClient has left our chat...
[2020-01-20 16:04:59][client.name]:
```

#### Client Interaction

```bash
$ nc $IP $port
[2020-01-20 16:03:43][client.name]:hello
[2020-01-20 16:03:46][client.name]:How are you?
[2020-01-20 16:04:15][NewClient]:Hi everyone!
[2020-01-20 16:04:32][NewClient]:How are you?
[2020-01-20 16:04:35][NewClient]:
[2020-01-20 16:04:41][client.name]:great, and you?
[2020-01-20 16:04:41][NewClient]:good!
[2020-01-20 16:04:44][NewClient]:alright, see ya!
[2020-01-20 16:04:50][NewClient]:
[2020-01-20 16:04:57][client.name]:bye-bye!
[2020-01-20 16:04:57][NewClient]:^C
```
## Authors
