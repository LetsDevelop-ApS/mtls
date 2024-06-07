package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
)

func LoadTLSConfig(certFile, keyFile, caCertFile string) (*tls.Config, *x509.CertPool, error) {
	certFile, err := filepath.Abs(certFile)
	if err != nil {
		return nil, nil, err
	}
	keyFile, err = filepath.Abs(keyFile)
	if err != nil {
		return nil, nil, err
	}
	caCertFile, err = filepath.Abs(caCertFile)
	if err != nil {
		return nil, nil, err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, err
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return config, caCertPool, nil
}
