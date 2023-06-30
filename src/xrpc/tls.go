package xrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// LoadTLSConfig This is self-signed TLS
// Normal TLS use credentials.NewServerTLSFromFile
func LoadTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
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
	return NewTLSConfig(caPEMBlock, certPEMBlock, keyPEMBlock)
}

// NewTLSConfig This is self-signed TLS
func NewTLSConfig(ca, cert, key []byte) (*tls.Config, error) {
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

// LoadTLSClientConfig This is self-signed TLS
// Normal TLS use credentials.NewClientTLSFromFile
func LoadTLSClientConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
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
	return NewTLSClientConfig(caPEMBlock, certPEMBlock, keyPEMBlock)
}

// NewTLSClientConfig This is self-signed TLS
func NewTLSClientConfig(ca, cert, key []byte) (*tls.Config, error) {
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
