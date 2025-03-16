package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math"
	"math/big"
	"net"
	"time"
)

const (
	ROOTCERT = "/etc/tlsguard/root.crt"
	ROOTKEY  = "/etc/tlsguard/root.key"
)

func GenClient() ([]byte, []byte, error) {
	rootcert, err := tls.LoadX509KeyPair(ROOTCERT, ROOTKEY)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	ca := rootcert.Leaf
	cakey := rootcert.PrivateKey

	clientkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	ser, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return []byte{}, []byte{}, err
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			Organization: ca.Issuer.Organization,
			CommonName:   "client",
		},
		SerialNumber: ser,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
	}
	clientcert, err := x509.CreateCertificate(rand.Reader, template, ca, &clientkey.PublicKey, cakey)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	cert := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: clientcert,
		},
	)

	key := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(clientkey),
		},
	)

	return cert, key, nil
}

func GenTls(addr string) (net.Listener, error) {
	crt, err := tls.LoadX509KeyPair(ROOTCERT, ROOTKEY)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AddCert(crt.Leaf)

	config := &tls.Config{
		Certificates: []tls.Certificate{crt},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return tls.Listen("tcp", addr, config)
}
