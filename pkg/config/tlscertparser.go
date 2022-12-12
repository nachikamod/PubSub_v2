package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

func ParseTransportCreds() (*credentials.TransportCredentials, error) {
	caPem, err := ioutil.ReadFile("ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair("server-cert.pem", "server-key.pem")
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	creds := credentials.NewTLS(tlsConfig)

	return &creds, nil
}
