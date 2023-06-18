package challengecheckerservice

import (
	"context"
	"fmt"
	"time"

	g "github.com/sdslabs/katana/configs"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployChallChecker(chall_name, port, team_namespace string) error {
	namespace := "katana"

	if err := GetClient(g.KatanaConfig.KubeConfig); err != nil {
		return err
	}
	cronJobSpec := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      chall_name + "-checker",
			Namespace: namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          "*/1 * * * *", // Every minute.
			ConcurrencyPolicy: batchv1.ForbidConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            chall_name + "-checker",
									Image:           chall_name + "-checker:latest",
									ImagePullPolicy: corev1.PullPolicy("Never"),
									Env: []corev1.EnvVar{
										{
											Name:  "URL",
											Value: "http://" + chall_name + "." + team_namespace + ".svc.cluster.local:" + port,
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	cronJob, _ := kubeclient.BatchV1().CronJobs(namespace).Get(context.TODO(), chall_name+"-checker", metav1.GetOptions{})
	if cronJob.Name == chall_name+"-checker" {
		fmt.Println("Challenge checker already exists for the challenge " + chall_name + " in namespace " + namespace)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := kubeclient.BatchV1().CronJobs(namespace).Create(ctx, cronJobSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Cron job created successfully!")

	// if err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println("Creating challenge checker (CronJob) for the challenge " + chall_name + " in namespace " + namespace)
	// result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})

	// if err != nil {
	// 	fmt.Println("Failed to create challenge checker for the challenge " + chall_name + " in namespace " + namespace)
	// 	panic(err)
	// }

	// fmt.Printf("Created cronjob %q.\n", result.GetObjectMeta().GetName()+" in namespace "+namespace+"\n")
	return nil
}
