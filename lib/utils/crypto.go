package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"os/exec"

	config "github.com/sdslabs/katana/configs"
)

// MD5 encodes string to hexadecimal of MD5 checksum.
func MD5(str string) string {
	return hex.EncodeToString(MD5Bytes(str))
}

// MD5Bytes encodes string to MD5 checksum.
func MD5Bytes(str string) []byte {
	m := md5.New()
	_, _ = m.Write([]byte(str))
	return m.Sum(nil)
}

// Base64Encode encodes string to base64.
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Generating certificate from domain name
func GenerateCert(domain string, basePath string) error {
	// Generate ca.key in harbor directory
	cmd := "openssl genrsa -out " + basePath + "/ca.key 4096"
	out := exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	// Generate ca.crt
	cmd = "openssl req -x509 -new -nodes -sha512 -days 3650 -subj '/C=IN/ST=Delhi/L=Delhi/O=Katana/CN=" + config.KatanaConfig.Harbor.Hostname + "' -key ca.key -out " + basePath + "/ca.crt"
	out = exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	// Generate private key
	cmd = "openssl genrsa -out " + basePath + "/" + config.KatanaConfig.Harbor.Hostname + ".key 4096"
	out = exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	// Generate certificate signing request
	cmd = "openssl req -sha512 -new -subj '/C=IN/ST=Delhi/L=Delhi/O=Katana/CN=" + config.KatanaConfig.Harbor.Hostname + "' -key " + basePath + "/" + config.KatanaConfig.Harbor.Hostname + ".key -out " + basePath + "/" + config.KatanaConfig.Harbor.Hostname + ".csr"
	out = exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	// Generate v3.ext file
	cmd = "echo 'authorityKeyIdentifier=keyid,issuer\nbasicConstraints=CA:FALSE\nkeyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment\nextendedKeyUsage = serverAuth\nsubjectAltName = @alt_names\n[alt_names]\nDNS.1=" + config.KatanaConfig.Harbor.Hostname + "' > " + basePath + "/v3.ext"
	out = exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	// Generate certificate
	cmd = "openssl x509 -req -sha512 -days 3650 -extfile " + basePath + "/v3.ext -CA " + basePath + "/ca.crt -CAkey " + basePath + "/ca.key -CAcreateserial -in " + basePath + "/" + config.KatanaConfig.Harbor.Hostname + ".csr -out " + basePath + "/" + config.KatanaConfig.Harbor.Hostname + ".crt"
	out = exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	return nil
}
