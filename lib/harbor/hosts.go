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

	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "trusted-ca", "kube-system")
	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "setup-script", "kube-system")
	utils.DeleteDaemonSetAndWait(kubeClient, kubeConfig, "node-custom-setup", "kube-system")

	serviceName := "ingress-nginx-controller"
	serviceNamespace := "katana"
	deplNamespace := "kube-system"

	basePath, _ := os.Getwd()
	pathToCaCrt := basePath + "/lib/harbor/certs/ca.crt"

	data, err := os.ReadFile(pathToCaCrt)
	if err != nil {
		return err
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trusted-ca",
			Namespace: "kube-system",
		},
		Data: map[string]string{
			"ca.crt": string(data),
		},
	}

	if _, err := kubeClient.CoreV1().ConfigMaps("kube-system").Create(context.Background(), configMap, metav1.CreateOptions{}); err != nil {
		return err
	}

	manifest := &bytes.Buffer{}

	tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.TemplatedManifestDir, "runtime", "harbor-daemonset.yml"))
	if err != nil {
		return err
	}

	deploymentConfig := utils.DeploymentConfig()

	service, err := kubeClient.CoreV1().Services(serviceNamespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	clusterIP := service.Spec.ClusterIP

	deploymentConfig.HarborIP = clusterIP

	if err := tmpl.Execute(manifest, deploymentConfig); err != nil {
		return err
	}

	if err = deployment.ApplyManifest(kubeConfig, kubeClient, manifest.Bytes(), deplNamespace); err != nil {
		return err
	}

	return nil
}
