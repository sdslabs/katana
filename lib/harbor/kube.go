package harbor

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func waitForLoadBalancerExternalIP(clientset *kubernetes.Clientset, serviceName string) error {
	watcher, err := clientset.CoreV1().Services("katana").Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + serviceName,
	})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		service, ok := event.Object.(*v1.Service)
		if !ok {
			continue
		}

		if service.Status.LoadBalancer.Ingress != nil && len(service.Status.LoadBalancer.Ingress) > 0 && service.Status.LoadBalancer.Ingress[0].IP != "" {
			return nil
		}
	}

	return nil
}

func waitForDeploymentReady(clientset *kubernetes.Clientset, deploymentName string) error {
	watcher, err := clientset.AppsV1().Deployments("katana").Watch(context.TODO(), metav1.ListOptions{
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
