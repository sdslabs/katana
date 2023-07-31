package harbor

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deployHarborClusterDaemonSet() error {
	kubeConfig, _ := utils.GetKubeConfig()
	kubeClient, _ := utils.GetKubeClient()

	namespace := "kube-system"

	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "trusted-ca", namespace)
	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "setup-script", namespace)
	utils.DeleteDaemonSetAndWait(kubeClient, kubeConfig, "node-custom-setup", namespace)

	serviceName := "ingress-nginx-controller"

	basePath, _ := os.Getwd()
	pathToCaCrt := basePath + "/lib/harbor/certs/ca.crt"

	data, err := os.ReadFile(pathToCaCrt)
	if err != nil {
		return err
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trusted-ca",
			Namespace: namespace,
		},
		Data: map[string]string{
			"ca.crt": string(data),
		},
	}

	if _, err := kubeClient.CoreV1().ConfigMaps(namespace).Create(context.Background(), configMap, metav1.CreateOptions{}); err != nil {
		return err
	}

	manifest := &bytes.Buffer{}

	tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.TemplatedManifestDir, "runtime", "harbor-daemonset.yml"))
	if err != nil {
		return err
	}

	deploymentConfig := utils.DeploymentConfig()

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	clusterIP := service.Spec.ClusterIP

	deploymentConfig.HarborIP = clusterIP

	if err := tmpl.Execute(manifest, deploymentConfig); err != nil {
		return err
	}

	if err = deployment.ApplyManifest(kubeConfig, kubeClient, manifest.Bytes(), namespace); err != nil {
		return err
	}

	return nil
}
