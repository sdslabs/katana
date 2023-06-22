package utils

import (
	"context"
	"fmt"
	"log"

	g "github.com/sdslabs/katana/configs"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployChallenge(ChallengeName, TeamName string, FirstPatch bool, replicas int32) error {

	TeamNamespace := TeamName + "-ns"
	kubeclient, err := GetClient(g.KatanaConfig.KubeConfig)
	if err != nil {
		return err
	}

	deploymentsClient := kubeclient.AppsV1().Deployments(TeamNamespace)
	imageName := ChallengeName + ":latest"
	if FirstPatch {
		/// Retrieve the existing deployment
		existingDeployment, err := deploymentsClient.Get(context.TODO(), ChallengeName, metav1.GetOptions{})
		if err != nil {
			fmt.Println("Error in retrieving existing deployment.")
			log.Println(err)
			return err
		}

		existingDeployment.Spec.Template.Spec.Containers[0].Image = TeamName + "/" + ChallengeName

		_, err = deploymentsClient.Update(context.TODO(), existingDeployment, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println("Error in updating deployment.")
			log.Println(err)
			return err
		}

		fmt.Println("Updated deployment with new image.")
		return nil
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: TeamNamespace,
			Name:      ChallengeName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": ChallengeName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": ChallengeName,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            ChallengeName + "-" + TeamName,
							Image:           imageName,
							ImagePullPolicy: v1.PullPolicy("Never"),
							Ports: []v1.ContainerPort{
								{
									Name:          "challenge-port",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		fmt.Println("Unable to create deployement")
		panic(err)
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName()+" in namespace "+TeamNamespace)
	return nil

}
