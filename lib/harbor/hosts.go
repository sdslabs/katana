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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

const (
	hostsFilePath = "/etc/hosts"
)

var hostsEntry string = configs.KatanaConfig.Harbor.Hostname

func checkHarborHostsEntryExists() bool {
	file, err := os.Open(hostsFilePath)
	if err != nil {
		return false
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), hostsEntry) {
			return true
		}
	}

	return false
}

func addHarborHostsEntry() error {
	client, err := utils.GetKubeClient()
	if err != nil {
		return err
	}

	serviceName := "harbor"
	namespace := "katana"

	service, err := client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	externalIP := service.Status.LoadBalancer.Ingress[0].IP

	hostsEntry = fmt.Sprintf("%s %s", externalIP, hostsEntry)

	file, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := fmt.Fprintf(file, "\n%s\n", hostsEntry); err != nil {
		return err
	}

	return nil
}

func setHostsInCluster() error {
	// Check if Harbor service is running
	client, err := utils.GetKubeClient()
	if err != nil {
		return err
	}

	serviceName := "harbor"
	namespace := "katana"

	serviceInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Services(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Services(namespace).Watch(context.TODO(), options)
			},
		},
		&v1.Service{},
		0,
		cache.Indexers{},
	)

	stop := make(chan struct{})
	go serviceInformer.Run(stop)

	serviceKey := fmt.Sprintf("%s/%s", namespace, serviceName)
	if !cache.WaitForCacheSync(stop, serviceInformer.HasSynced) {
		return fmt.Errorf("timed out waiting for caches to sync")
	}

	// Check if the service is already running
	if _, exists, _ := serviceInformer.GetIndexer().GetByKey(serviceKey); exists {
		deployHarborClusterDaemonSet()
	}

	serviceEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			service := obj.(*v1.Service)
			if service.Name == serviceName {
				deployHarborClusterDaemonSet()
				stop <- struct{}{}
			}
		},
	}

	serviceInformer.AddEventHandler(serviceEventHandler)

	<-stop

	return nil
}

func deployHarborClusterDaemonSet() error {
	kubeConfig, _ := utils.GetKubeConfig()
	kubeClient, _ := utils.GetKubeClient()
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

	tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.ManifestDir, "harbor-daemonset.yml"))
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
