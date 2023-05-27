package controllers

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gofiber/fiber/v2"
	deployer "github.com/sdslabs/katana/services/challengedeployerservice"
)

// Run testdeploy for the basic pod copying test,change in 59/59 and 100/101 line
func testdeploy(dirPath, challengename, challengetype string) {

	dirPath, _ = os.Getwd()
	pattern := `^(.*)/[^/]+/?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(dirPath)
	parentPath := matches[1]

	fmt.Println(parentPath)

	localFilePath := parentPath + "/" + challengename + ".tar.gz"
	pathInPod := "/opt/katana/katana_" + challengetype + "_" + challengename + ".tar.gz"

	deployer.DeployToAll(localFilePath, pathInPod)

}

func challdeploy(dirPath, challengename, challengetype string) {

	//TO-DO : Put the challenge in dirPath and remove testdeploy
	localFilePath := dirPath + "/" + challengename + ".tar.gz"
	pathInPod := "/opt/katana/katana_" + challengetype + "_" + challengename + ".tar.gz"
	//fmt.Println("Testing" + localFilePath + "....and..." + pathInPod)
	deployer.DeployToAll(localFilePath, pathInPod)

}

func Deploy(c *fiber.Ctx) error {

	challengetype := "web"
	challengename := "notekeeper"
	//TODO : Change /Deploy to postreq and update the challengename whatever the name of file entered is.

	basePath, _ := os.Getwd()
	dirPath := basePath + "/chall"

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Chall directory does not exist ,creating directory")
			os.Mkdir(dirPath, 0777)
		} else if os.IsPermission(err) {
			fmt.Println("Error opening challenge directory. Permission Issue", err)
		} else {
			fmt.Println("Error opening challenge directory:", err)
		}
	}
	defer dir.Close()

	// Create a new challenge directory to keep challenge
	newDirPath := dirPath + "/" + challengename
	fmt.Println("Creating directory :", challengename)
	err = os.Mkdir(newDirPath, 0777)
	if err != nil {
		return c.SendString("A challenge with the same name exists")
	}

	challdeploy(newDirPath, challengename, challengetype)
	//testdeploy(dirPath, challengename, challengetype)

	return c.SendString("Deployed")
}
