package challengecheckerservice

import (
	"context"
	"fmt"
	"time"

	g "github.com/sdslabs/katana/configs"
	utils "github.com/sdslabs/katana/lib/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployChallChecker(challName, port, teamNamespace string) error {
	namespace := "katana"
	kubeclient, _ := utils.GetKubeClient(g.KatanaConfig.KubeConfig)

	cronJobSpec := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      challName + "-checker",
			Namespace: namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          g.KatanaConfig.CronJobSchedule, // Every minute.
			ConcurrencyPolicy: batchv1.ForbidConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            challName + "-checker",
									Image:           challName + "-checker:latest",
									ImagePullPolicy: corev1.PullPolicy("Never"),
									Env: []corev1.EnvVar{
										{
											Name:  "URL",
											Value: "http://" + challName + "." + teamNamespace + ".svc.cluster.local:" + port,
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
	cronJob, _ := kubeclient.BatchV1().CronJobs(namespace).Get(context.TODO(), challName+"-checker", metav1.GetOptions{})
	if cronJob.Name == challName+"-checker" {
		fmt.Println("Challenge checker already exists for the challenge " + challName + " in namespace " + namespace)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(g.KatanaConfig.TimeOut)*time.Second)
	defer cancel()

	_, err := kubeclient.BatchV1().CronJobs(namespace).Create(ctx, cronJobSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Cron job created successfully!")
	return nil
}
