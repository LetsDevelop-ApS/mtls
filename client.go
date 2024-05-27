package mtls

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Client struct {
	transport *ClientTransport
	stopChan  chan os.Signal
}

func NewClient() *Client {
	client := &Client{
		transport: NewClientTransport(),
		stopChan:  make(chan os.Signal, 1),
	}

	// Set up signal handler for graceful shutdown
	signal.Notify(client.stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-client.stopChan
		log.Println("Shutting down client...")
		client.Close()
		os.Exit(0)
	}()

	return client
}

func (c *Client) Connect(cc *ClientConfig) error {
	return c.transport.Connect(cc)
}

func (c *Client) Send(data []byte) error {
	return c.transport.Send(data)
}

func (c *Client) Receive(callback func([]byte)) error {
	return c.transport.Receive(callback)
}

func (c *Client) Close() {
	c.transport.Close()
}
