package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

type ServerTransport struct {
	Listener          net.Listener
	MessageHandler    func(conn net.Conn, message []byte)
	DisconnectHandler func(conn net.Conn)
}

func NewServerTransport() *ServerTransport {
	return &ServerTransport{}
}

func (s *ServerTransport) Start(config *ServerConfig) (net.Listener, error) {
	certFile, err := filepath.Abs(config.CertFile)
	if err != nil {
		return nil, err
	}
	keyFile, err := filepath.Abs(config.KeyFile)
	if err != nil {
		return nil, err
	}
	caCertFile, err := filepath.Abs(config.CaCertFile)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	configTLS := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", config.Port), configTLS)
	if err != nil {
		return nil, err
	}

	s.Listener = listener
	return listener, nil
}

func (s *ServerTransport) SetMessageHandler(handler func(conn net.Conn, message []byte)) {
	s.MessageHandler = handler
}

func (s *ServerTransport) SetDisconnectHandler(handler func(conn net.Conn)) {
	s.DisconnectHandler = handler
}

func (s *ServerTransport) Close() {
	if s.Listener != nil {
		s.Listener.Close()
	}
}

func (s *ServerTransport) HandleConnection(conn net.Conn) {
	defer func() {
		if s.DisconnectHandler != nil {
			s.DisconnectHandler(conn)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected gracefully")
			} else {
				log.Printf("Client disconnected with error: %v", err)
			}
			return
		}

		// Pass the received data to the message handler
		if s.MessageHandler != nil {
			s.MessageHandler(conn, buf[:n])
		}
	}
}
