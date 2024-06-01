package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
)

type Message struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

type Server struct {
	Host    string
	Port    string
	Clients []*Client
	mu      sync.Mutex
}

type Client struct {
	Conn     net.Conn
	ID       uuid.UUID
	Nickname string
}

type Config struct {
	Host string
	Port string
}

func New(config *Config) *Server {
	return &Server{
		Host:    config.Host,
		Port:    config.Port,
		Clients: make([]*Client, 0),
	}
}

func (server *Server) Broadcast(message string, sender string) {
	msg := Message{
		Content: message,
		Sender:  sender,
	}

	encodedMsg, err := json.Marshal(msg)

	if err != nil {
		log.Fatal(err)

		return
	}

	encodedMsg = append(encodedMsg, '\n')

	server.mu.Lock()
	defer server.mu.Unlock()

	for _, client := range server.Clients {
		_, err = client.Conn.Write([]byte(encodedMsg))

		if err != nil {
			log.Println("Error broadcasting message to client:", err)
			client.Conn.Close()
		}
	}
}

func (server *Server) RemoveClient(client *Client) {
	fmt.Println("Client disconnected:", client.Nickname)
	server.mu.Lock()
	defer server.mu.Unlock()

	for i, c := range server.Clients {
		if c.ID == client.ID {
			server.Clients = append(server.Clients[:i], server.Clients[i+1:]...)
			break
		}
	}
}

func (server *Server) OnClientConnected(client *Client) {
	defer func() {
		client.Conn.Close()
		server.RemoveClient(client)
	}()

	reader := bufio.NewReader(client.Conn)

	for {
		reply, err := reader.ReadString('\n')
		if err != nil {
			client.Conn.Close()
			break
		}

		if strings.HasPrefix(reply, "/nick") {
			client.Nickname = strings.TrimPrefix(reply, "/nick")
			client.Nickname = strings.TrimSpace(client.Nickname)
			client.Nickname = strings.Replace(client.Nickname, "\n", "", -1)

			fmt.Println("Client connected with nickname:", client.Nickname)

			successMsg := Message{
				Sender:  "server",
				Content: "nickname set",
			}

			encodedMsg, err := json.Marshal(successMsg)
			if err != nil {
				log.Fatal(err)
			}

			_, err = client.Conn.Write(append(encodedMsg, '\n'))
			if err != nil {
				log.Fatal(err)
			}

			server.Broadcast(fmt.Sprintf("%s connected", client.Nickname), "server")

			server.mu.Lock()
			server.Clients = append(server.Clients, client)
			server.mu.Unlock()
			continue
		}

		reply = strings.Replace(reply, "\n", "", -1)
		server.Broadcast(reply, client.Nickname)
		fmt.Printf("%s: %s\n", client.Nickname, reply)
	}
}

func (server *Server) Run() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", fmt.Sprintf("%s:%s", server.Host, server.Port), config)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			Conn: conn,
			ID:   id,
		}

		go server.OnClientConnected(client)
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := New(&Config{
		Host: "0.0.0.0",
		Port: port,
	})

	fmt.Printf("Server running on %s:%s\n", server.Host, server.Port)
	server.Run()
}
