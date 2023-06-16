package challengecheckerservice

import (
	"context"
	"fmt"
	"log"

	g "github.com/sdslabs/katana/configs"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployChallChecker(chall_name, port, team_namespace string) error {
	namespace := "katana"

	if err := GetClient(g.KatanaConfig.KubeConfig); err != nil {
		return err
	}

	jobs := kubeclient.BatchV1beta1().CronJobs(namespace)

	jobSpec := &batchv1beta.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      chall_name + "-checker",
			Namespace: namespace,
		},
		Spec: batchv1beta.CronJobSpec{
			Schedule: "* */2 * * *",
			JobTemplate: batchv1beta.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{

									Name:            chall_name + "-checker",
									Image:           chall_name + "-checker:latest",
									ImagePullPolicy: v1.PullPolicy("Never"),
									Env: []v1.EnvVar{
										{
											Name:  "URL",
											Value: "http://" + chall_name + "." + team_namespace + ".svc.cluster.local:" + port,
										},
									},
								},
							},
							RestartPolicy: v1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}
	cronJob, err := jobs.Get(context.TODO(), chall_name+"-checker", metav1.GetOptions{})
	if cronJob.Name == chall_name+"-checker" {
		fmt.Println("Challenge checker already exists for the challenge " + chall_name + " in namespace " + namespace)
		return nil
	}

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Creating challenge checker (CronJob) for the challenge " + chall_name + " in namespace " + namespace)
	result, err := jobs.Create(context.Background(), jobSpec, metav1.CreateOptions{})

	if err != nil {
		fmt.Println("Failed to create challenge checker for the challenge " + chall_name + " in namespace " + namespace)
		panic(err)
	}

	fmt.Printf("Created cronjob %q.\n", result.GetObjectMeta().GetName()+" in namespace "+namespace+"\n")
	return nil
}
