package mtls

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Transport            *ServerTransport
	Mu                   sync.Mutex
	Clients              map[string]net.Conn
	ClientMessageHandler func(clientID string, message []byte)
	Config               *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		Transport: NewServerTransport(),
		Clients:   make(map[string]net.Conn),
		Config:    config,
	}
}

func (s *Server) Start() {
	listener, err := s.Transport.Start(s.Config)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	s.ClientMessageHandler = s.Config.ClientMessageHandler
	s.Transport.SetMessageHandler(s.HandleMessage)
	s.Transport.SetDisconnectHandler(s.HandleDisconnect)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		clientID := s.AddClient(conn)
		log.Printf("Client connected with ID: %s", clientID)

		go s.Transport.HandleConnection(conn)
	}
}

func (s *Server) AddClient(conn net.Conn) string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	clientID := generateUniqueID()
	s.Clients[clientID] = conn
	return clientID
}

func (s *Server) HandleDisconnect(conn net.Conn) {
	clientID := s.GetClientID(conn)
	if clientID != "" {
		log.Printf("Client with ID %s disconnected", clientID)
	}
	conn.Close()
	s.RemoveClient(clientID)
}

func (s *Server) HandleMessage(conn net.Conn, message []byte) {
	clientID := s.GetClientID(conn)
	if clientID != "" {
		log.Printf("Received from %s: %s", clientID, string(message))
		if s.ClientMessageHandler != nil {
			s.ClientMessageHandler(clientID, message)
		}
	}
}

func (s *Server) GetClientID(conn net.Conn) string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for clientID, clientConn := range s.Clients {
		if clientConn == conn {
			return clientID
		}
	}
	return ""
}

func (s *Server) RemoveClient(clientID string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Clients, clientID)
}

func (s *Server) SendMessage(clientID string, message []byte) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	conn, exists := s.Clients[clientID]
	if !exists {
		return fmt.Errorf("client %s not found", clientID)
	}

	_, err := conn.Write(message)
	return err
}

func generateUniqueID() string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return hex.EncodeToString(hash.Sum(nil))
}
