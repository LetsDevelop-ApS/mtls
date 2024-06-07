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

func main() {

	conf := mtls.DefaultPeerConfig().
		WithCertFile("certs/client/cert.pem").
		WithKeyFile("certs/client/key.key").
		WithCaCertFile("certs/client/ca.key").
		WithMessageTypes([]mtls.MessageInterface{
			&CustomMessage{},
		}).
		WithMessageHandler(func(m mtls.MessageInterface) {
			switch msg := m.(type) {
			case *mtls.RegisterSuccessMessage:
				fmt.Println("Successfully registered with peer", msg.GetSenderID())
			case *CustomMessage:
				fmt.Printf("Got custom message from %s: %s", msg.GetSenderID(), msg.MyCustomProperty)
			}

		}).
		Build()

	client := mtls.NewPeerTransport(conf)
	fmt.Println("I have peer ID", client.PeerID)
	conn, err := client.Connect("localhost", 3001)
	if err != nil {
		log.Fatalf("Peer2 (Client) failed to connect to Peer1 (Server): %v", err)
	}

	// Send initial messages
	err = conn.Send(mtls.NewRegisterMessage())
	if err != nil {
		fmt.Println(err)
	}

	// Run indefinitely, handling connections and messages
	select {}
}
