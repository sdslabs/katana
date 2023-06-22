package controllers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mholt/archiver/v3"
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

func buildimage(foldername string) {
	// Build the challenge with Dockerfile
	dirPath, _ := os.Getwd()
	fmt.Println("Dockerfile for the image is at :")
	fmt.Println(dirPath + "/challenges/" + foldername + "/" + foldername)
	cmd := exec.Command("docker", "build", "-t", foldername, dirPath+"/chall/"+foldername+"/"+foldername)
	cmd2 := exec.Command("minikube", "image", "load", foldername)
	cmd.Run()
	cmd2.Run()
}

func createfolder(challengename string) (message int, newDirPath string) {

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
			return 2, newDirPath
		} else {
			fmt.Println("Error opening challenge directory:", err)
			//Some other error
			return 2, newDirPath
		}
	}
	defer dir.Close()

	// Create a new challenge directory to keep challenge
	newDirPath = dirPath + "/" + challengename
	fmt.Println("Creating directory :", challengename)
	err = os.Mkdir(newDirPath, 0777)
	if err != nil {
		//Directory already exists with same name
		return 1, newDirPath
	}
	//Successfully created directory
	return 0, newDirPath
}

func Deploy(c *fiber.Ctx) error {

	challengetype := "web"
	foldername := ""
	fmt.Println("Starting")
	if form, err := c.MultipartForm(); err == nil {

		files := form.File["challenge"]

		// Loops through all challenges, if multiple uploaded :
		for _, file := range files {

			//creates folders for each challenge
			pattern := `([^/]+)\.tar\.gz$`
			regex := regexp.MustCompile(pattern)
			match := regex.FindStringSubmatch(file.Filename)
			foldername = match[1]

			response, newDirPath := createfolder(foldername)
			if response == 1 {
				fmt.Println("Directory already exists with same name")
				return c.SendString("Directory already exists with same name")
			} else if response == 2 {
				fmt.Println("Issue with creating chall directory.Check permissions")
				return c.SendString("Issue with creating chall directory.Check permissions")
			}

			//save to disk in that directory
			if err := c.SaveFile(file, fmt.Sprintf("./challenges/%s/%s", foldername, file.Filename)); err != nil {
				return err
			}

			//extract the tar.gz file
			err := archiver.Unarchive("./challenges/"+foldername+"/"+file.Filename, "./chall/"+foldername)
			if err != nil {
				fmt.Println("Error in unarchiving", err)
				return c.SendString("Error in unarchiving")
			}

			fmt.Println("Building docker image with tag", foldername)
			buildimage(foldername)
			fmt.Println("Docker image built successfully")

			//Get no.of teams and DEPLOY CHALLENGE to each namespace (assuming they exist and /createTeams has been called)
			//For only testing this and not the /createTeams route, create 3 namespaces (katana-team-0-ns) (katana-team-1-ns) (katana-team-2-ns) manually
			clusterConfig := g.ClusterConfig
			numberOfTeams := clusterConfig.TeamCount
			res := make([][]string, 0)
			for i := 0; i < int(numberOfTeams); i++ {
				fmt.Println("-----------Deploying challenge for team: " + strconv.Itoa(i) + " --------")
				team_name := "katana-team-" + strconv.Itoa(i)
				utils.DeployChallenge(foldername, team_name, false, 1)
				url, err := deployer.CreateService(foldername, team_name)
				if err != nil {
					res = append(res, []string{team_name, err.Error()})
				} else {
					res = append(res, []string{team_name, url})
				}
			}

			//Copy challenge in pods and etc.
			challcopy(newDirPath, foldername, challengetype)

			return c.JSON(res)
		}
	}
	fmt.Println("Ending")

	//In case only want to test the deployer wih empty post request
	// foldername = "notekeeper"
	// clusterConfig := g.ClusterConfig
	// numberOfTeams := clusterConfig.TeamCount
	// for i := 0; i < int(numberOfTeams); i++ {
	// 	deployer.DeployChallenge(foldername, "team-"+strconv.Itoa(i))
	//  deployer.CreateService(foldername, "team-"+strconv.Itoa(i))
	// }

	//deployer.CreateService("notekeeper", "team-0")

	return c.SendString("Wrong file")
}

func DeleteChallenge(c *fiber.Ctx) error {

	chall_name := c.Params("chall_name")
	//Delete chall directory also ?
	fmt.Println("Challenge name is : " + chall_name)

	deployer.DeleteChallenge(chall_name)
	return c.SendString("Deleted challenge" + chall_name + "in all namespaces.")
}
