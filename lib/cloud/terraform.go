package cloud

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

// "github.com/hashicorp/terraform-exec/tfexec"

func InitializeTf() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	tfExecDir := workingDir + relativePathToTfexe
	err = os.Mkdir(tfExecDir, 0755)
	if err != nil {
		return err
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

// func obtainTfexec(path string) (*tfexec.Terraform, error) {
// 	workingDir, err := os.Getwd()
// 	if err != nil {
// 		return nil, err
// 	}

// 	workingDir += pathToCloudPackage + path
// 	return tfexec.NewTerraform(workingDir, execPath)
// }
