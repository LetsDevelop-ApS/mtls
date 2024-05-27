package main

import (
	"log"

	"github.com/LetsDevelop-ApS/mtls"
)

func main() {
	c := mtls.NewClient()
	config := mtls.DefaultClientConfig().
		WithCertFile("examples/client/certs/cert.pem").
		WithKeyFile("examples/client/certs/key.key").
		WithCaCertFile("examples/client/certs/ca.key").
		Build()
	err := c.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Example: Send data
	err = c.Send([]byte("Hello, server!\n"))
	if err != nil {
		log.Fatalf("Failed to send data: %v", err)
	}

	// Example: Receive data with a callback
	err = c.Receive(func(data []byte) {
		log.Printf("Received from server: %s", string(data))
	})
	if err != nil {
		log.Fatalf("Failed to receive data: %v", err)
	}

	c.Close()
}
