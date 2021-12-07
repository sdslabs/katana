package cloud

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

func InitializeTf() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	tfExecDir := workingDir + relativePathToTfexe
	err = os.Mkdir(tfExecDir, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			err = DestroyTf()
			if err != nil {
				return err
			}
			err = os.Mkdir(tfExecDir, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	execPath, err = tfinstall.Find(context.Background(), tfinstall.LatestVersion(tfExecDir, false))
	if err != nil {
		return err
	}

	return nil
}

func DestroyTf() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.RemoveAll(workingDir + relativePathToTfexe)

	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No Terraform binary found")
		} else {
			return err
		}
	}

	return nil
}

func obtainTfexec(path string) (*tfexec.Terraform, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	execPath = workingDir + relativePathToTfexe + "/terraform"
	workingDir += pathToCloudPackage + path
	return tfexec.NewTerraform(workingDir, execPath)
}
