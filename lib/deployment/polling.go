package deployment

import (
	"context"

	g "github.com/sdslabs/katana/configs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

var (
	pingers = []ResourcePinger{
		PingDeployments,
		PingStatefulSets,
	}
)

func PingStatefulSets(ctx context.Context, kubeclientset *kubernetes.Clientset, opts map[string]string) ([]*ResourceStatus, bool, error) {
	setsInterface := kubeclientset.AppsV1().StatefulSets(g.KatanaConfig.KubeNameSpace)
	sets, err := setsInterface.List(ctx, metav1.ListOptions{LabelSelector: labels.Set(opts).AsSelector().String()})
	if err != nil {
		return nil, false, err
	}

	pingresult := make([]*ResourceStatus, len(sets.Items))
	allReady := true
	for i, set := range sets.Items {
		status := &ResourceStatus{
			Name:          set.Name,
			TotalReplicas: *set.Spec.Replicas,
			ReadyReplicas: set.Status.ReadyReplicas,
			Ready:         set.Status.ReadyReplicas == *set.Spec.Replicas,
		}
		pingresult[i] = status
		allReady = allReady && status.Ready
	}

	return pingresult, allReady, nil
}

func PingDeployments(ctx context.Context, kubeclientset *kubernetes.Clientset, opts map[string]string) ([]*ResourceStatus, bool, error) {
	deploymentsInterface := kubeclientset.AppsV1().Deployments(g.KatanaConfig.KubeNameSpace)
	deployments, err := deploymentsInterface.List(ctx, metav1.ListOptions{LabelSelector: labels.Set(opts).AsSelector().String()})
	if err != nil {
		return nil, false, err
	}

	pingresult := make([]*ResourceStatus, len(deployments.Items))
	allReady := true
	for i, deployment := range deployments.Items {
		status := &ResourceStatus{
			Name:          deployment.Name,
			TotalReplicas: *deployment.Spec.Replicas,
			ReadyReplicas: deployment.Status.ReadyReplicas,
			Ready:         deployment.Status.ReadyReplicas == *deployment.Spec.Replicas,
		}
		pingresult[i] = status
		allReady = allReady && status.Ready
	}

	return pingresult, allReady, nil
}

func PollDeployments(kubeclientset *kubernetes.Clientset, done chan<- error) {
	ctx := context.Background()
	var opts map[string]string
	go func() {
		for {
			allReady := true
			var results []*ResourceStatus
			for _, pinger := range pingers {
				result, ready, err := pinger(ctx, kubeclientset, opts)
				if err != nil {
					done <- err
				}
				allReady = allReady && ready
				results = append(results, result...)
			}
			if allReady {
				close(done)
			}
		}
	}()
}
