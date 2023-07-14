package challengedeployerservice

import (
	"context"
	"log"
	"os/exec"

	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateService(chall_name, team_name string) (string, error) {

	team_namespace := team_name + "-ns"
	kubeclient, err := utils.GetKubeClient()
	if err != nil {
		return "", err
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

	//Get all services

	//Check if service already exists
	services, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
	for _, service := range services.Items {
		if service.Name == chall_name {
			log.Println("Service already exists for the challenge " + chall_name + " in namespace " + team_namespace)
			return "", nil
		}
	}

	if err != nil {
		log.Println(" Error in getting services. ")
		//return err
		panic(err)
	}

	// Create Service
	log.Println("Creating service...")
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		log.Println("Error creating service.. ")
		return "", err
		// panic(err)
	}

	log.Printf("Created service %q.\n", result.GetObjectMeta().GetName()+" in namespace "+team_namespace)

	// expose service to localhost
	// TODO: change implementation when deploying on cluster
	url, err := ExposeService(chall_name, team_namespace)
	if err != nil {
		log.Printf("Error in exposing service %s for namespace %s", chall_name, team_namespace)
		log.Println("Error: ", err)
		return "", err
	}

	log.Printf("Challenge for %s is deployed at %s", team_name, url)

	return url, nil
}

func ExposeService(service_name, namespace string) (string, error) {
	// run command to expose service minikube service <service_name> -n <namespace> --url
	out := exec.Command("minikube", "service", service_name, "-n", namespace, "--url")
	url, err := out.Output()
	if err != nil {
		return "", err
	}
	return string(url), nil
}
