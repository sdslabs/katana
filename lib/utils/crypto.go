package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
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
