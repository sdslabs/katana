package challengedeployerservice

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteChallenge(chall_name string) error {

	//Delete Challenge Folder
	dirPath, _ := os.Getwd()
	challengePath := dirPath + "/challenges/" + chall_name
	log.Println("Deleting challenge folder :", challengePath)
	err := os.RemoveAll(challengePath)
	if err != nil {
		log.Println("Error in deleting challenge folder")
		return err
	}

	totalteams := utils.GetTeamNumber()

	for i := 0; i < totalteams; i++ {

		team_name := "team-" + strconv.Itoa(i)

		team_namespace := "katana-" + team_name + "-ns"
		kubeclient, err := utils.GetKubeClient()
		if err != nil {
			return err
		}

		log.Println("---------------Deleting challenge for team: ", team_namespace)
		serviceClient := kubeclient.CoreV1().Services(team_namespace)
		deploymentsClient := kubeclient.AppsV1().Deployments(team_namespace)

		//Get deployment
		deps, err := deploymentsClient.Get(context.TODO(), chall_name, metav1.GetOptions{})
		if err != nil {
			log.Println(" Error in getting deployments associated with the challenge. ")
			continue
			//panic(err)
		}

		//Delete deployments
		if deps.Name != chall_name {
			log.Println("Deployment does not exist. Create one using /deploy route.")
			return nil
		} else {
			log.Println("Deleting deployment...")
			deletePolicy := metav1.DeletePropagationForeground
			err = deploymentsClient.Delete(context.TODO(), chall_name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
			if err != nil {
				log.Println("Error in deleting deployment.")
				continue
			}
		}

		//Check if service exists
		services, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Println(" Error in getting services for" + chall_name + "in the namespace " + team_namespace)
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
			log.Println("Service does not exist for the " + chall_name + " in namespace " + team_namespace)
			continue
		}

		//Delete service
		log.Println("Deleting services associated with this challenge...")
		err = serviceClient.Delete(context.TODO(), chall_name, metav1.DeleteOptions{})
		if err != nil {
			log.Println("Error in deleting service for "+chall_name+" in namespace "+team_namespace, err)
			continue
		}

	}

	log.Println("Process completed")
	return nil
}
