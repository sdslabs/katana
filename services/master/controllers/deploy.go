package controllers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	archiver "github.com/mholt/archiver/v3"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"

	deployer "github.com/sdslabs/katana/services/challengedeployerservice"
)

func challcopy(dirPath, challengename, challengetype string) {

	localFilePath := dirPath + "/" + challengename + ".tar.gz"
	pathInPod := "/opt/katana/katana_" + challengetype + "_" + challengename + ".tar.gz"
	fmt.Println("Testing" + localFilePath + "....and..." + pathInPod)
	deployer.CopyInPod(localFilePath, pathInPod)

}

func createfolder(challengename string) (message int, challengePath string) {

	basePath, _ := os.Getwd()
	dirPath := basePath + "/challenges" //basepath is .../katana

	// Open the challenges directory to check if it exists , create if not
	dir, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Challenges directory does not exist ,creating directory")
			os.Mkdir(dirPath, 0777)
		} else if os.IsPermission(err) {
			fmt.Println("Error opening challenge directory. Permission Issue", err)
			//Permission issue
			return 2, challengePath
		} else {
			fmt.Println("Error opening challenge directory:", err)
			//Some other error
			return 2, challengePath
		}
	}
	defer dir.Close()

	// Create a new challenge directory to keep challenge
	challengePath = dirPath + "/" + challengename
	fmt.Println("Creating directory :", challengename)
	err = os.Mkdir(challengePath, 0777)
	if err != nil {
		//Directory already exists with same name
		return 1, challengePath
	}
	//Successfully created directory
	return 0, challengePath
}

func buildimage(folderName string) {
	// Build the challenge with Dockerfile
	dirPath, _ := os.Getwd()
	fmt.Println("Dockerfile for the image is at :")
	fmt.Println(dirPath + "/challenges/" + folderName + "/" + folderName)
	cmd := exec.Command("docker", "build", "-t", g.HarborConfig.Hostname+"/katana/"+folderName, dirPath+"/chall/"+folderName+"/"+folderName)
	cmd2 := exec.Command("docker", "push", g.HarborConfig.Hostname+"/katana/"+folderName)
	cmd.Run()
	cmd2.Run()
}

func Deploy(c *fiber.Ctx) error {

	challengetype := "web"
	folderName := ""
	patch := false
	replicas := g.KatanaConfig.TeamDeployement
	fmt.Println("Starting")
	if form, err := c.MultipartForm(); err == nil {

		files := form.File["challenge"]

		// Loops through all challenges, if multiple uploaded :
		for _, file := range files {

			//creates folders for each challenge
			pattern := `([^/]+)\.tar\.gz$`
			regex := regexp.MustCompile(pattern)
			match := regex.FindStringSubmatch(file.Filename)
			folderName = match[1]

			response, challengePath := createfolder(folderName)
			if response == 1 {
				fmt.Println("Directory already exists with same name")
				return c.SendString("Directory already exists with same name")
			} else if response == 2 {
				fmt.Println("Issue with creating chall directory.Check permissions")
				return c.SendString("Issue with creating chall directory.Check permissions")
			}

			//save to disk in that directory
			if err := c.SaveFile(file, fmt.Sprintf("./challenges/%s/%s", folderName, file.Filename)); err != nil {
				return err
			}

			//Create folder inside the challenge folder
			err = os.Mkdir(challengePath+"/"+folderName, 0777)
			if err != nil {
				fmt.Println("Error in creating folder inside challenge folder")
				return c.SendString("Error in creating folder inside challenge folder")
			}

			//extract the tar.gz file
			err := archiver.Unarchive("./challenges/"+folderName+"/"+file.Filename, "./challenges/"+folderName)
			if err != nil {
				fmt.Println("Error in unarchiving", err)
				return c.SendString("Error in unarchiving")
			}

			//Update challenge path to get dockerfile
			challengePath = challengePath + "/" + folderName

			utils.BuildDockerImage(folderName, challengePath)

			//Get no.of teams and DEPLOY CHALLENGE to each namespace (assuming they exist and /createTeams has been called)
			clusterConfig := g.ClusterConfig
			numberOfTeams := clusterConfig.TeamCount
			res := make([][]string, 0)
			for i := 0; i < int(numberOfTeams); i++ {
				fmt.Println("-----------Deploying challenge for team: " + strconv.Itoa(i) + " --------")
				teamName := "katana-team-" + strconv.Itoa(i)
				utils.DeployChallenge(folderName, teamName, patch, replicas)
				url, err := deployer.CreateService(folderName, teamName)
				if err != nil {
					res = append(res, []string{teamName, err.Error()})
				} else {
					res = append(res, []string{teamName, url})
				}
			}

			//Copy challenge in pods and etc.
			challcopy(challengePath, folderName, challengetype)

			return c.JSON(res)
		}
	}
	fmt.Println("Ending")

	return c.SendString("Wrong file")
}

func DeleteChallenge(c *fiber.Ctx) error {

	chall_name := c.Params("chall_name")
	//Delete chall directory also ?
	fmt.Println("Challenge name is : " + chall_name)

	deployer.DeleteChallenge(chall_name)
	return c.SendString("Deleted challenge" + chall_name + "in all namespaces.")
}
