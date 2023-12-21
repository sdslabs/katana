package challengedeployerservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gofiber/fiber/v2"
	archiver "github.com/mholt/archiver/v3"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Deploy(c *fiber.Ctx) error {
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
				utils.BuildDockerImage(folderName, challengePath+"/"+folderName+"/"+folderName)

				clusterConfig := g.ClusterConfig
				numberOfTeams := clusterConfig.TeamCount
				for i := 0; i < int(numberOfTeams); i++ {
					log.Println("-----------Deploying challenge for team: " + strconv.Itoa(i) + " --------")
					teamName := "katana-team-" + strconv.Itoa(i)
					deployment.DeployChallengeToCluster(folderName, teamName, patch, replicas)
					url, err := createServiceForChallenge(folderName, teamName, 3000, i)
					if err != nil {
						res = append(res, []string{teamName, err.Error()})
					} else {
						res = append(res, []string{teamName, url})
					}
				}
			}
			copyChallengeIntoTsuka(challengePath, folderName, challengeType)
			copyFlagDataIntoKashira(challengePath, folderName)
			copyChallengeCheckerIntoKissaki(challengePath, folderName)
		}
	}
	return c.JSON(res)
}

func DeployChallenge(c *fiber.Ctx) error {
	challengeType := "web"
	folderName := ""
	patch := false
	replicas := int32(1)
	log.Println("Starting")
	if form, err := c.MultipartForm(); err == nil {

		files := form.File["challenge"]

		// Loops through all challenges, if multiple uploaded :
		for _, file := range files {

			//Creates folders for each challenge
			pattern := `([^/]+)\.tar\.gz$`
			regex := regexp.MustCompile(pattern)
			match := regex.FindStringSubmatch(file.Filename)
			folderName = match[1]

			response, challengePath := createFolder(folderName)
			if response == 1 {
				log.Println("Directory already exists with same name")
				return c.SendString("Directory already exists with same name")
			} else if response == 2 {
				log.Println("Issue with creating chall directory.Check permissions")
				return c.SendString("Issue with creating chall directory.Check permissions")
			}

			//Save to disk in that directory
			if err := c.SaveFile(file, fmt.Sprintf("./challenges/%s/%s", folderName, file.Filename)); err != nil {
				return err
			}

			//Create folder inside the challenge folder
			err = os.Mkdir(challengePath+"/"+folderName, 0777)
			if err != nil {
				log.Println("Error in creating folder inside challenge folder")
				return c.SendString("Error in creating folder inside challenge folder")
			}
			//Extract the tar.gz file
			err := archiver.Unarchive("./challenges/"+folderName+"/"+file.Filename, "./challenges/"+folderName)
			if err != nil {
				log.Println("Error in unarchiving", err)
				return c.SendString("Error in unarchiving")
			}

			//Update challenge path to get dockerfile
			utils.BuildDockerImage(folderName, challengePath+"/"+folderName)

			//Get no.of teams and DEPLOY CHALLENGE to each namespace (assuming they exist and /createTeams has been called)
			clusterConfig := g.ClusterConfig
			numberOfTeams := clusterConfig.TeamCount
			res := make([][]string, 0)
			for i := 0; i < int(numberOfTeams); i++ {
				log.Println("-----------Deploying challenge for team: " + strconv.Itoa(i) + " --------")
				teamName := "katana-team-" + strconv.Itoa(i)
				deployment.DeployChallengeToCluster(folderName, teamName, patch, replicas)
				url, err := createServiceForChallenge(folderName, teamName, 3000, i)
				if err != nil {
					res = append(res, []string{teamName, err.Error()})
				} else {
					res = append(res, []string{teamName, url})
				}
			}
			copyChallengeIntoTsuka(challengePath, folderName, challengeType)
			copyFlagDataIntoKashira(challengePath, folderName)
			copyChallengeCheckerIntoKissaki(challengePath, folderName)

			return c.JSON(res)
		}
	}
	log.Println("Ending")
	return c.SendString("Wrong file")
}

func ChallengeUpdate(c *fiber.Ctx) error {
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
		log.Println("Error pulling changes:", err)
	}
	katanaDir, err := utils.GetKatanaRootPath()
	imageName := strings.Replace(dir, "/", "-", -1)

	log.Println("Pull successful for", teamName, ". Building image...")
	firstPatch := !utils.DockerImageExists(imageName)
	utils.BuildDockerImage(imageName, katanaDir+"/teams/"+dir)

	if firstPatch {
		log.Println("First Patch for", teamName)
		deployment.DeployChallengeToCluster(challengeName, teamName, patch, replicas)
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
		}
	}
	log.Println("Image built for", teamName)
	return c.SendString("Challenge updated")
}

func DeleteChallenge(c *fiber.Ctx) error {

	challengeName := c.Params("challengeName")
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
			//panic(err)
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
			continue
			//panic(err)
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

	return c.SendString("Deleted challenge" + challengeName + "in all namespaces.")
}
