package wireguard

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetupWireguard() error {
	kubeConfig, _ := utils.GetKubeConfig()
	kubeClient, _ := utils.GetKubeClient()

	namespace := "katana"

	manifest := &bytes.Buffer{}

	tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.TemplatedManifestDir, "runtime", "wireguard.yml"))
	if err != nil {
		return err
	}

	deploymentConfig := utils.DeploymentConfig()

	serviceName := "wireguard"

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	wireguard_lbIP := service.Status.LoadBalancer.Ingress

	deploymentConfig.WireguardIP = wireguard_lbIP[0].IP

	if err := tmpl.Execute(manifest, deploymentConfig); err != nil {
		return err
	}

	if err = deployment.ApplyManifest(kubeConfig, kubeClient, manifest.Bytes(), namespace); err != nil {
		return err
	}

	noOfTeams := int(configs.ClusterConfig.TeamCount)

	configPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	configPath = configPath + "/peer_configs"
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
		} else {
			fmt.Println("peer_configs folder created successfully.")
		}
	} else if err != nil {
		fmt.Printf("Error checking folder existence: %v\n", err)
	}

	for i := 0; i < noOfTeams; i++ {
		if err := GetConfigFiles(strconv.Itoa(i + 1)); err != nil {
			return err
		}
	}

	return nil
}

func GetConfigFiles(team_number string) error {

	client, _ := utils.GetKubeClient()

	deploymentNames := []string{
		"wireguard-deployment",
	}
	namespace := "katana"

	for _, deploymentName := range deploymentNames {
		if err := utils.WaitForDeploymentReady(client, deploymentName, namespace); err != nil {
			log.Printf("Error testing deployment '%s': %v\n", deploymentName, err)
		}
	}

	//get pod in the wireguard deployment
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=wireguard",
	})
	if err != nil {
		log.Printf("Error getting pod: %s\n", err)
	}

	wireguardPod := pods.Items[0]
	wireguardContainer := wireguardPod.Spec.Containers[0]

	pathInPod := "/config/peer" + team_number + "/peer" + team_number + ".conf"
	localFilePath := "./peer_configs/peer" + team_number + ".conf"

	//wait for container ready
	time.Sleep(1 * time.Minute)
	// [TODO] : Replace time.Sleep with waitForContainerRunning
	// if err := waitForContainerRunning(client, wireguardPod.Name, wireguardContainer.Name, namespace, 10*time.Minute); err != nil {
	// 	log.Printf("Error waiting for container to become running: %v\n", err)
	// }

	if err := utils.CopyFromPod(wireguardPod.Name, wireguardContainer.Name, pathInPod, localFilePath, namespace); err != nil {
		log.Println(err)
		return err
	}

	return nil

}
