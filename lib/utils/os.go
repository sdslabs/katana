package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Tar(src string, writers ...io.Writer) error {

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(file, src+string(filepath.Separator))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})
}

func RunCommand(cmd string) error {
	out := exec.Command("bash", "-c", cmd)
	err := out.Run()
	if err != nil {
		return err
	}

	return nil
}

func CheckOpenSSLVersion() (bool, error) {
	cmd := "openssl version"
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return false, err
	}
	opensslVersion:= string(output)

	versionParts := strings.Fields(opensslVersion)
	if len(versionParts) < 2 {
		return false,fmt.Errorf("unable to determine OpenSSL version")
	}

	majorVersionStr := strings.Split(versionParts[1], ".")[0]
	if majorVersionStr == "" {
		return false,fmt.Errorf("unable to determine OpenSSL major version")
	}

	majorVersion := 0
	_, err = fmt.Sscanf(majorVersionStr, "%d", &majorVersion)
	if err != nil {
		return false,fmt.Errorf("error parsing OpenSSL major version: %v", err)
	}

	if majorVersion >= 3 {
		return true,nil
	}else{
		fmt.Println("OpenSSL version 3 or higher is required")
		return false,nil
	}
}

func GetKatanaRootPath() (string, error) {
	katanaDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return katanaDir, nil
}