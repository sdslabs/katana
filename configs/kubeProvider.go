package configs

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
func LoadKubeElements() error {
	var pathToCfg string
	if KatanaConfig.KubeConfig == "" {
		pathToCfg = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	} else {
		pathToCfg = KatanaConfig.KubeConfig
	}
	config,err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		return err
	}
	GlobalKubeConfig = config
	kubeclient,err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	GlobalKubeClient = kubeclient
	return nil
}

var 
(
	GlobalKubeConfig *rest.Config	
 	GlobalKubeClient *kubernetes.Clientset
)

