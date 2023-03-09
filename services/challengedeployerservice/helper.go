package challengedeployerservice

import (
	"context"
	"os"
	"path/filepath"

	g "github.com/sdslabs/katana/configs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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

func getPods(lbls map[string]string, ns ...string) ([]v1.Pod, error) {
	var namespace string
	if len(ns) == 0 {
		namespace = g.KatanaConfig.KubeNameSpace
	} else {
		namespace = ns[0]
	}

	selector := labels.SelectorFromSet(lbls)
	pods, err := kubeclient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}
