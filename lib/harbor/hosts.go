package harbor

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
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

const (
	hostsFilePath = "/etc/hosts"
)

var hostsEntry string = configs.KatanaConfig.Harbor.Hostname

func addHarborHostsEntry() error {
	client, err := utils.GetKubeClient()
	if err != nil {
		return err
	}

	serviceName := "harbor"
	deploymentName := "katana-release-harbor-core"
	namespace := "katana"

	err = waitForLoadBalancerExternalIP(client, serviceName)
	if err != nil {
		return err
	}

	err = waitForDeploymentReady(client, deploymentName)
	if err != nil {
		return err
	}

	service, err := client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	externalIP := service.Status.LoadBalancer.Ingress[0].IP

	// Check if hosts entry already exists
	file, err := os.Open(hostsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, hostsEntry) {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				fields[0] = externalIP
				line = strings.Join(fields, " ")
				found = true
			}
		}
		lines = append(lines, line)
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s %s", externalIP, hostsEntry))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write to file
	file, err = os.OpenFile(hostsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}

func deployHarborClusterDaemonSet() error {
	kubeConfig, _ := utils.GetKubeConfig()
	kubeClient, _ := utils.GetKubeClient()

	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "trusted-ca", "kube-system")
	utils.DeleteConfigMapAndWait(kubeClient, kubeConfig, "setup-script", "kube-system")
	utils.DeleteDaemonSetAndWait(kubeClient, kubeConfig, "node-custom-setup", "kube-system")

	serviceName := "harbor"
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

	tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.ManifestRuntimeDir, "harbor-daemonset.yml"))
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
