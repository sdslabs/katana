package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/xdg-go/pbkdf2"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"

	"strings"
	"time"
)

// MD5 encodes string to hexadecimal of MD5 checksum.
func MD5(str string) string {
	m := md5.New()
	_, _ = m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// Base64Encode encodes string to base64.
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// V3Ext represents a v3.ext file
type V3Ext struct {
	AuthorityKeyIdentifier string
	BasicConstraintsValid bool
	IsCA       bool
	KeyUsage               string
	ExtKeyUsage            string
	DNSNames               []string
}

func GenerateCerts(domain string, basePath string) error {
	basePath += "/"
	// Generate a new private key for the CA
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// Set up the certificate template for the CA
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"SDSLabs"},
			Country:       []string{"IN"},
			Province:      []string{"Delhi"},
			Locality:      []string{"Delhi"},
			StreetAddress: []string{"smoking jawahar"},
			PostalCode:    []string{"110080"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create the CA certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return err
	}

	// Save the CA private key
	keyFile, err := os.Create(basePath + "ca.key")
	if err != nil {
		return err
	}
	defer keyFile.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(caPrivateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	_, err = keyFile.Write(privPEM)
	if err != nil {
		return err
	}

	// Save the CA certificate
	certFile, err := os.Create(basePath + "ca.crt")
	if err != nil {
		return err
	}
	defer certFile.Close()

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	_, err = certFile.Write(certPEM)
	if err != nil {
		return err
	}

	// Generate a new private key for the server
	serverPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// Set up the CSR template for the server
	csrTemplate := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: domain,
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames:           []string{domain},
	}

	// Create the CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, serverPrivateKey)
	if err != nil {
		return err
	}

	// Save the CSR
	csrFile, err := os.Create(basePath + domain + ".csr")
	if err != nil {
		return err
	}
	defer csrFile.Close()

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	_, err = csrFile.Write(csrPEM)
	if err != nil {
		return err
	}

	// Save the server private key
	serverKeyFile, err := os.Create(basePath + domain + ".key")
	if err != nil {
		return err
	}
	defer serverKeyFile.Close()

	serverPrivBytes := x509.MarshalPKCS1PrivateKey(serverPrivateKey)
	serverPrivPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: serverPrivBytes})
	_, err = serverKeyFile.Write(serverPrivPEM)
	if err != nil {
		return err
	}

	// Define your v3.ext
	v3ext := V3Ext{
		AuthorityKeyIdentifier: "keyid,issuer",
		BasicConstraintsValid: true,
		IsCA:       false,
		KeyUsage:               "digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment",
		ExtKeyUsage:            "serverAuth",
		DNSNames:               []string{"harbor.katana.local"},
	}

	// Set up the certificate template for the server
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: domain,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		KeyUsage:  keyUsage(v3ext.KeyUsage),
		ExtKeyUsage: []x509.ExtKeyUsage{
			extKeyUsage(v3ext.ExtKeyUsage),
		},
		DNSNames: v3ext.DNSNames,
		BasicConstraintsValid: v3ext.BasicConstraintsValid,
		IsCA: v3ext.IsCA,
	}

	// Create the server certificate
	serverBytes, err := x509.CreateCertificate(rand.Reader, serverTemplate, caTemplate, &serverPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return err
	}

	// Save the server certificate
	serverCertFile, err := os.Create(basePath + domain + ".crt")
	if err != nil {
		return err
	}
	defer serverCertFile.Close()

	serverCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverBytes,
	})

	_, err = serverCertFile.Write(serverCertPEM)
	if err != nil {
		return err
	}
	return nil
}

func keyUsage(s string) x509.KeyUsage {
	var ku x509.KeyUsage
	if strings.Contains(s, "digitalSignature") {
		ku |= x509.KeyUsageDigitalSignature
	}
	if strings.Contains(s, "nonRepudiation") {
		ku |= x509.KeyUsageContentCommitment
	}
	if strings.Contains(s, "keyEncipherment") {
		ku |= x509.KeyUsageKeyEncipherment
	}
	if strings.Contains(s, "dataEncipherment") {
		ku |= x509.KeyUsageDataEncipherment
	}
	return ku
}

func extKeyUsage(s string) x509.ExtKeyUsage {
	var eku x509.ExtKeyUsage
	if strings.Contains(s, "serverAuth") {
		eku = x509.ExtKeyUsageServerAuth
	}
	return eku
}

func HashPassword(password string) (string, error) {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashWithPassword(hashedPassword, password string) bool {
	hash := []byte(hashedPassword)
	pass := []byte(password)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	return err == nil
}

// EncodePassword encodes password using PBKDF2 SHA256 with given salt.
func EncodePassword(password, salt string) string {
	newPasswd := pbkdf2.Key([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return fmt.Sprintf("%x", newPasswd)
}

func SHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", hash)
}
