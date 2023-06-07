package challengedeployerservice

import (
	"context"
	"fmt"

	g "github.com/sdslabs/katana/configs"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployChallenge(chall_name, team_name string) error {

	team_namespace := "katana-" + team_name + "-ns"

	if err := GetClient(g.KatanaConfig.KubeConfig); err != nil {
		return err
	}

	//Change namespace
	deploymentsClient := kubeclient.AppsV1().Deployments(team_namespace)

	// Creates Deployment object and deploys it
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: team_namespace,
			Name:      chall_name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": chall_name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": chall_name,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            chall_name + "-" + team_name,
							Image:           chall_name + ":latest",
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

	//Check if deployment already exists
	deps, err := deploymentsClient.Get(context.TODO(), chall_name, metav1.GetOptions{})
	if deps.Name == chall_name {
		fmt.Println("Deployment already exists for the challenge " + chall_name + " in namespace " + team_namespace)
		return nil
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		fmt.Println(" FAT GYA..SADGE :( ")
		panic(err)
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName()+" in namespace "+team_namespace)
	return nil

}

func int32Ptr(i int32) *int32 { return &i }
