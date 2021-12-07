package cloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
	restclient "k8s.io/client-go/rest"
)

type azure struct {
}

func ApplyAzureCluster() error {

	tf, err := obtainTfexec(pathToAzureTf)
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

func DestroyAzureCluster() error {
	tf, err := obtainTfexec(pathToAzureTf)
	if err != nil {
		return err
	}

	err = tf.Destroy(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func GetAzureKubeConfig() (*restclient.Config, error) {
	tf, err := obtainTfexec(pathToAzureTf)
	if err != nil {
		return nil, err
	}
	output, err := tf.Output(context.Background())

	if err != nil {
		return nil, err
	}
	fmt.Println(output)
	for k, v := range output {
		fmt.Println(k, v.Value)
	}
	return nil, nil
}
