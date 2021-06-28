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

// Get kubernetes REST config
func getKubeConfig() (*rest.Config, error) {
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

// Get kubernetes clientset
func GetKubeClient() (*kubernetes.Clientset, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func GetPodByName(clientset *kubernetes.Clientset, podName string) (*v1.Pod, error) {
	client := clientset.CoreV1()
	podsInterface := client.Pods(g.KatanaConfig.KubeNameSpace)
	return podsInterface.Get(context.Background(), podName, metav1.GetOptions{})
}

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
