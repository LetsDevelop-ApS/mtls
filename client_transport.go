package mtls

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
)

type ClientTransport struct {
	conn net.Conn
}

func NewClientTransport() *ClientTransport {
	return &ClientTransport{}
}

func (c *ClientTransport) Connect(cc *ClientConfig) error {
	certFile, err := filepath.Abs(cc.CertFile)
	if err != nil {
		return err
	}
	keyFile, err := filepath.Abs(cc.KeyFile)
	if err != nil {
		return err
	}
	caCertFile, err := filepath.Abs(cc.CaCertFile)
	if err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", cc.ServerAddr, cc.ServerPort), config)
	if err != nil {
		return err
	}
	c.conn = conn

	log.Printf("Connected to server at %s:%d", cc.ServerAddr, cc.ServerPort)
	return nil
}

func (c *ClientTransport) Send(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *ClientTransport) Receive(callback func([]byte)) error {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		callback(line)
	}
}

func (c *ClientTransport) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
