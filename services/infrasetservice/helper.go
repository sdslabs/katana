package infrasetservice

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mysql"
	utils "github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func GenerateCertsforHarbor() error {
	path, _ := os.Getwd()
	path = path + "/lib/harbor/certs"

	// Delete the directory if it already exists
	if _, err := os.Stat(path); os.IsExist(err) {
		errDir := os.RemoveAll(path)
		if errDir != nil {
			log.Fatal(err)
			return err
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		errDir := os.Mkdir(path, 0755)
		if errDir != nil {
			log.Fatal(err)
			return err
		}
	}

	// Generate the certificates
	if err := utils.GenerateCerts("harbor.katana.local", path); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func CreateTeamCredentials(teamNumber int) (string, types.CTFTeam, error) {
	teamlabels := utils.GetTeamPodLabels()
	podName := teamlabels + "-team-master-pod-0"
	gogs, err := utils.GetKatanaLoadbalancer()
	if err != nil {
		log.Fatal(err)
		return "", types.CTFTeam{}, err
	}
	gogs = gogs + ":3000"
	pwd := utils.RandomString(configs.SSHProviderConfig.PasswordLen)
	hashed, err := utils.HashPassword(pwd)
	if err != nil {
		log.Fatal(err)
		return "", types.CTFTeam{}, err
	}
	podNamespace := "katana-team-" + fmt.Sprint(teamNumber)
	// start watching for container events

	envVariables(gogs, pwd, podNamespace)
	team := types.CTFTeam{
		Index:    teamNumber,
		Name:     podNamespace,
		PodName:  podName,
		Password: hashed,
	}
	err = mysql.CreateGogsUser(team.Name, pwd)
	if err != nil {
		log.Fatal(err)
		return "", types.CTFTeam{}, err
	}
	err = mysql.CreateAccessToken(team.Name, pwd)
	if err != nil {
		log.Fatal(err)
		return "", types.CTFTeam{}, err
	}
	return pwd, team, nil
}

func envVariables(gogs string, pwd string, podNamespace string) error {
	kubeClientset, _ := utils.GetKubeClient()
	kubeConfig, _ := utils.GetKubeConfig()
	var wg sync.WaitGroup
	wg.Add(1)
	go utils.WaitForPodReady(kubeClientset, podNamespace, &wg)
	wg.Wait()
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

	return nil
}

func BuildKatanaServices() error {
	katanaDir, errDir := utils.GetKatanaRootPath()
	if errDir != nil {
		log.Fatal(errDir)
		return errDir
	}

	katanaServicesDir := katanaDir + "/katana-services"

	services, err := os.ReadDir(katanaServicesDir)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, service := range services {
		if service.Name() == ".github" {
			continue
		}
		if service.IsDir() {
			log.Println("Building " + service.Name())
			imageName := strings.ToLower(service.Name())
			err := utils.BuildDockerImage(imageName, katanaServicesDir+"/"+service.Name())
			if err != nil {
				return err
			}

		}
	}
	return nil
}
