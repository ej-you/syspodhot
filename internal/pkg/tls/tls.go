package tls

import (
	"crypto/tls"
	"net"
)

// WrapWithTLS wraps given net.Listener with given TLS-certificate.
func WrapWithTLS(listener net.Listener, certFile, keyFile string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	return tls.NewListener(listener, tlsConfig), nil
}
