package infrasetservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mysql"
	utils "github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateCertsforHarbor() {
	path, _ := os.Getwd()
	path = path + "/lib/harbor/certs"

	log.Println("CHECK 1")
	// Delete the directory if it already exists
	if _, err := os.Stat(path); os.IsExist(err) {
		errDir := os.RemoveAll(path)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	log.Println("CHECK 2")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errDir := os.Mkdir(path, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	log.Println("CHECK 3")
	// Generate the certificates
	if err := utils.GenerateCerts("harbor.katana.local", path); err != nil {
		log.Fatal(err)
	}
	log.Println("CHECK 4")
}

func createTeamCredentials(teamNumber int) (string, types.CTFTeam) {
	teamlabels := utils.GetTeamPodLabels()
	podName := teamlabels + "-team-master-pod-0"
	gogs := utils.GetKatanaLoadbalancer() + ":3000"
	pwd := utils.RandomString(configs.SSHProviderConfig.PasswordLen)
	hashed, err := utils.HashPassword(pwd)
	if err != nil {
		log.Fatal(err)
	}
	podNamespace := "katana-team-" + fmt.Sprint(teamNumber)
	// start watching for container events
	go envVariables(gogs, pwd, podNamespace)
	team := types.CTFTeam{
		Index:    teamNumber,
		Name:     podNamespace,
		PodName:  podName,
		Password: hashed,
	}
	mysql.CreateGogsUser(team.Name, pwd)
	mysql.CreateAccessToken(team.Name, pwd)
	return pwd, team
}

func envVariables(gogs string, pwd string, podNamespace string) {
	kubeClientset, _ := utils.GetKubeClient()
	kubeConfig, _ := utils.GetKubeConfig()
	watch, _ := kubeClientset.CoreV1().Pods(podNamespace+"-ns").Watch(context.Background(), metav1.ListOptions{})
	for event := range watch.ResultChan() {
		p, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Fatal("unexpected type")
		}
		if p.Status.Phase != "Pending" {
			log.Println("Pod created")
			command := []string{"bash", "-c", "echo 'export GOGS=" + gogs + "' >> /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			command = []string{"bash", "-c", "echo 'export PASSWORD=" + pwd + "' >> /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			command = []string{"bash", "-c", "echo 'export USERNAME=" + podNamespace + "' >> /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			command = []string{"bash", "-c", "echo 'export BACKEND_URL=" + configs.KatanaConfig.BackendUrl + "/api/v1/admin/challengeUpdate' >> /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			command = []string{"bash", "-c", "echo 'export ADMIN=" + configs.AdminConfig.Username + "' >> /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			command = []string{"bash", "-c", "echo 'cd /opt/katana' >> /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			command = []string{"bash", "-c", "source /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			break
		}
	}
}

func buildKatanaServices() {
	katanaDir, err := utils.GetKatanaRootPath()
	if err != nil {
		log.Fatal(err)
	}
	katanaServicesDir := katanaDir + "/katana-services"

	services, err := os.ReadDir(katanaServicesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range services {

		invalidServiceNames := []string{".github", ".git", ".gitignore"}
		found := false
		for _, invalidName := range invalidServiceNames {
			if service.Name() == invalidName {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if service.IsDir() {
			log.Println("Building " + service.Name())
			imageName := strings.ToLower(service.Name())
			utils.BuildDockerImage(imageName, katanaServicesDir+"/"+service.Name())
		}
	}
}
