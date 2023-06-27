package challengedeployerservice

import (
	"context"
	"fmt"
	"strconv"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteChallenge(chall_name string) error {

	totalteams := utils.GetTeamNumber()

	for i := 0; i < totalteams; i++ {

		team_name := "team-" + strconv.Itoa(i)

		team_namespace := "katana-" + team_name + "-ns"
		kubeclient, err := utils.GetKubeClient()
		if err != nil {
			return err
		}

		fmt.Println("---------------Deleting challenge for team: ", team_namespace)
		serviceClient := kubeclient.CoreV1().Services(team_namespace)
		deploymentsClient := kubeclient.AppsV1().Deployments(team_namespace)

		//Get deployment
		deps, err := deploymentsClient.Get(context.TODO(), chall_name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(" Error in getting deployments associated with the challenge. ")
			continue
			//panic(err)
		}

		//Delete deployments
		if deps.Name != chall_name {
			fmt.Println("Deployment does not exist. Create one using /deploy route.")
			return nil
		} else {
			fmt.Println("Deleting deployment...")
			deletePolicy := metav1.DeletePropagationForeground
			err = deploymentsClient.Delete(context.TODO(), chall_name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
			if err != nil {
				fmt.Println("Error in deleting deployment.")
				continue
			}
		}

		//Check if service exists
		services, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(" Error in getting services for" + chall_name + "in the namespace " + team_namespace)
			continue
			//panic(err)
		}

		flag := 0
		for _, service := range services.Items {
			if service.Name == chall_name {
				flag = 1
			}
		}
		if flag == 0 {
			fmt.Println("Service does not exist for the " + chall_name + " in namespace " + team_namespace)
			continue
		}

		//Delete service
		fmt.Println("Deleting services associated with this challenge...")
		err = serviceClient.Delete(context.TODO(), chall_name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("Error in deleting service for "+chall_name+" in namespace "+team_namespace, err)
			continue
		}

	}

	fmt.Println("Process completed")
	return nil
}
