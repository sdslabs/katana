package cloud

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type Azure struct {
}

func (az Azure) CreateCluster() error {

	tf, err := obtainTfexec(PathToAzureTf)
	if err != nil {
		return err
	}

	if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
		return err
	}

	err = tf.Apply(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (az Azure) DestroyCluster() error {
	tf, err := obtainTfexec(PathToAzureTf)
	if err != nil {
		return err
	}

	err = tf.Destroy(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (az Azure) ObtainKubeConfig() error {
	tf, err := obtainTfexec(PathToAzureTf)
	if err != nil {
		return err
	}
	output, err := tf.Output(context.Background())

	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	str := string(output["kube_config"].Value[1 : len(output["kube_config"].Value)-1])
	str2 := strings.Replace(str, "\\n", "\n", -1)

	err = os.WriteFile(workingDir+PathToCloudPackage+PathToAzureTf+"/kubeconfig",
		[]byte(str2), 0644)
	return nil
}
