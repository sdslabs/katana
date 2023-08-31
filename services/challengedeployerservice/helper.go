package challengedeployerservice

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
)

func copyChallengeIntoTsuka(dirPath, challengeName string, challengeType string) error {
	srcFilePath := dirPath +"/"+ challengeName + "/" + challengeName
	pathInPod := "/opt/katana/challenge"
	
	filename := challengeName

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
			URLs: []string{"http://sdslabs@" + utils.GetKatanaLoadbalancer() + ":3000" + "/" + path}}
		_, err = repo.CreateRemote(remoteConfig)

		if err != nil {
			log.Println(err)
		}
		podsInTeam, err := utils.GetPods(map[string]string{
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
		
		if err := utils.CopyTarIntoPodNew(pod.Name, g.TeamVmConfig.ContainerName, pathInPod, srcFilePath , pod.Namespace); err != nil {
			log.Println(err)
			return err
		}
	}
	
	return nil
}


func copyChallengeIntoTsukaOriginal(dirPath string, challengeName string, challengeType string) error {

	localFilePath := dirPath +"/"+ challengeName + ".tar.gz"
	pathInPod := "/opt/katana/katana_" + challengeType + "_" + challengeName + ".tar.gz"
	log.Println("Testing" + localFilePath + "....and..." + pathInPod)

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
			URLs: []string{"http://sdslabs@" + utils.GetKatanaLoadbalancer() + ":3000" + "/" + path}}
		_, err = repo.CreateRemote(remoteConfig)

		if err != nil {
			log.Println(err)
		}
		podsInTeam, err := utils.GetPods(map[string]string{
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

func createServiceForChallenge(challengeName, teamName string, targetPort int32, teamNumber int) (string, error) {
	kubeclient, _ := utils.GetKubeClient()
	serviceName := challengeName + "-svc-" + strconv.Itoa(teamNumber)
	teamNamespace := teamName + "-ns"
	port := int32(80)
	selector := map[string]string{
		"app": challengeName,
	}

	utils.CreateService(kubeclient, serviceName, teamNamespace, port, targetPort, selector)

	log.Printf("Created service %s for challenge %s in namespace %s", serviceName, challengeName, teamNamespace)

	return serviceName, nil
}

func createFolder(challengeName string) (message int, challengePath string) {

	basePath, _ := os.Getwd()
	dirPath := basePath + "/challenges" //basepath is .../katana

	// Open the challenges directory to check if it exists , create if not
	dir, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Challenges directory does not exist ,creating directory")
			os.Mkdir(dirPath, 0777)
		} else if os.IsPermission(err) {
			log.Println("Error opening challenge directory. Permission Issue", err)
			//Permission issue
			return 2, challengePath
		} else {
			log.Println("Error opening challenge directory:", err)
			//Some other error
			return 2, challengePath
		}
	}
	defer dir.Close()

	// Create a new challenge directory to keep challenge
	challengePath = dirPath + "/" + challengeName
	log.Println("Creating directory :", challengeName)
	err = os.Mkdir(challengePath, 0777)
	if err != nil {
		//Directory already exists with same name
		return 1, challengePath
	}
	//Successfully created directory
	return 0, challengePath
}

func copyChallengeCheckerIntoKissaki(dirPath string, challengeName string, challengeType string) error {
	srcFilePath := dirPath +"/"+ challengeName + "/" + challengeName + "_cc" 
	pathInPod := "/opt/kissaki/challenge-data"

	log.Println("Testing" + srcFilePath + "....and..." + pathInPod)
		
		if err := utils.CopyTarIntoPodNew("kissaki-0", "kissaki", pathInPod, srcFilePath, "katana"); err != nil {
			log.Println(err)
			return err
		}
	return nil
}

func copyFlagDataIntoKashira(dirPath string, challengeName string, challengeType string) error {
	srcFilePath := dirPath +"/"+ challengeName + "/" + "flag_data" 
	pathInPod := "/opt/kashira/flag-data"
	log.Println("Testing" + srcFilePath + "....and..." + pathInPod)
		
		if err := utils.CopyTarIntoPodNew("kashira-0", "kashira", pathInPod, srcFilePath, "katana"); err != nil {
			log.Println(err)
			return err
		}
	return nil
}