package flaghandlerservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"os/exec"

	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func checkSubmission(submission *types.Submission) bool {
	_, err := mongo.FetchSubmission(submission.SubmittedBy, submission.Flag)
	return err != nil
}

func submitFlag(value string, team types.CTFTeam) (bool, int) {
	flag, err := mongo.FetchFlag(value)
	var points int
	if err != nil {
		return false, 0
	}
	if flag.TeamID == team.Index {
		return false, 0
	}
	submission := &types.Submission{}
	submission.ChallengeID = flag.ChallengeID
	submission.Flag = flag.Value
	submission.SubmittedBy = team.Index
	submission.Time = time.Now()
	if !checkSubmission(submission) {
		return false, 0
	}
	if res, err := mongo.FetchChallenge(flag.ChallengeID); err != nil {
		return false, 0
	} else {
		team.Score = team.Score + res.Points
		points = res.Points
	}
	if err := mongo.UpdateOne(context.Background(), mongo.TeamsCollection, bson.M{"id": team.Index}, team, options.FindOneAndUpdate().SetUpsert(false)); err != nil {
		return false, 0
	}
	if _, err := mongo.InsertOne(context.Background(), mongo.SubmissionsCollection, submission); err != nil {
		return false, 0
	}
	return true, points
}

func getflag(container string, script string) {
	// hard coded script name and container, will become parameter later
	pathToCfg := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	fmt.Println(g.KatanaConfig.KubeNameSpace)
	config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, pod := range pods.Items {
		podName := pod.Name
		cmd := exec.Command("kubectl", "cp", script, podName+":"+script, "-c", container)
		err := cmd.Run()
		if err != nil {
			fmt.Println("error aaya hai bhaiya")
			fmt.Println(err)
		}
		out, erro := exec.Command("kubectl", "exec", podName, "-c", container, "--", "bash", "-c", "./"+script).Output()
		if erro != nil {
			fmt.Println(err)
		}
		output := string(out)
		_, err = mongo.FetchFlag(output)
		if err != nil {
			fmt.Println(podName + " ka container " + container + "is not running")
		}
		err = exec.Command("kubectl", "exec", podName, "-c", container, "--", "rm", script).Run()
		if err != nil {
			fmt.Println("error aaya hai bhaiya")
			fmt.Println(err)
		}
	}
	return
}
