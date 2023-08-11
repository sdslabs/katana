package sshproviderservice

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/mysql"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func CreateTeams(teamnumber int) error {
	teamlabels := utils.GetTeamPodLabels()
	var teams []interface{}
	credsFile, err := os.Create(configs.SSHProviderConfig.CredsFile)
	if err != nil {
		return err
	}
	podName := teamlabels + "-team-master-pod-0"
	gogs := utils.GetKatanaLoadbalancer() + ":3000"
	for i := 0; i < teamnumber; i++ {
		pwd := utils.RandomString(configs.SSHProviderConfig.PasswordLen)
		hashed, err := utils.HashPassword(pwd)
		if err != nil {
			return err
		}
		podNamespace := "katana-team-" + fmt.Sprint(i)
		// start watching for container events
		go envVariables(gogs, pwd, podNamespace)
		team := types.CTFTeam{
			Index:    i,
			Name:     podNamespace,
			PodName:  podName,
			Password: hashed,
		}
		mysql.CreateGogsUser(team.Name, pwd)
		mysql.CreateAccessToken(team.Name, pwd)
		fmt.Fprintf(credsFile, "Team: %d, Username: %s, Password: %s\n", i, team.Name, pwd)
		teams = append(teams, team)
	}
	_, err = mongo.CreateTeams(teams)
	return err
}

func envVariables(gogs string, pwd string, podNamespace string) {
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
			command = []string{"bash", "-c", "source /etc/profile"}
			utils.Podexecutor(command, kubeClientset, kubeConfig, podNamespace)
			break
		}
	}
}
