package flaghandlerservice

import (
	"context"
	"fmt"

	utils "github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Server() {
	// Get the ticker
	ticker := utils.GetTicker()

	// Get Kubernetes client
	client, err := utils.GetKubeClient()
	if err != nil {
		fmt.Println(err)
	}

	// Define the pod name and namespace
	podName := "kashira"
	namespace := "katana"

	go func() {
		for range ticker.C {
			// Get the pod
			pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
			if err != nil {
				fmt.Println(err)
			}

			// Modify annotations
			pod.Annotations["tick"] = "true"

			// Update the pod
			_, err = client.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
}
