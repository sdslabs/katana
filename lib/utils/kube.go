package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
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

func GetTeamPodLabels() string {
	return string(g.ClusterConfig.DeploymentLabel)
}

func GetTeamNumber() int {
	return int(g.ClusterConfig.TeamCount)
}

func GetMongoIP() string {
	client, err := GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}
	service, err := client.CoreV1().Services("default").Get(context.TODO(), "mongo-nodeport-svc", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the IP address of the service
	fmt.Println(service.Spec.ClusterIP)
	return service.Spec.ClusterIP
}

func CopyIntoPod(podName string, containerName string, pathInPod string, localFilePath string, ns ...string) error {
	config, err := GetKubeConfig()
	if err != nil {
		return err
	}

	client, err := GetKubeClient()
	if err != nil {
		return err
	}

	localFile, err := os.Open(localFilePath)
	if err != nil {
		log.Printf("Error opening local file: %s\n", err)
		return err
	}

	namespace := "default"
	if len(ns) > 0 {
		namespace = ns[0]
	}

	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting pod: %s\n", err)
		return err
	}

	// Find the container in the pod
	var container *corev1.Container
	for _, c := range pod.Spec.Containers {
		if c.Name == containerName {
			container = &c
			break
		}
	}

	if container == nil {
		log.Printf("Container not found in pod\n")
		return err
	}

	// Create a stream to the container
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName)

	req.VersionedParams(&corev1.PodExecOptions{
		Container: containerName,
		Command:   []string{"bash", "-c", "cat > " + pathInPod},
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Printf("Error creating executor: %s\n", err)
		return err
	}

	// Stream the file
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  localFile,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		log.Printf("Error streaming the file: %s\n", err)
		return err
	}

	log.Println("File copied successfully")
	return nil
}

func GetGogsIp() string {
	client, err := GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}
	service, err := client.CoreV1().Services("gogs").Get(context.TODO(), "gogs-svc", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return service.Spec.ClusterIP
}

func DeploymentConfig() types.ManifestConfig {
	clusterConfig := g.ClusterConfig
	config := types.ManifestConfig{
		FluentHost:            fmt.Sprintf("\"elasticsearch.%s.svc.cluster.local\"", g.KatanaConfig.KubeNameSpace),
		KubeNameSpace:         g.KatanaConfig.KubeNameSpace,
		TeamCount:             clusterConfig.TeamCount,
		TeamLabel:             clusterConfig.TeamLabel,
		BroadcastCount:        clusterConfig.BroadcastCount,
		BroadcastLabel:        clusterConfig.BroadcastLabel,
		BroadcastPort:         g.ServicesConfig.ChallengeDeployer.BroadcastPort,
		TeamPodName:           g.TeamVmConfig.TeamPodName,
		ContainerName:         g.TeamVmConfig.ContainerName,
		ChallengDir:           g.TeamVmConfig.ChallengeDir,
		TempDir:               g.TeamVmConfig.TempDir,
		InitFile:              g.TeamVmConfig.InitFile,
		DaemonPort:            g.TeamVmConfig.DaemonPort,
		ChallengeDeployerHost: g.ServicesConfig.ChallengeDeployer.Host,
		ChallengeArtifact:     g.ServicesConfig.ChallengeDeployer.ArtifactLabel,
	}
	return config
}

func Podexecutor(command []string, kubeClientset *kubernetes.Clientset, kubeConfig *rest.Config, podNamespace string) {
	req := kubeClientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name("katana-team-master-pod-0").
		Namespace(podNamespace + "-ns").
		SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		log.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(stdout.String())
	log.Println(stderr.String())
}
