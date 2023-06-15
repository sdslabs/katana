package challengecheckerservice

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeclient *kubernetes.Clientset

// Get kubernetes client
func GetClient(pathToCfg string) error {
	if pathToCfg == "" {
		pathToCfg = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}
	config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	kubeclient = client
	return nil
}
