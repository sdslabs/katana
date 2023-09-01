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
	
	openSSLabove3,err :=CheckOpenSSLVersion()
	if err != nil {	
		return err
	}
	// Generate ca.key in harbor directory
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
	cmd = "openssl genrsa -out " + basePath + "/" + domain + ".key 4096"
	if err := RunCommand(cmd); err != nil {
		return err
	}

	if(openSSLabove3){
		cmd := "openssl rsa -in " + basePath+"/"+domain + ".key -out "+ basePath +"/"+domain + ".key -traditional"
		if err := RunCommand(cmd); err != nil {
			return err
		}
	}

	if(openSSLabove3){
		cmd := "openssl rsa -in " + basePath+"/"+domain + ".key -out "+ basePath +"/"+domain + ".key -traditional"
		if err := RunCommand(cmd); err != nil {
			return err
		}
	}

	if(openSSLabove3){
		cmd := "openssl rsa -in " + basePath+"/"+domain + ".key -out "+ basePath +"/"+domain + ".key -traditional"
		if err := RunCommand(cmd); err != nil {
			return err
		}
	}

	if(openSSLabove3){
		cmd := "openssl rsa -in " + basePath+"/"+domain + ".key -out "+ basePath +"/"+domain + ".key -traditional"
		if err := RunCommand(cmd); err != nil {
			return err
		}
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

// [DO NOT REMOVE....WIP....]

// Generating certificates without cmd commands
// func GenerateCerts(domain string, basePath string) error {
// 	log.Println("Cert 1")
// 	// Generate ca.key in harbor directory
// 	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
// 	if err != nil {
// 		return err
// 	}
// 	caKeyBytes := x509.MarshalPKCS1PrivateKey(caKey)
// 	caKeyFile, err := os.Create(basePath + "/ca.key")
// 	if err != nil {
// 		return err
// 	}
// 	defer caKeyFile.Close()
// 	if err := pem.Encode(caKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: caKeyBytes}); err != nil {
// 		return err
// 	}

// 	log.Println("Cert 2")
// 	// Generate ca.crt
// 	caTemplate := &x509.Certificate{
// 		SerialNumber: big.NewInt(1),
// 		Subject: pkix.Name{
// 			Country:            []string{"IN"},
// 			Organization:       []string{"Katana"},
// 			OrganizationalUnit: []string{"Katana CA"},
// 			Locality:           []string{"Delhi"},
// 			Province:           []string{"Delhi"},
// 			CommonName:         domain,
// 		},
// 		NotBefore: time.Now(),
// 		NotAfter:  time.Now().AddDate(10, 0, 0), // 10 years validity
// 		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
// 		ExtKeyUsage: []x509.ExtKeyUsage{
// 			x509.ExtKeyUsageServerAuth,
// 		},
// 		BasicConstraintsValid: true,
// 		IsCA:                  true,
// 	}
// 	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
// 	if err != nil {
// 		return err
// 	}
// 	caCertFile, err := os.Create(basePath + "/ca.crt")
// 	if err != nil {
// 		return err
// 	}
// 	defer caCertFile.Close()
// 	if err := pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
// 		return err
// 	}

// 	log.Println("Cert 3")
// 	// Generate private key
// 	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
// 	if err != nil {
// 		return err
// 	}
// 	privateKeyFile, err := os.Create(basePath + "/" + domain + ".key")
// 	if err != nil {
// 		return err
// 	}
// 	defer privateKeyFile.Close()
// 	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
// 	if err != nil {
// 		return err
// 	}
// 	if err := pem.Encode(privateKeyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
// 		return err
// 	}

// 	log.Println("Cert 4")
// 	// Generate certificate signing request
// 	csrTemplate := &x509.CertificateRequest{
// 		Subject: pkix.Name{
// 			Country:            []string{"IN"},
// 			Organization:       []string{"Katana"},
// 			OrganizationalUnit: []string{"Katana"},
// 			Locality:           []string{"Delhi"},
// 			Province:           []string{"Delhi"},
// 			CommonName:         domain,
// 		},
// 		DNSNames: []string{domain},
// 	}
// 	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privateKey)
// 	if err != nil {
// 		return err
// 	}
// 	csrFile, err := os.Create(basePath + "/" + domain + ".csr")
// 	if err != nil {
// 		return err
// 	}
// 	defer csrFile.Close()
// 	if err := pem.Encode(csrFile, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}); err != nil {
// 		return err
// 	}

// 	log.Println("Cert 5")
// 	// Generate v3.ext file
// 	extFileContent := fmt.Sprintf("authorityKeyIdentifier=keyid,issuer\nbasicConstraints=CA:FALSE\nkeyUsage=digitalSignature,keyEncipherment\nextendedKeyUsage=serverAuth\nsubjectAltName=DNS:%s", domain)
// 	extFile, err := os.Create(basePath + "/v3.ext")
// 	if err != nil {
// 		return err
// 	}
// 	defer extFile.Close()
// 	if _, err := extFile.WriteString(extFileContent); err != nil {
// 		return err
// 	}

// 	log.Println("Cert 6")
// 	// Generate certificate
// 	certTemplate := &x509.Certificate{
// 		SerialNumber: big.NewInt(2),
// 		Subject: pkix.Name{
// 			Country:            []string{"IN"},
// 			Organization:       []string{"Katana"},
// 			OrganizationalUnit: []string{"Katana"},
// 			Locality:           []string{"Delhi"},
// 			Province:           []string{"Delhi"},
// 			CommonName:         domain,
// 		},
// 		NotBefore: time.Now(),
// 		NotAfter:  time.Now().AddDate(10, 0, 0), // 10 years validity
// 		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
// 		ExtKeyUsage: []x509.ExtKeyUsage{
// 			x509.ExtKeyUsageServerAuth,
// 		},
// 		BasicConstraintsValid: true,
// 	}
// 	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, caTemplate, &privateKey.PublicKey, caKey)
// 	if err != nil {
// 		return err
// 	}
// 	certFile, err := os.Create(basePath + "/" + domain + ".crt")
// 	if err != nil {
// 		return err
// 	}
// 	defer certFile.Close()
// 	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
// 		return err
// 	}

// 	log.Println("Cert 7")
// 	return nil
// }

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
