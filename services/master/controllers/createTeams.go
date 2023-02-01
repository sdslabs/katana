package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTeam() types.CTFTeam {
	ctfTeam := new(types.CTFTeam)
	return *ctfTeam
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
				Name: "ctf-team-" + strconv.Itoa(i),
			},
		}

		_, err = client.CoreV1().Namespaces().Create(c.Context(), nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	return c.SendString("No. of Teams Created:" + c.Params("number") + "\n")
}
