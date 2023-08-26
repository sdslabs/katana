package flaghandlerservice


	import (
		"log"
		
		g "github.com/sdslabs/katana/configs"
		"github.com/sdslabs/katana/lib/utils"
		_"k8s.io/api/core/v1"
	)
	
	
func copyFlagSetterIntoKashira(dirPath string, flagSetter string) error {
	podName := "kashira-0"
	namespace := "katana"
	localFilePath := dirPath + "/" + flagSetter + ".tar.gz"
	pathInPod := "/opt/kashira" + flagSetter + ".tar.gz"
	log.Println("Testing" + localFilePath + "....and..." + pathInPod)

	//regex to find challenge name since localFilePath[12:22] is hardcoded
	// regexPattern := `\/([^\/]+)\.tar\.gz$`
	// regex := regexp.MustCompile(regexPattern)
	// matches := regex.FindStringSubmatch(localFilePath)
	// filename := matches[1]

	// Get pods from different namespaces
	
		// Copy file into pod
		if err := utils.CopyIntoPod(podName, g.TeamVmConfig.ContainerName, pathInPod, localFilePath, namespace); err != nil {
			log.Println(err)
			return err
		}
	

	return nil
}