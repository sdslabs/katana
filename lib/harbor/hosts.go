package harbor

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func updateOrAddHost() error {
	filePath := "/etc/hosts"
	hostname := "harbor.katana.local"
	ipAddress := utils.GetKatanaLoadbalancer()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	found := false

	// Search for the hostname and update if found
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) > 1 && parts[1] == hostname {
			found = true
			lines = append(lines, ipAddress+" "+hostname)
			lines = append(lines, "")
		} else {
			lines = append(lines, line)
		}
	}

	// If hostname not found, add it
	if !found {
		lines = append(lines, ipAddress+" "+hostname)
		lines = append(lines, "")
	}

	output := strings.Join(lines, "\n")
	if err := os.WriteFile(filePath, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

func deployHarborClusterDaemonSet() error {
	kubeConfig, _ := utils.GetKubeConfig()
	kubeClient, _ := utils.GetKubeClient()

	namespace := "kube-system"

	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "trusted-ca", namespace)
	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "setup-script", namespace)
	utils.DeleteDaemonSetAndWait(kubeClient, kubeConfig, "node-custom-setup", namespace)

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

	serviceName := "harbor"
	namespace = "katana"

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
