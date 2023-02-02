package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/lib/utils"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTeam(namespace string, i string) *coreV1.Pod {
	// ctfTeam := new(types.CTFTeam)
	log.Println("Creating pod for team " + i)
	teamPod := &coreV1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "katana-team-" + i,
			Namespace: namespace,
		},
		Spec: coreV1.PodSpec{
			Containers: []coreV1.Container{
				{
					Name:  "teamvm",
					Image: "scar26/sdskatanad",
					Env: []coreV1.EnvVar{
						{
							Name: "CHALLENGE_DIR",
							ValueFrom: &coreV1.EnvVarSource{
								ConfigMapKeyRef: &coreV1.ConfigMapKeySelector{
									LocalObjectReference: coreV1.LocalObjectReference{
										Name: "teamvm-config",
									},
									Key: "challenge_dir",
								},
							},
						},
						{
							Name: "TMP_DIR",
							ValueFrom: &coreV1.EnvVarSource{
								ConfigMapKeyRef: &coreV1.ConfigMapKeySelector{
									LocalObjectReference: coreV1.LocalObjectReference{
										Name: "teamvm-config",
									},
									Key: "tmp_dir",
								},
							},
						},
						{
							Name: "INIT_FILE",
							ValueFrom: &coreV1.EnvVarSource{
								ConfigMapKeyRef: &coreV1.ConfigMapKeySelector{
									LocalObjectReference: coreV1.LocalObjectReference{
										Name: "teamvm-config",
									},
									Key: "init_file",
								},
							},
						},
						{
							Name: "DAEMON_PORT",
							ValueFrom: &coreV1.EnvVarSource{
								ConfigMapKeyRef: &coreV1.ConfigMapKeySelector{
									LocalObjectReference: coreV1.LocalObjectReference{
										Name: "teamvm-config",
									},
									Key: "daemon_port",
								},
							},
						},
					},
				},
			},
		},
	}

	return teamPod
}

func CreateTeams(c *fiber.Ctx) error {
	client, err := utils.GetKubeClient()
	if err != nil {
		log.Println(err)
	}
	noOfTeams, err := strconv.Atoi(c.Params("number"))

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < noOfTeams; i++ {
		log.Println("Creating Team: " + strconv.Itoa(i))
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "katana-team-ns-" + strconv.Itoa(i),
			},
		}

		_, err = client.CoreV1().Namespaces().Create(c.Context(), nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}
		pod := CreateTeam("katana-team-ns-"+strconv.Itoa(i), strconv.Itoa(i))
		_, err = client.CoreV1().Pods(pod.Namespace).Create(c.Context(), pod, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	return c.SendString("No. of Teams Created:" + c.Params("number") + "\n")
}
