package utils

import (
	"context"
	"os"
	"path/filepath"

	g "github.com/sdslabs/katana/configs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	appLabelKey        = "app"
	deploymentLabelKey = "deployment"
)

// GetKubeConfig returns a kubernetes REST config object
func GetKubeConfig() (*rest.Config, error) {
	var pathToCfg string
	if g.KatanaConfig.KubeConfig == "" {
		pathToCfg = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	} else {
		pathToCfg = g.KatanaConfig.KubeConfig
	}

	return clientcmd.BuildConfigFromFlags("", pathToCfg)
}

// GetKubeClient returns a kubernetes clientset
func GetKubeClient() (*kubernetes.Clientset, error) {
	config, err := GetKubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// Returns pods in the cluster by name
func GetPodByName(clientset *kubernetes.Clientset, podName string) (*v1.Pod, error) {
	client := clientset.CoreV1()
	podsInterface := client.Pods(g.KatanaConfig.KubeNameSpace)
	return podsInterface.Get(context.Background(), podName, metav1.GetOptions{})
}

// Returns pods in the cluster by labels
func GetPods(clientset *kubernetes.Clientset, lbls map[string]string) ([]v1.Pod, error) {
	client := clientset.CoreV1()
	podsInterface := client.Pods(g.KatanaConfig.KubeNameSpace)
	filter := metav1.ListOptions{
		LabelSelector: labels.Set(lbls).AsSelector().String(),
	}
	pods, err := podsInterface.List(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

// Returns labels of all team pods in the cluster
func GetTeamPodLabels() map[string]string {
	return map[string]string{
		appLabelKey: g.ClusterConfig.TeamLabel,
	}
}
