package chalDeployerService

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/services/challengedeployerservice"
	"github.com/sdslabs/katana/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ChallengeUpdate() error {
	replicas := int32(1)
	client, _ := utils.GetKubeClient()
	patch := true

	//http connection configuration for 30 min

	var p types.GogsRequest
	if err := c.BodyParser(&p); err != nil {
		return err
	}

	dir := p.Repository.FullName
	s := strings.Split(dir, "/")
	challengeName := s[1]
	teamName := s[0]
	namespace := teamName + "-ns"
	log.Println("Challenge update request received for", challengeName, "by", teamName)
	repo, err := git.PlainOpen("teams/" + dir)
	if err != nil {
		log.Println(err)
		return err
	}

	auth := &http.BasicAuth{
		Username: g.AdminConfig.Username,
		Password: g.AdminConfig.Password,
	}

	worktree, err := repo.Worktree()
	worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil {
		return fmt.Errorf("error pulling changes:%s", err)

	}
	katanaDir, err := utils.GetKatanaRootPath()
	if err != nil {
		return err
	}
	imageName := strings.Replace(dir, "/", "-", -1)

	log.Println("Pull successful for", teamName, ". Building image...")
	firstPatch := !utils.DockerImageExists(imageName)
	err = utils.BuildDockerImage(imageName, katanaDir+"/teams/"+dir)
	if err != nil {
		return err
	}

	if firstPatch {
		log.Println("First Patch for", teamName)
		err = deployment.DeployChallengeToCluster(challengeName, teamName, patch, replicas)
		if err != nil {
			return err
		}
	} else {
		log.Println("Not the first patch for", teamName, ". Simply deploying the image...")
		labelSelector := metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": challengeName,
			},
		}
		// Delete the challenge pod
		err = client.CoreV1().Pods(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(&labelSelector),
		})
		if err != nil {
			log.Println("Error")
			log.Println(err)
			return err
		}
	}
	log.Println("Image built for", teamName)
	return nil
}

func DeployChallenge() error {
	patch := false
	replicas := int32(1)
	challengeType := "web"
	log.Println("Starting")

	//Read folder challenge by os
	dir, err := os.Open("./challenges")

	//Loop over all subfolders in the challenge folder
	if err != nil {
		log.Println("Error in opening challenges folder")
		return err
	}
	defer dir.Close()

	//Read all challenges in the folder
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Println("Error in reading challenges folder")
		return err
	}

	res := make([][]string, 0)

	//Loop over all folders
	for _, fileInfo := range fileInfos {
		//Check if it is a directory
		if fileInfo.IsDir() {
			//Get the challenger name
			folderName := fileInfo.Name()
			log.Println("Folder name is : " + folderName)
			//Update challenge path to be absolute path
			challengePath, _ := os.Getwd()
			challengePath = challengePath + "/challenges/" + folderName
			log.Println("Challenge path is : " + challengePath)
			log.Println(challengePath + "/" + folderName + "/" + folderName)

			//Check if the folder has a Dockerfile
			if _, err := os.Stat(challengePath + "/" + folderName + "/" + folderName); err != nil {
				log.Println("Dockerfile not found in the " + folderName + " challenge folder. Please follow proper format.")
			} else {
				//Update challenge path to get dockerfile
				err := utils.BuildDockerImage(folderName, challengePath+"/"+folderName+"/"+folderName)
				if err != nil {
					return err
				}
				clusterConfig := g.ClusterConfig
				numberOfTeams := clusterConfig.TeamCount
				for i := 0; i < int(numberOfTeams); i++ {
					log.Println("-----------Deploying challenge for team: " + strconv.Itoa(i) + " --------")
					teamName := "katana-team-" + strconv.Itoa(i)
					err = deployment.DeployChallengeToCluster(folderName, teamName, patch, replicas)
					if err != nil {
						return err
					}
					url, err := challengedeployerservice.CreateServiceForChallenge(folderName, teamName, 3000, i)
					if err != nil {
						return err
					} else {
						res = append(res, []string{teamName, url})
					}
				}
			}
			err = challengedeployerservice.CopyChallengeIntoTsuka(challengePath, folderName, challengeType)
			if err != nil {
				return err
			}
			err = challengedeployerservice.CopyFlagDataIntoKashira(challengePath, folderName)
			if err != nil {
				return err
			}
			err = challengedeployerservice.CopyChallengeCheckerIntoKissaki(challengePath, folderName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteChallenge(challengeName string) error {

	//Delete chall directory also ?
	log.Println("Challenge name is : " + challengeName)

	dirPath, _ := os.Getwd()
	challengePath := dirPath + "/challenges/" + challengeName
	log.Println("Deleting challenge folder :", challengePath)
	err := os.RemoveAll(challengePath)
	if err != nil {
		log.Println("Error in deleting challenge folder")
		return err
	}

	totalTeams := utils.GetTeamNumber()

	for i := 0; i < totalTeams; i++ {

		teamName := "team-" + strconv.Itoa(i)

		teamNamespace := "katana-" + teamName + "-ns"
		kubeclient, err := utils.GetKubeClient()
		if err != nil {
			return err
		}

		log.Println("---------------Deleting challenge for team: ", teamNamespace)
		serviceClient := kubeclient.CoreV1().Services(teamNamespace)
		deploymentsClient := kubeclient.AppsV1().Deployments(teamNamespace)

		//Get deployment
		deps, err := deploymentsClient.Get(context.TODO(), challengeName, metav1.GetOptions{})
		if err != nil {
			log.Println(" Error in getting deployments associated with the challenge. ")
			continue
		}

		//Delete deployments
		if deps.Name != challengeName {
			log.Println("Deployment does not exist. Create one using /deploy route.")
			return nil
		} else {
			log.Println("Deleting deployment...")
			deletePolicy := metav1.DeletePropagationForeground
			err = deploymentsClient.Delete(context.TODO(), challengeName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
			if err != nil {
				log.Println("Error in deleting deployment.")
				continue
			}
		}

		//Check if service exists
		services, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Println(" Error in getting services for" + challengeName + "in the namespace " + teamNamespace)
			return err
		}

		flag := 0
		for _, service := range services.Items {
			if service.Name == challengeName {
				flag = 1
			}
		}
		if flag == 0 {
			log.Println("Service does not exist for the " + challengeName + " in namespace " + teamNamespace)
			continue
		}

		//Delete service
		log.Println("Deleting services associated with this challenge...")
		err = serviceClient.Delete(context.TODO(), challengeName, metav1.DeleteOptions{})
		if err != nil {
			log.Println("Error in deleting service for "+challengeName+" in namespace "+teamNamespace, err)
			continue
		}
	}

	log.Println("Process completed")

	return nil
}
