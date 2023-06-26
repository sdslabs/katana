package controllers

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	archiver "github.com/mholt/archiver/v3"
	checker "github.com/sdslabs/katana/services/challengecheckerservice"
)

func DeployCheckers(c *fiber.Ctx) error {
	folderName := ""
	fmt.Println("Starting")
	if form, err := c.MultipartForm(); err == nil {
		checkerFiles := form.File["challengeChecker"]

		// Loops through all challenges, if multiple uploaded :
		for _, file := range checkerFiles {
			pattern := `([^/]+)\.tar\.gz$`
			regex := regexp.MustCompile(pattern)
			match := regex.FindStringSubmatch(file.Filename)
			folderName = match[1]

			response, _ := createfolder(folderName + "-checker")
			if response == 1 {
				fmt.Println("Directory for challenge checker already exists with same name")
				return c.SendString("Directory for challenge checker already exists with same name")
			} else if response == 2 {
				fmt.Println("Issue with creating challenge checkor directory.Check permissions")
				return c.SendString("Issue with creating challenge checkor directory.Check permissions")
			}

			if err := c.SaveFile(file, fmt.Sprintf("./challenges/%s/%s", folderName+"-checker", file.Filename)); err != nil {
				return err
			}

			err = archiver.Unarchive("./challenges/"+folderName+"-checker/"+file.Filename, "./chall/"+folderName+"-checker")
			if err != nil {
				fmt.Println("Error in unarchiving", err)
				return c.SendString("Error in unarchiving")
			}

			fmt.Println("Building docker image with tag", folderName+"-checker")
			buildimage(folderName + "-checker")
			fmt.Println("Docker images built successfully")

			checker.DeployChallChecker(folderName, "5000", "katana-team-0-ns")
			if err != nil {
				return c.SendString(err.Error())
			} else {
				// return c.SendString(url)
				return c.SendString("Checker deployed successfully")
			}
		}
	}
	fmt.Println("Ending")
	return c.SendString("Wrong file")
}
