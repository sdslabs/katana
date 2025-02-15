package infrasetservice

import (
	"context"
	"errors"
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

	// Delete the directory if it already exists
	_, err := os.Stat(path)
	if err == nil {
		//If it exists, delete it
		errDir := os.RemoveAll(path)
		if errDir != nil {
			logger.Fatal().Msgf("Failed to remove directory: %v", errDir)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		// If there is an error other than "does not exist", log it and exit
		logger.Fatal().Msgf("Failed to access directory: %v", err)
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		//creating directory
		errDir := os.Mkdir(path, 0755)
		if errDir != nil {
			logger.Fatal().Msgf("Failed to create directory: %v", errDir)
		}
	}

	// Generate the certificates
	if err := utils.GenerateCerts("harbor.katana.local", path); err != nil {
		logger.Fatal().Err(err)
	}
}

func createTeamCredentials(teamNumber int) (string, types.CTFTeam) {
	teamlabels := utils.GetTeamPodLabels()
	podName := teamlabels + "-team-master-pod-0"
	gogs := utils.GetKatanaLoadbalancer() + ":3000"
	pwd := utils.RandomString(configs.SSHProviderConfig.PasswordLen)
	hashed := utils.SHA256(pwd)
	podNamespace := "katana-team-" + fmt.Sprint(teamNumber)
	// start watching for container events
	go envVariables(gogs, pwd, podNamespace)
	team := types.CTFTeam{
		Index:      teamNumber,
		Name:       podNamespace,
		PodName:    podName,
		Password:   hashed,
		Score:      0,
		Challenges: []types.Challenge{},
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
			logger.Fatal().Msgf("unexpected type")
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
		logger.Fatal().Err(err)
	}
	katanaServicesDir := katanaDir + "/katana-services"

	services, err := os.ReadDir(katanaServicesDir)
	if err != nil {
		logger.Fatal().Err(err)
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
