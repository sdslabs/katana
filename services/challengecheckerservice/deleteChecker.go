package challengecheckerservice

import (
	"context"
	"fmt"
	"log"

	g "github.com/sdslabs/katana/configs"
	utils "github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteChecker(checkerName string) error {
	namespace := "katana"
	fmt.Println("Creating kube client...")
	kubeclient, err := utils.GetKubeClient(g.KatanaConfig.KubeConfig)
	if err != nil {
		log.Fatal("failed to create kube client", err)
	}
	fmt.Println("Kube client created...")
	fmt.Println("---------------Deleting checker for challenge: ", checkerName)
	cronJobClient := kubeclient.BatchV1().CronJobs(namespace)
	propagationPolicy := metav1.DeletePropagationBackground
	err = cronJobClient.Delete(context.Background(), checkerName, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy})

	if err != nil {

		log.Fatal("failed to delete cron job", err)

	}
	return nil
}
