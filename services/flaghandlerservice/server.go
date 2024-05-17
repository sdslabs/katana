package flaghandlerservice

import (
	"context"
	"fmt"

	utils "github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespace string = "katana"
var podName string = "kashira"
var containerName string = "kashira"

func SendFlagCheckerAndUpdaterToKashira(localFilePath string) {
	pathInPod := "/opt/kashira/tmp"
	utils.CopyIntoPod(podName, containerName, pathInPod, localFilePath, namespace)
}

func Server() {
	// Get the ticker
	ticker := utils.GetTicker()

	// Get Kubernetes client
	client, err := utils.GetKubeClient()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for range ticker.C {
			pod, err := client.AppsV1().StatefulSets(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
			if err != nil {
				fmt.Println(err)
				return
			}
			pod.Annotations["tick"] = "true"
			_, err = client.AppsV1().StatefulSets(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()
}
