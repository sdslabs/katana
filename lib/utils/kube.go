package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
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

func GetPods(lbls map[string]string, ns ...string) ([]corev1.Pod, error) {
	var namespace string
	if len(ns) == 0 {
		namespace = g.KatanaConfig.KubeNameSpace
	} else {
		namespace = ns[0]
	}

	selector := labels.SelectorFromSet(lbls)
	kubeclient, err := GetKubeClient()
	if err != nil {
		return nil, err
	}
	pods, err := kubeclient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func GetTeamNumber() int {
	return int(g.ClusterConfig.TeamCount)
}

func GetTeamPodLabels() string {
	return string(g.ClusterConfig.DeploymentLabel)
}

func GetMongoIP() string {
	client, err := GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}
	service, err := client.CoreV1().Services("katana").Get(context.TODO(), "mongo-nodeport-svc", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the IP address of the service
	log.Println(service.Spec.ClusterIP)
	return service.Spec.ClusterIP
}

func CopyFromPod(podName string, containerName string, pathInPod string, localFilePath string, ns ...string) error {
	config, err := GetKubeConfig()
	if err != nil {
		return err
	}

	client, err := GetKubeClient()
	if err != nil {
		return err
	}

	namespace := "katana"
	if len(ns) > 0 {
		namespace = ns[0]
	}

	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting pod: %s\n", err)
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
		err = fmt.Errorf("container not found in pod")
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
		Command:   []string{"cat", pathInPod},
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Printf("Error creating executor: %s\n", err)
		return err
	}

	localFile, err := os.Create(localFilePath)
	if err != nil {
		log.Printf("Error creating local file: %s\n", err)
		return err
	}
	defer localFile.Close()

	// Stream the file
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: localFile,
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

func GetKatanaLoadbalancer() string {
	client, err := GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}

	// Check if the service is ready
	err = WaitForLoadBalancerExternalIP(client, "katana-lb", "katana")
	if err != nil {
		log.Fatal(err)
	}

	service, err := client.CoreV1().Services("katana").Get(context.TODO(), "katana-lb", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	externalIP := service.Status.LoadBalancer.Ingress[0].IP
	return externalIP
}

func DeploymentConfig() types.ManifestConfig {
	config := types.ManifestConfig{
		KubeNameSpace:     g.KatanaConfig.KubeNameSpace,
		TeamCount:         g.ClusterConfig.TeamCount,
		TeamLabel:         g.ClusterConfig.TeamLabel,
		TeamPodName:       g.TeamVmConfig.TeamPodName,
		ContainerName:     g.TeamVmConfig.ContainerName,
		ChallengDir:       g.TeamVmConfig.ChallengeDir,
		TempDir:           g.TeamVmConfig.TempDir,
		InitFile:          g.TeamVmConfig.InitFile,
		DaemonPort:        g.TeamVmConfig.DaemonPort,
		MongoUsername:     Base64Encode(g.KatanaConfig.Mongo.Username),
		MongoPassword:     Base64Encode(g.KatanaConfig.Mongo.Password),
		MySQLPassword:     g.KatanaConfig.MySQL.Password,
		SSHPassword:       "",
		HarborKey:         "",
		HarborCrt:         "",
		HarborCaCrt:       "",
		HarborIP:          "",
		WireguardIP:       "",
		NodeAffinityValue: "",
	}

	basePath, _ := os.Getwd()

	harborKey, err := ioutil.ReadFile(basePath + "/lib/harbor/certs/harbor.katana.local.key")
	if err != nil {
		log.Fatal(err)
	}
	harborCrt, err := ioutil.ReadFile(basePath + "/lib/harbor/certs/harbor.katana.local.crt")
	if err != nil {
		log.Fatal(err)
	}
	harborCaCrt, err := ioutil.ReadFile(basePath + "/lib/harbor/certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	config.HarborKey = Base64Encode(string(harborKey))
	config.HarborCrt = Base64Encode(string(harborCrt))
	config.HarborCaCrt = Base64Encode(string(harborCaCrt))

	return config
}

func Podexecutor(command []string, kubeClientset *kubernetes.Clientset, kubeConfig *rest.Config, podNamespace string) {
	req := kubeClientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name("katana-team-master-pod-0").
		Namespace(podNamespace + "-ns").
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
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
}

func DeleteDaemonSetAndWait(kubeClientset *kubernetes.Clientset, kubeConfig *rest.Config, daemonSetName string, daemonSetNamespace string) {
	listOptions := metav1.ListOptions{
		FieldSelector:   "metadata.name=" + daemonSetName,
		Watch:           true,
		ResourceVersion: "0",
	}

	watcher, err := kubeClientset.AppsV1().DaemonSets(daemonSetNamespace).Watch(context.Background(), listOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = kubeClientset.AppsV1().DaemonSets(daemonSetNamespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return
		}
		log.Fatal(err)
	}

	for event := range watcher.ResultChan() {
		// Check if DaemonSet exists
		daemonSetName := event.Object.(*appsv1.DaemonSet).Name
		if daemonSetName == "" {
			break
		}
		if event.Type == watch.Deleted {
			break
		}
	}

	watcher.Stop()
}

func DeleteConfigMapAndWait(kubeClientset *kubernetes.Clientset, kubeConfig *rest.Config, configMapName string, configMapNamespace string) {
	// Wait for the configmap to be deleted
	listOptions := metav1.ListOptions{
		FieldSelector:   "metadata.name=" + configMapName,
		Watch:           true,
		ResourceVersion: "0",
	}

	watcher, err := kubeClientset.CoreV1().ConfigMaps(configMapNamespace).Watch(context.Background(), listOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = kubeClientset.CoreV1().ConfigMaps(configMapNamespace).Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return
		}
		log.Fatal(err)
	}

	for event := range watcher.ResultChan() {
		// Check if ConfigMap exists
		configMapName := event.Object.(*corev1.ConfigMap).Name
		if configMapName == "" {
			break
		}
		if event.Type == watch.Deleted {
			break
		}
	}

	watcher.Stop()
}

func WaitForLoadBalancerExternalIP(clientset *kubernetes.Clientset, serviceName string, namespace string) error {
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if service.Status.LoadBalancer.Ingress != nil && len(service.Status.LoadBalancer.Ingress) > 0 && service.Status.LoadBalancer.Ingress[0].IP != "" {
		return nil
	}

	watcher, err := clientset.CoreV1().Services(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + serviceName,
	})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		service, ok := event.Object.(*corev1.Service)
		if !ok {
			continue
		}

		if service.Status.LoadBalancer.Ingress != nil && len(service.Status.LoadBalancer.Ingress) > 0 && service.Status.LoadBalancer.Ingress[0].IP != "" {
			return nil
		}
	}

	return nil
}

func WaitForDeploymentReady(clientset *kubernetes.Clientset, deploymentName string, namespace string) error {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deployment.Status.ReadyReplicas > 0 {
		return nil
	}

	watcher, err := clientset.AppsV1().Deployments(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + deploymentName,
	})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		deployment, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			continue
		}

		if deployment.Status.ReadyReplicas > 0 {
			return nil
		}
	}

	return nil
}

func CreateService(clientset *kubernetes.Clientset, serviceName string, namespace string, port int32, targetPort int32, selector map[string]string) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: selector,
			Ports: []corev1.ServicePort{
				{
					Port:       port,
					TargetPort: intstr.FromInt(int(targetPort)),
				},
			},
		},
	}

	_, err := clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetNodes(clientset *kubernetes.Clientset) ([]corev1.Node, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodes.Items, nil
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

	namespace := "katana"
	if len(ns) > 0 {
		namespace = ns[0]
	}

	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		err := Tar(localFilePath, writer)
		cmdutil.CheckErr(err)
	}()

	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting pod: %s\n", err)
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
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Printf("Error creating executor: %s\n", err)
		return err
	}

	// Stream the file
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  reader,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		log.Printf("Error streaming the file: %s\n", err)
		return err
	}

	log.Println("File copied successfully at " + pathInPod)
	return nil
}
