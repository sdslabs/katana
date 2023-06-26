package controllers

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	archiver "github.com/mholt/archiver/v3"
	checker "github.com/sdslabs/katana/services/challengecheckerservice"
)

func DeployCheckers(c *fiber.Ctx) error {
	foldername := ""
	fmt.Println("Starting")
	if form, err := c.MultipartForm(); err == nil {
		checkerFiles := form.File["challengeChecker"]

		// Loops through all challenges, if multiple uploaded :
		for _, file := range checkerFiles {
			pattern := `([^/]+)\.tar\.gz$`
			regex := regexp.MustCompile(pattern)
			match := regex.FindStringSubmatch(file.Filename)
			foldername = match[1]

			response, _ := createfolder(foldername + "-checker")
			if response == 1 {
				fmt.Println("Directory for challenge checker already exists with same name")
				return c.SendString("Directory for challenge checker already exists with same name")
			} else if response == 2 {
				fmt.Println("Issue with creating challenge checkor directory.Check permissions")
				return c.SendString("Issue with creating challenge checkor directory.Check permissions")
			}

			if err := c.SaveFile(file, fmt.Sprintf("./challenges/%s/%s", foldername+"-checker", file.Filename)); err != nil {
				return err
			}

			err = archiver.Unarchive("./challenges/"+foldername+"-checker/"+file.Filename, "./chall/"+foldername+"-checker")
			if err != nil {
				fmt.Println("Error in unarchiving", err)
				return c.SendString("Error in unarchiving")
			}

			fmt.Println("Building docker image with tag", foldername+"-checker")
			buildimage(foldername + "-checker")
			fmt.Println("Docker images built successfully")

			checker.DeployChallChecker(foldername, "5000", "katana-team-0-ns")
			// url, err := checker.CreateService(foldername, "katana-team-0")
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
