package xrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// LoadServerTLSConfig This is self-signed TLS
// Normal TLS use credentials.NewServerTLSFromFile
func LoadServerTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	caPEMBlock, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPEMBlock, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return NewServerTLSConfig(caPEMBlock, certPEMBlock, keyPEMBlock)
}

// NewServerTLSConfig This is self-signed TLS
func NewServerTLSConfig(ca, cert, key []byte) (*tls.Config, error) {
	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("x509: certPool.AppendCertsFromPEM failed")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}, nil
}

// LoadClientTLSConfig This is self-signed TLS
// Normal TLS use credentials.NewClientTLSFromFile
func LoadClientTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	caPEMBlock, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPEMBlock, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return NewClientTLSConfig(caPEMBlock, certPEMBlock, keyPEMBlock)
}

// NewClientTLSConfig This is self-signed TLS
func NewClientTLSConfig(ca, cert, key []byte) (*tls.Config, error) {
	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("x509: certPool.AppendCertsFromPEM failed")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   "",
		RootCAs:      caCertPool,
	}, nil
}
