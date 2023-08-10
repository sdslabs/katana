package deployment

import (
	"context"
	"log"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

var (
	pingers = []types.ResourcePinger{
		PingDeployments,
		PingStatefulSets,
	}
)

func PingStatefulSets(ctx context.Context, kubeclientset *kubernetes.Clientset, opts map[string]string) ([]*types.ResourceStatus, bool, error) {
	setsInterface := kubeclientset.AppsV1().StatefulSets(g.KatanaConfig.KubeNameSpace)
	sets, err := setsInterface.List(ctx, metav1.ListOptions{LabelSelector: labels.Set(opts).AsSelector().String()})
	if err != nil {
		return nil, false, err
	}

	pingresult := make([]*types.ResourceStatus, len(sets.Items))
	allReady := true
	for i, set := range sets.Items {
		status := &types.ResourceStatus{
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

func PingDeployments(ctx context.Context, kubeclientset *kubernetes.Clientset, opts map[string]string) ([]*types.ResourceStatus, bool, error) {
	deploymentsInterface := kubeclientset.AppsV1().Deployments(g.KatanaConfig.KubeNameSpace)
	deployments, err := deploymentsInterface.List(ctx, metav1.ListOptions{LabelSelector: labels.Set(opts).AsSelector().String()})
	if err != nil {
		return nil, false, err
	}

	pingresult := make([]*types.ResourceStatus, len(deployments.Items))
	allReady := true
	for i, deployment := range deployments.Items {
		status := &types.ResourceStatus{
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
			// var results []*types.ResourceStatus
			for _, pinger := range pingers {
				result, ready, err := pinger(ctx, kubeclientset, opts)
				if err != nil {
					done <- err
				}
				allReady = allReady && ready
				// temporary fix for lint complains
				log.Printf("result: %v", result)
				// results = append(results, result...)
			}
			if allReady {
				close(done)
			}
		}
	}()

}
