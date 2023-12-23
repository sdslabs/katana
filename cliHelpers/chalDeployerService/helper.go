package chalDeployerService

import (
	"context"
	"log"
	"os"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sdslabs/katana/configs"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/services/challengedeployerservice"
)

func DeployChallenge() error {
	patch := false
	replicas := int32(1)
	challengeType := "web"
	log.Println("Starting")

	katanaDir, err := utils.GetKatanaRootPath()
	if err != nil {
		return err
	}
	challengesDir:=katanaDir+"/challenges"	

	//Read folder challenge by os
	dir, err := os.Open(challengesDir)

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
		kubeclient := configs.GlobalKubeClient
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
			if service.Name == challengeName+"-svc-"+strconv.Itoa(i) {
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
