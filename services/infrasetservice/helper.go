package infrasetservice

import (
	"log"
	"os"
	"strings"

	utils "github.com/sdslabs/katana/lib/utils"
	ssh "github.com/sdslabs/katana/services/sshproviderservice"
)

func generateCertsforHarbor() {
	path, _ := os.Getwd()
	path = path + "/lib/harbor/certs"

	// Delete the directory if it already exists
	if _, err := os.Stat(path); os.IsExist(err) {
		errDir := os.RemoveAll(path)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		errDir := os.Mkdir(path, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	// Generate the certificates
	if err := utils.GenerateCerts("harbor.katana.local", path); err != nil {
		log.Fatal(err)
	}
}

func SSH(noOfTeams int) {
	ssh.CreateTeams(noOfTeams)
	startServer()
}

func startServer() {
	x := ssh.Server()
	go func() {
		x.ListenAndServe()
	}()
	log.Println("Server up and running")
	for {
	}
}

func buildKatanaServices() {
	katanaServicesDir := "./katana-services"

	services, err := os.ReadDir(katanaServicesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range services {
		if service.Name() == ".github" {
			continue
		}
		if service.IsDir() {
			log.Println("Building " + service.Name())
			imageName := strings.ToLower(service.Name())
			utils.BuildDockerImage(imageName, katanaServicesDir+"/"+service.Name())
		}
	}
}
