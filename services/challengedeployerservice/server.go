package challengedeployerservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployToAll(localFilePath string, pathInPod string) error {

	if err := GetClient(g.KatanaConfig.KubeConfig); err != nil {
		return err
	}

	//regex to find challenge name since localFilePath[12:22] is hardcoded
	regexPattern := `\/([^\/]+)\.tar\.gz$`
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(localFilePath)
	filename := matches[1]

	// Get pods from different namespaces
	var pods []v1.Pod
	numberOfTeams := utils.GetTeamNumber()
	for i := 0; i < numberOfTeams; i++ {
		path := "katana-team-" + fmt.Sprint(i) + "/" + filename
		err := os.Mkdir("teams/"+path, 0755)
		if err != nil {
			log.Println(err)
		}
		git.PlainInit("teams/"+path, false)
		repo, err := git.PlainOpen("teams/" + path)
		if err != nil {
			log.Println(err)
		}
		remoteConfig := &config.RemoteConfig{
			Name: "origin",
			URLs: []string{"http://sdslabs@" + utils.GetGogsIp() + ":18080" + "/" + path}}
		_, err = repo.CreateRemote(remoteConfig)

		if err != nil {
			log.Println(err)
		}
		err = exec.Command("touch teams/" + path + "/challenge.yaml").Run()
		if err != nil {
			log.Println(err)
		}
		podsInTeam, err := getPods(map[string]string{
			"app": g.ClusterConfig.TeamLabel,
		}, "katana-team-"+fmt.Sprint(i)+"-ns")
		if err != nil {
			log.Println(err)
			return err
		}
		pods = append(pods, podsInTeam...)
	}
	// Loop over pods
	for _, pod := range pods {
		// Copy file into pod
		if err := utils.CopyIntoPod(pod.Name, g.TeamVmConfig.ContainerName, pathInPod, localFilePath, pod.Namespace); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func DeployChallenge() error {

	fmt.Println("TEST 1")
	if err := GetClient(g.KatanaConfig.KubeConfig); err != nil {
		return err
	}
	fmt.Println("TEST 2")
	//Change namespace
	deploymentsClient := kubeclient.AppsV1().Deployments(v1.NamespaceDefault)

	fmt.Println("TEST 3")
	// // Read the deployment YAML file
	// pwd, _ := os.Getwd()
	// fmt.Println(pwd)
	// deploymentYAML, err := ioutil.ReadFile("web_challenge.yaml")
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("TEST 3.5")
	// // Open the deployment YAML file
	// file, err := os.Open("web_challenge.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	fmt.Println("TEST 4")

	// Create Deployment struct dynamically when we switch to various challenges other than web.
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "chall-notekeeper",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "chall-notekeeper",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "chall-notekeeper",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "chall-notekeeper-team-0",
							Image:           "notekeeper:latest",
							ImagePullPolicy: v1.PullPolicy("Never"),
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
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

	// decoder := yaml.NewYAMLOrJSONDecoder(file, 4096)
	// err = decoder.Decode(deployment)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("TEST 5")
	// err = yaml.Unmarshal(deploymentYAML, deployment)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("TEST 6")
	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("TEST 7")
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
