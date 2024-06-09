package mtls

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type PeerTransport struct {
	listener    net.Listener
	config      *PeerConfig
	connections map[string]*PeerConnection
	mutex       sync.Mutex
	PeerID      string
}

func NewPeerTransport(config *PeerConfig) *PeerTransport {
	peerID := getPeerID()
	gob.Register(&Message{})
	gob.Register(&RegisterMessage{})
	gob.Register(&RegisterSuccessMessage{})
	for _, t := range config.MessageTypes {
		gob.Register(t)
	}
	return &PeerTransport{
		connections: make(map[string]*PeerConnection),
		config:      config,
		PeerID:      peerID,
	}
}

type MessageInterface interface {
	GetConn() *PeerConnection
	SetConn(*PeerConnection)
	GetSenderID() string
	SetSenderID(string)
	Reply(MessageInterface) error
}

type Message struct {
	conn     *PeerConnection
	SenderID string
}

func (m Message) GetConn() *PeerConnection {
	if m.conn != nil {
		return m.conn
	}
	return &PeerConnection{}
}

func (m *Message) SetConn(conn *PeerConnection) {
	m.conn = conn
}

func (m *Message) GetSenderID() string {
	return m.SenderID
}

func (m *Message) SetSenderID(ID string) {
	m.SenderID = ID
}

func (m *Message) Reply(msg MessageInterface) error {
	return m.conn.Send(msg)
}

type RegisterMessage struct {
	*Message
}

func NewRegisterMessage() *RegisterMessage {
	return &RegisterMessage{
		Message: &Message{},
	}
}

type RegisterSuccessMessage struct {
	*Message
}

func NewRegisterSuccessMessage() *RegisterSuccessMessage {
	return &RegisterSuccessMessage{
		Message: &Message{},
	}
}

type PeerConnection struct {
	net.Conn
	Transport *PeerTransport
	Inbound   bool
	Outbound  bool
}

func (pc *PeerConnection) Send(msg MessageInterface) error {
	msg.SetSenderID(pc.Transport.PeerID)
	return gob.NewEncoder(pc.Conn).Encode(&msg)
}

func (p *PeerTransport) Listen(address string, port uint16) error {
	tlsConfig, caCertPool, err := LoadTLSConfig(p.config.CertFile, p.config.KeyFile, p.config.CaCertFile)
	if err != nil {
		return err
	}

	tlsConfig.ClientCAs = caCertPool
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	listener, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", address, port), tlsConfig)
	if err != nil {
		return err
	}

	p.listener = listener
	go p.acceptConnections()
	return nil
}

func (p *PeerTransport) Connect(peerAddr string, peerPort int) (*PeerConnection, error) {
	tlsConfig, caCertPool, err := LoadTLSConfig(p.config.CertFile, p.config.KeyFile, p.config.CaCertFile)
	if err != nil {
		return nil, err
	}

	tlsConfig.RootCAs = caCertPool

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", peerAddr, peerPort), tlsConfig)
	if err != nil {
		return nil, err
	}
	peerConn := &PeerConnection{
		Conn:      conn,
		Transport: p,
		Outbound:  true,
		Inbound:   false,
	}

	p.mutex.Lock()
	p.connections[peerConn.RemoteAddr().String()] = peerConn
	p.mutex.Unlock()

	go p.handleConnection(peerConn)
	return peerConn, nil
}

func (p *PeerTransport) acceptConnections() error {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			return err
		}
		peerConn := &PeerConnection{
			Conn:      conn,
			Transport: p,
			Outbound:  false,
			Inbound:   true,
		}
		p.mutex.Lock()
		p.connections[peerConn.RemoteAddr().String()] = peerConn
		p.mutex.Unlock()
		go p.handleConnection(peerConn)
	}
}

func (p *PeerTransport) handleConnection(conn *PeerConnection) error {
	defer func() {
		p.mutex.Lock()
		delete(p.connections, conn.RemoteAddr().String())
		p.mutex.Unlock()
		if conn.Inbound {
			if p.config.InboundDisconnectionHandler != nil {
				p.config.InboundDisconnectionHandler(conn)
			}
		} else if conn.Outbound {
			if p.config.OutboundDisconnectionHandler != nil {
				p.config.OutboundDisconnectionHandler(conn)
			}
		}
	}()

	for {
		var msg MessageInterface
		err := gob.NewDecoder(conn.Conn).Decode(&msg)
		if err != nil {
			if err == io.EOF {
				if conn.Inbound {
					log.Printf("Peer disconnected gracefully")
				} else if conn.Outbound {
					log.Printf("Peer terminated connection gracefully")
				}
			} else {
				if conn.Inbound {
					log.Printf("Peer disconnected with error: %v", err)
				} else if conn.Outbound {
					log.Printf("Peer terminated connection with error: %v", err)
				}
			}
			return nil
		}

		if p.config.MessageHandler != nil {
			msg.SetConn(conn)
			p.config.MessageHandler(msg)
		}
	}
}

func (p *PeerTransport) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.listener != nil {
		p.listener.Close()
	}
	for _, conn := range p.connections {
		conn.Close()
	}
	p.connections = make(map[string]*PeerConnection)
}

func getPeerID() string {
	macAddress := getMacAddress()
	hash := sha256.New()
	hash.Write([]byte(macAddress))
	newPeerID := hex.EncodeToString(hash.Sum(nil))
	return newPeerID
}

func getMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}
	for _, i := range interfaces {
		mac := i.HardwareAddr.String()
		if mac != "" {
			return mac
		}
	}
	return "unknown"
}
