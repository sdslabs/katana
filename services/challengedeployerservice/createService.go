package challengedeployerservice

import (
	"context"
	"fmt"

	g "github.com/sdslabs/katana/configs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateService(chall_name, team_name string) error {

	team_namespace := "katana-" + team_name + "-ns"

	if err := GetClient(g.KatanaConfig.KubeConfig); err != nil {
		return err
	}

	serviceClient := kubeclient.CoreV1().Services(team_namespace)

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: team_namespace,
			Name:      chall_name,
		},

		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": chall_name,
			},
			Ports: []v1.ServicePort{
				{
					Name:     "http",
					Protocol: v1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}

	//Check if service already exists
	services, err := serviceClient.Get(context.TODO(), chall_name, metav1.GetOptions{})
	if services.Name == chall_name {
		fmt.Println("Service already exists for the challenge " + chall_name + " in namespace " + team_namespace)
		return nil
	}
	if err != nil {
		fmt.Println(" Error in getting services. ")
		return err
		// panic(err)
	}

	// Create Service
	fmt.Println("Creating service...")
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error creating service.. ")
		return err
		// panic(err)
	}

	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName()+" in namespace "+team_namespace)

	return nil
}
