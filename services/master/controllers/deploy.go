package controllers

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

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
	fmt.Println("Testing" + localFilePath + "....and..." + pathInPod)
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
			fmt.Println("Creating chall1 directory")
			fmt.Println(dirPath)
			os.Mkdir(dirPath, 0777)
			dirPath = dirPath + "/chall1"
			fmt.Println(dirPath)
			err = os.Mkdir(dirPath, 0777)
			//challdeploy(dirPath,challengename,challengetype)
			testdeploy(dirPath, challengename, challengetype)
			//Direct call to create dir and create chall1 and deploy cause the remaining code not req to run for first repo
			return c.SendString("Deployed")
		} else if os.IsPermission(err) {
			fmt.Println("Error opening challenge directory. Permission Issue", err)
		} else {
			fmt.Println("Error opening challenge directory:", err)
		}
	}
	defer dir.Close()

	// Read directory entries
	entries, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading directory entries:", err)
		return c.SendString("Error reading directory entries:")
	}

	//Gets the largest number in the chall directory
	re := regexp.MustCompile(`chall(\d+)`)
	maxNum := 0
	for _, element := range entries {
		name := element.Name()
		match := re.FindStringSubmatch(name)
		//fmt.Println(match)
		num, _ := strconv.Atoi(match[1])
		if num > maxNum {
			maxNum = num
		}
	}

	// Create a new directory with the next number
	newDirName := fmt.Sprintf("chall%d", maxNum+1)
	fmt.Println("Creating directory:", newDirName)
	newDirPath := dirPath + "/" + newDirName
	err = os.Mkdir(newDirPath, 0777)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return c.SendString("creating directory")
	}

	//challdeploy(newDirPath,challengename,challengetype)
	testdeploy(dirPath, challengename, challengetype)

	return c.SendString("Deployed")
}
