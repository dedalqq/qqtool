package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

func getSerialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	return serialNumber
}

func getTLSServerConn(c net.Conn) net.Conn {

	template := x509.Certificate{
		SerialNumber: getSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{"qqtool"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(10 * 365 * 24 * time.Hour),
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			{
				PrivateKey:  priv,
				Certificate: [][]byte{derBytes},
			},
		},
		InsecureSkipVerify: true,
	}

	c = tls.Server(c, tlsConfig)

	return c
}
