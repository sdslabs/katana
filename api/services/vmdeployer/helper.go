package vmdeployer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/nomad/api"
)

var path, _ = os.Getwd()
var artifactsPath = filepath.Join(path, "artifacts")

func checkAllocation(c chan status, task string, allocID string, allocs api.Allocations) {
	var res status
	for {
		alloc, _, err := allocs.Info(allocID, &api.QueryOptions{})

		if err != nil {
			res.Error = fmt.Sprintf("%s : Failed to fetch Information of the Allocation", err)
			c <- res
			return
		}

		if alloc.DesiredStatus == "stop" {
			res.Error = fmt.Sprintf("Allocation is not running")
			c <- res
			return
		} else if alloc.ClientStatus == "failed" {
			res.Error = fmt.Sprintf("Allocation failed to run")
			c <- res
			return
		} else if alloc.ClientStatus == "running" {
			break
		}

	}

	filePath := fmt.Sprintf("%s/%s-%s", getTempDirectory(), task, allocID)
	data, errh := getIPFromFile(filePath)
	res.Error = errh
	res.Data = data
	c <- res
	return
}

func getIPFromFile(fileName string) (allocation, interface{}) {

	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return allocation{}, fmt.Sprintf("%s : Error reading file", err)
	}

	data := allocation{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return allocation{}, fmt.Sprintf("%s : Failed to unmarshal file", err)
	}

	return data, nil

}

func checkAndPullImage(chImage chan bool) {
	imagePath := filepath.Join(artifactsPath, "vmlinux")
	if _, err := os.Stat(imagePath); err != nil {
		if os.IsNotExist(err) {
			if err = getImageAndBootDisk("vmlinux-5.4.0-rc5.tar.gz", "kernels"); err != nil {
				chImage <- false
			}
		} else {
			chImage <- false
		}
	}
	if err := unZipFile("vmlinux-5.4.0-rc5.tar.gz"); err != nil {
		chImage <- false
	}
	chImage <- true
}

func checkAndPullBootDisk(chDisk chan bool) {
	imagePath := filepath.Join(path, "artifacts/rootfs.ext4")
	if _, err := os.Stat(imagePath); err != nil {
		if os.IsNotExist(err) {
			if err = getImageAndBootDisk("ubuntu18.04.rootfs.tar.gz", "rootfs"); err != nil {
				chDisk <- false
			}
		} else {
			chDisk <- false
		}
	}
	if err := unZipFile("ubuntu18.04.rootfs.tar.gz"); err != nil {
		chDisk <- false
	}
	chDisk <- true
}

func getImageAndBootDisk(fileName, fileType string) error {
	cmd := exec.Command("wget", fmt.Sprintf("https://firecracker-%s.s3-sa-east-1.amazonaws.com/%s", fileType, fileName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = artifactsPath
	return cmd.Run()
}

func unZipFile(fileName string) error {
	cmd := exec.Command("tar", "xzf", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = artifactsPath
	return cmd.Run()
}

func getKernelImagePath() string {
	imagePath := filepath.Join(path, "artifacts/vmlinux")
	return imagePath
}

func getBootDiskPath() string {
	bootDiskPath := filepath.Join(path, "artifacts/rootfs.ext4")
	return bootDiskPath

}

func getTempDirectory() string {
	return os.TempDir()
}
