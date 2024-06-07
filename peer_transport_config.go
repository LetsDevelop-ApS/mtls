package mtls

type PeerConfig struct {
	CertFile                     string
	KeyFile                      string
	CaCertFile                   string
	PeerID                       string
	MessageTypes                 []MessageInterface
	IncomingPort                 int
	MessageHandler               func(message MessageInterface)
	InboundConnectionHandler     func(conn *PeerConnection)
	InboundDisconnectionHandler  func(conn *PeerConnection)
	OutboundConnectionHandler    func(conn *PeerConnection)
	OutboundDisconnectionHandler func(conn *PeerConnection)
}

func DefaultPeerConfig() *PeerConfig {
	return &PeerConfig{}
}

func (pc *PeerConfig) WithPeerID(peerID string) *PeerConfig {
	pc.PeerID = peerID
	return pc
}

func (pc *PeerConfig) WithCertFile(cf string) *PeerConfig {
	pc.CertFile = cf
	return pc
}

func (pc *PeerConfig) WithKeyFile(kf string) *PeerConfig {
	pc.KeyFile = kf
	return pc
}

func (pc *PeerConfig) WithCaCertFile(ccf string) *PeerConfig {
	pc.CaCertFile = ccf
	return pc
}

func (pc *PeerConfig) WithMessageTypes(ct []MessageInterface) *PeerConfig {
	pc.MessageTypes = ct
	return pc
}

func (pc *PeerConfig) WithMessageHandler(handler func(message MessageInterface)) *PeerConfig {
	pc.MessageHandler = handler
	return pc
}

func (pc *PeerConfig) WithInboundConnectionHandler(handler func(conn *PeerConnection)) *PeerConfig {
	pc.InboundConnectionHandler = handler
	return pc
}

func (pc *PeerConfig) WithInboundDisconnectionHandler(handler func(conn *PeerConnection)) *PeerConfig {
	pc.InboundDisconnectionHandler = handler
	return pc
}

func (pc *PeerConfig) WithOutboundConnectionHandler(handler func(conn *PeerConnection)) *PeerConfig {
	pc.OutboundConnectionHandler = handler
	return pc
}

func (pc *PeerConfig) WithOutboundDisconnectionHandler(handler func(conn *PeerConnection)) *PeerConfig {
	pc.OutboundDisconnectionHandler = handler
	return pc
}

func (pc *PeerConfig) Build() *PeerConfig {
	return pc
}
