package cloud

import (
	restclient "k8s.io/client-go/rest"
)

// "github.com/hashicorp/terraform-exec/tfexec"
// "context"

type azure struct {
}

func (az azure) ApplyCluster() error {

	// tf, err := obtainTfexec(pathToAzureTf)
	// if err != nil {
	// 	return err
	// }

	// if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
	// 	return err
	// }

	// err = tf.Apply(context.Background())
	// if err != nil {
	// 	return err
	// }
	//output, err := tf.Output(context.Background())

	// if err != nil {
	// 	return err
	// }

	// for k, v := range output {

	// }
	return nil
}

func (az azure) DestroyCluster() error {
	//tf, err := obtainTfexec(pathToAzureTf)
	// if err != nil {
	// 	return err
	// }

	// err = tf.Destroy(context.Background())
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (az azure) GetKubeConfig() (*restclient.Config, error) {

	return nil, nil
}
