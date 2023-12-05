package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/xdg-go/pbkdf2"
	"golang.org/x/crypto/bcrypt"
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

func GenerateCerts(domain string, basePath string) error {
	// Generate ca.key in harbor directory
	//add -traditional flag to make it run on minikube
	cmd := "openssl genrsa -out " + basePath + "/ca.key 4096"
	if err := RunCommand(cmd); err != nil {
		return err
	}

	// Generate ca.crt
	cmd = "openssl req -x509 -new -nodes -sha512 -days 3650 -subj '/C=IN/ST=Delhi/L=Delhi/O=Katana/CN=" + domain + "' -key " + basePath + "/ca.key -out " + basePath + "/ca.crt"
	if err := RunCommand(cmd); err != nil {
		return err
	}

	// Generate private key
	//add -traditional flag to make it run on minikube
	cmd = "openssl genrsa -out " + basePath + "/" + domain + ".key 4096"
	if err := RunCommand(cmd); err != nil {
		return err
	}

	// Generate certificate signing request
	cmd = "openssl req -sha512 -new -subj '/C=IN/ST=Delhi/L=Delhi/O=Katana/CN=" + domain + "' -key " + basePath + "/" + domain + ".key -out " + basePath + "/" + domain + ".csr"
	if err := RunCommand(cmd); err != nil {
		return err
	}

	// Generate v3.ext file
	cmd = "echo 'authorityKeyIdentifier=keyid,issuer\nbasicConstraints=CA:FALSE\nkeyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment\nextendedKeyUsage = serverAuth\nsubjectAltName = @alt_names\n[alt_names]\nDNS.1=" + domain + "' > " + basePath + "/v3.ext"
	if err := RunCommand(cmd); err != nil {
		return err
	}

	// Generate certificate
	cmd = "openssl x509 -req -sha512 -days 3650 -extfile " + basePath + "/v3.ext -CA " + basePath + "/ca.crt -CAkey " + basePath + "/ca.key -CAcreateserial -in " + basePath + "/" + domain + ".csr -out " + basePath + "/" + domain + ".crt"
	if err := RunCommand(cmd); err != nil {
		return err
	}

	return nil
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
