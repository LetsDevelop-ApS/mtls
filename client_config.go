package mtls

type ClientConfig struct {
	CertFile       string
	KeyFile        string
	CaCertFile     string
	ServerAddr     string
	ServerPort     int
	MessageHandler func([]byte)
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ServerPort: 3000,
		ServerAddr: "localhost",
	}
}

func (cc *ClientConfig) WithCertFile(cf string) *ClientConfig {
	cc.CertFile = cf
	return cc
}

func (cc *ClientConfig) WithKeyFile(kf string) *ClientConfig {
	cc.KeyFile = kf
	return cc
}

func (cc *ClientConfig) WithCaCertFile(ccf string) *ClientConfig {
	cc.CaCertFile = ccf
	return cc
}

func (cc *ClientConfig) WithServerAddr(serverAddr string) *ClientConfig {
	cc.ServerAddr = serverAddr
	return cc
}

func (cc *ClientConfig) WithServerPort(sp int) *ClientConfig {
	cc.ServerPort = sp
	return cc
}

func (cc *ClientConfig) WithMessageHandler(handler func([]byte)) *ClientConfig {
	cc.MessageHandler = handler
	return cc
}

func (cc *ClientConfig) Build() *ClientConfig {
	return cc
}
