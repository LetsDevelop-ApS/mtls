package mtls

type ServerConfig struct {
	CertFile             string
	KeyFile              string
	CaCertFile           string
	Port                 int
	ClientMessageHandler func(clientID string, message []byte)
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: 3000,
	}
}

func (sc *ServerConfig) WithCertFile(cf string) *ServerConfig {
	sc.CertFile = cf
	return sc
}

func (sc *ServerConfig) WithKeyFile(sk string) *ServerConfig {
	sc.KeyFile = sk
	return sc
}

func (sc *ServerConfig) WithCaCertFile(ccf string) *ServerConfig {
	sc.CaCertFile = ccf
	return sc
}

func (sc *ServerConfig) WithPort(p int) *ServerConfig {
	sc.Port = p
	return sc
}

func (sc *ServerConfig) WithClientMessageHandler(h func(clientID string, message []byte)) *ServerConfig {
	sc.ClientMessageHandler = h
	return sc
}

func (sc *ServerConfig) Build() *ServerConfig {
	return sc
}
