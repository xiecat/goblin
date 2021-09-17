package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func MakeDefautCert(certFile, keyFile string) error {

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	bits := 4096
	pk, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	serialNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}
	subject := pkix.Name{
		Organization:       []string{"default"},
		OrganizationalUnit: []string{"default"},
		CommonName:         "default",
	}

	rootTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(3650 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		//IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &pk.PublicKey, pk)
	if err != nil {
		return err
	}
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	certOut.Close()

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
	return nil
}
