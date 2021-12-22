package cloud

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hashicorp/terraform-exec/tfexec"
	g "github.com/sdslabs/katana/configs"
)

type GCP struct {
}

func (gcp GCP) CreateCluster() error {
	err := createGCPTerraformFile()

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

func (gcp GCP) DestroyCluster() error {
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

func createGCPTerraformFile() error {
	gcpConfig := g.KatanaConfig.GCPConfig

	if gcpConfig.ProjectID == "" {
		return errors.New("Project ID for GCP is not set")
	}

	if gcpConfig.CredentialsFile == "" {
		return errors.New("Credentials file for GCP is not set")
	}

	if gcpConfig.ClusterName == "" {
		gcpConfig.ClusterName = "katana-cluster"
	}

	gcpConfig.ProjectID = "\"" + gcpConfig.ProjectID + "\""
	gcpConfig.CredentialsFile = "\"" + gcpConfig.CredentialsFile + "\""

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	pathToTemplate := filepath.Join(workingDir, "manifests", "templates", "gcp_template.tf")

	tmpl, err := template.ParseFiles(pathToTemplate)
	if err != nil {
		return err
	}
	manifest := &bytes.Buffer{}

	err = tmpl.Execute(manifest, gcpConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(workingDir, "main.tf"), manifest.Bytes(), 0644)
	return nil
}
