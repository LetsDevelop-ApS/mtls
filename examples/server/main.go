package main

import (
	"fmt"
	"log"

	"github.com/LetsDevelop-ApS/mtls"
)

type CustomMessage struct {
	*mtls.Message
	MyCustomProperty string
}

func NewCustomMessage(s string) *CustomMessage {
	return &CustomMessage{
		Message:          &mtls.Message{},
		MyCustomProperty: s,
	}
}

func main() {

	conf := mtls.DefaultPeerConfig().
		WithCertFile("certs/server/cert.pem").
		WithKeyFile("certs/server/key.key").
		WithCaCertFile("certs/server/ca.key").
		WithMessageTypes([]mtls.MessageInterface{
			&CustomMessage{},
		}).
		WithMessageHandler(func(msg mtls.MessageInterface) {
			switch msg.(type) {
			case *mtls.RegisterMessage:
				msg.Reply(mtls.NewRegisterSuccessMessage())
				msg.Reply(NewCustomMessage("Hello world!"))
			}
		}).
		Build()

	server := mtls.NewPeerTransport(conf)
	fmt.Println("I have peer ID", server.PeerID)
	err := server.Listen("127.0.0.1", 3001)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer server.Close()

	select {}
}
