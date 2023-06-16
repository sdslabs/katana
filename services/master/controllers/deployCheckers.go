package controllers

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/mholt/archiver/v3"
	checker "github.com/sdslabs/katana/services/challengecheckerservice"
)

func DeployCheckers(c *fiber.Ctx) error {
	foldername := ""
	fmt.Println("Starting")
	if form, err := c.MultipartForm(); err == nil {

		//sort this
		if token := form.Value["token"]; len(token) > 0 {
			// Get key value:
			fmt.Println(token[0])
			c.SendString("Test a")
		}
		checkerFiles := form.File["challengeCheckeryOOO"]
		// if len(challengeFiles) != len(checkerFiles) {
		// 	return c.SendString(fmt.Sprintf("No. of challenge files and checker files should be same! Got %d challenge files and %d checker files", len(challengeFiles), len(checkerFiles)))
		// }

		// Loops through all challenges, if multiple uploaded :
		for _, file := range checkerFiles {
			// checkerFile := checkerFiles[i]

			//creates folders for each challenge
			pattern := `([^/]+)\.tar\.gz$`
			regex := regexp.MustCompile(pattern)
			match := regex.FindStringSubmatch(file.Filename)
			foldername = match[1]

			// response, newDirPath := createfolder(foldername)
			// if response == 1 {
			// 	fmt.Println("Directory for challenge already exists with same name")
			// 	return c.SendString("Directory for challenge already exists with same name")
			// } else if response == 2 {
			// 	fmt.Println("Issue with creating challenge directory.Check permissions")
			// 	return c.SendString("Issue with creating challenge directory.Check permissions")
			// }

			response, _ := createfolder(foldername + "-checker")
			if response == 1 {
				fmt.Println("Directory for challenge checker already exists with same name")
				return c.SendString("Directory for challenge checker already exists with same name")
			} else if response == 2 {
				fmt.Println("Issue with creating challenge checkor directory.Check permissions")
				return c.SendString("Issue with creating challenge checkor directory.Check permissions")
			}

			// Save both files to disk in respective directories
			// if err := c.SaveFile(file, fmt.Sprintf("./chall/%s/%s", foldername, file.Filename)); err != nil {
			// 	return err
			// }
			if err := c.SaveFile(file, fmt.Sprintf("./chall/%s/%s", foldername+"-checker", file.Filename)); err != nil {
				return err
			}

			// Extract both tar.gz files
			// err := archiver.Unarchive("./chall/"+foldername+"/"+file.Filename, "./chall/"+foldername)
			// if err != nil {
			// 	fmt.Println("Error in unarchiving", err)
			// 	return c.SendString("Error in unarchiving")
			// }
			err = archiver.Unarchive("./chall/"+foldername+"-checker/"+file.Filename, "./chall/"+foldername+"-checker")
			if err != nil {
				fmt.Println("Error in unarchiving", err)
				return c.SendString("Error in unarchiving")
			}

			// fmt.Println("Building docker image with tag", foldername)
			// buildimage(foldername)
			fmt.Println("Building docker image with tag", foldername+"-checker")
			buildimage(foldername + "-checker")
			fmt.Println("Docker images built successfully")

			//Get no.of teams and DEPLOY CHALLENGE to each namespace (assuming they exist and /createTeams has been called)
			//For only testing this and not the /createTeams route, create 3 namespaces (katana-team-0-ns) (katana-team-1-ns) (katana-team-2-ns) manually
			// clusterConfig := g.ClusterConfig
			// numberOfTeams := clusterConfig.TeamCount
			// res := make([][]string, 0)
			// for i := 0; i < int(numberOfTeams); i++ {
			// 	fmt.Println("-----------Deploying challenge for team: " + strconv.Itoa(i) + " --------")
			// 	team_name := "katana-team-" + strconv.Itoa(i)
			// 	deployer.DeployChallenge(foldername, team_name)
			// 	checker.DeployChallChecker(foldername, "5000", team_name+"-ns")
			// 	url, err := deployer.CreateService(foldername, team_name)
			// 	if err != nil {
			// 		res = append(res, []string{team_name, err.Error()})
			// 	} else {
			// 		res = append(res, []string{team_name, url})
			// 	}
			// }

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

// func DeleteChallenge(c *fiber.Ctx) error {

// 	chall_name := c.Params("chall_name")
// 	//Delete chall directory also ?
// 	fmt.Println("Challenge name is : " + chall_name)

// 	deployer.DeleteChallenge(chall_name)
// 	return c.SendString("Deleted challenge" + chall_name + "in all namespaces.")
// }
