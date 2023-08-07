package challengedeployerservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/sdslabs/katana/configs"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func copyChallengeIntoTsuka(dirPath string, challengeName string, challengeType string) error {
	localFilePath := dirPath + "/" + challengeName + ".tar.gz"
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
			URLs: []string{"http://sdslabs@" + utils.GetKatanaLoadbalancer() + ":80" + "/" + path}}
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

func createServiceAndIngressRuleForChallenge(challengeName, teamName string, targetPort int32, teamNumber int) (string, error) {
	kubeclient, _ := utils.GetKubeClient()
	serviceName := challengeName + "-svc"
	teamNamespace := teamName + "-ns"
	port := int32(80)
	selector := map[string]string{
		"app": challengeName,
	}

	utils.CreateService(kubeclient, serviceName, teamNamespace, port, targetPort, selector)

	log.Printf("Created service %s for challenge %s in namespace %s", serviceName, challengeName, teamNamespace)

	// Get team ingress
	ingressName := "team-" + strconv.Itoa(teamNumber) + "-ingress"
	teamIngress, err := kubeclient.NetworkingV1().Ingresses(teamNamespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	additionalRules := networkingv1.IngressRule{
		Host: fmt.Sprintf("%s.%s.%s", challengeName, teamName, configs.KatanaConfig.IngressHost),
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{
					{
						Path: "/",
						PathType: func() *networkingv1.PathType {
							pt := networkingv1.PathTypePrefix
							return &pt
						}(),
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: serviceName,
								Port: networkingv1.ServiceBackendPort{
									Number: port,
								},
							},
						},
					},
				},
			},
		},
	}

	teamIngress.Spec.Rules = append(teamIngress.Spec.Rules, additionalRules)

	_, err = kubeclient.NetworkingV1().Ingresses(teamNamespace).Update(context.Background(), teamIngress, metav1.UpdateOptions{})
	if err != nil {
		return "", err
	}

	log.Printf("Added ingress rule for challenge %s in namespace %s", challengeName, teamNamespace)

	return fmt.Sprintf("%s.%s.%s", challengeName, teamName, configs.KatanaConfig.IngressHost), nil
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
