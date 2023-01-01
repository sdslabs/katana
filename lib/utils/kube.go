package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	g "github.com/sdslabs/katana/configs"
	corev1 "k8s.io/api/core/v1"
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

func GetPodByName(clientset *kubernetes.Clientset, podName string) (*corev1.Pod, error) {
	client := clientset.CoreV1()
	podsInterface := client.Pods(g.KatanaConfig.KubeNameSpace)
	return podsInterface.Get(context.Background(), podName, metav1.GetOptions{})
}

func GetPods(clientset *kubernetes.Clientset, lbls map[string]string) ([]corev1.Pod, error) {
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

func GetTeamPodLabels() map[string]string {
	return map[string]string{
		appLabelKey: g.ClusterConfig.TeamLabel,
	}
}

func GetClusterIP() string {
	clientset, err := GetKubeClient()
	if err != nil {
		log.Println(err)
		// handle error
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	var minikubeNode *corev1.Node
	for _, node := range nodes.Items {
		if node.Spec.ProviderID == "minikube" {
			minikubeNode = &node
			break
		}
	}

	var minikubeIP string
	fmt.Print()
	fmt.Println(minikubeNode)
	for _, address := range minikubeNode.Status.Addresses {
		if address.Type == corev1.NodeExternalIP {
			fmt.Println(address.Address)
			minikubeIP = address.Address
			break
		}
	}
	log.Println(minikubeIP)
	return minikubeIP
}
