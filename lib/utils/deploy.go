package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	g "github.com/sdslabs/katana/configs"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func DeployChallenge(chall_name, team_name string, patched bool) error {

	team_namespace := team_name + "-ns"

	log.Println(team_namespace)
	kubeclient,err := GetClient(g.KatanaConfig.KubeConfig)
	if err != nil {
		return err
	}

	deploymentsClient := kubeclient.AppsV1().Deployments(team_namespace)
	imageName := chall_name + ":latest"
	if(patched){
		//Delete deployment
		fmt.Println("Deleting initial deployment...")
		deletePolicy := metav1.DeletePropagationForeground
		err = deploymentsClient.Delete(context.TODO(), chall_name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
		if(err != nil){
			fmt.Println("Error in deleting deployment.")
			log.Println(err)
		}
	
		imageName = team_name + "/" + chall_name
	}
	
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
							Image:            imageName,
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
		fmt.Println(" FAT GYA..SADGE :( ")
		panic(err)
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName()+" in namespace "+team_namespace)
	return nil

}

func int32Ptr(i int32) *int32 { return &i }

func GetClient(pathToCfg string) (*kubernetes.Clientset,error) {
	if pathToCfg == "" {
		pathToCfg = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}
	config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		return nil,err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil,err
	}
	return client,nil
}

func GetPods(lbls map[string]string, ns ...string) ([]v1.Pod, error) {
	var namespace string
	if len(ns) == 0 {
		namespace = g.KatanaConfig.KubeNameSpace
	} else {
		namespace = ns[0]
	}

	selector := labels.SelectorFromSet(lbls)
	kubeclient,err := GetClient(g.KatanaConfig.KubeConfig)
	if err != nil {
		return nil,err
	}
	pods, err := kubeclient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}
