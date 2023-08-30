package wireguard

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"path/filepath"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetupWireguard() error {
	kubeConfig, _ := utils.GetKubeConfig()
	kubeClient, _ := utils.GetKubeClient()

	namespace := "katana"

	manifest := &bytes.Buffer{}

	tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.TemplatedManifestDir, "runtime", "wireguard.yml"))
	if err != nil {
		return err
	}

	deploymentConfig := utils.DeploymentConfig()

	serviceName := "wireguard"

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	wireguard_lbIP := service.Status.LoadBalancer.Ingress

	//print the wireguard_lbIP
	log.Println(wireguard_lbIP[0].IP)
	deploymentConfig.WireguardIP = wireguard_lbIP[0].IP

	if err := tmpl.Execute(manifest, deploymentConfig); err != nil {
		return err
	}

	if err = deployment.ApplyManifest(kubeConfig, kubeClient, manifest.Bytes(), namespace); err != nil {
		return err
	}

	return nil
}
