package cloud

import (
	"bytes"
	"context"
	"os"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-exec/tfexec"
	g "github.com/sdslabs/katana/configs"
)

type Azure struct {
}

func (az Azure) CreateCluster() error {
	err := createAzureTerraformFile()
	if err != nil {
		return err
	}

	tf, err := obtainTfexec()
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
	tf, err := obtainTfexec()
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
	tf, err := obtainTfexec()
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

	err = os.WriteFile(workingDir+"/kubeconfig",
		[]byte(str2), 0644)
	return nil
}

func createAzureTerraformFile() error {
	azureConfig := g.KatanaConfig.AzureConfig

	if azureConfig.ResourceGroupName == "" {
		azureConfig.ResourceGroupName = "katana"
	}

	if azureConfig.ClusterName == "" {
		azureConfig.ClusterName = "katanaCluster"
	}

	if azureConfig.Location == "" {
		azureConfig.Location = "centralindia"
	}

	azureConfig.ResourceGroupName = "\"" + azureConfig.ResourceGroupName + "\""
	azureConfig.ClusterName = "\"" + azureConfig.ClusterName + "\""
	azureConfig.Location = "\"" + azureConfig.Location + "\""

	pathToTemplate := "./manifests/templates/azure_template.tf"

	tmpl, err := template.ParseFiles(pathToTemplate)
	if err != nil {
		return err
	}
	manifest := &bytes.Buffer{}

	err = tmpl.Execute(manifest, azureConfig)
	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.WriteFile(workingDir+"/main.tf", manifest.Bytes(), 0644)
	return nil
}
