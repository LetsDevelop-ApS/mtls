package main

import (
	"log"

	"github.com/LetsDevelop-ApS/mtls"
)

func HandleClientMessage(clientID string, message []byte) {
	log.Printf("Business logic handling message from %s: %s", clientID, string(message))
	// Implement your business logic here

	// Example: Send a response back to the client
	err := mServer.SendMessage(clientID, []byte("Your response message here\n"))
	if err != nil {
		log.Printf("Failed to send response to %s: %v", clientID, err)
	}
}

var mServer *mtls.Server

func main() {
	config := mtls.DefaultServerConfig().
		WithCertFile("examples/server/certs/cert.pem").
		WithKeyFile("examples/server/certs/key.key").
		WithCaCertFile("examples/server/certs/ca.key").
		WithClientMessageHandler(HandleClientMessage).
		Build()

	mServer = mtls.NewServer(config)
	mServer.Start()
}
