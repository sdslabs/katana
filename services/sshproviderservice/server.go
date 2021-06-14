package sshproviderservice

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/gliderlabs/ssh"
	g "github.com/sdslabs/katana/configs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	config       g.ChallengeDeployerConfig
	katanaConfig *g.KatanaCfg
	kubeclient   *kubernetes.Clientset
)

func ExecCmdExample(s ssh.Session) {

	pathToCfg := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		log.Fatal(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	cmd := []string{
		"/bin/bash",
	}

	pods, err := client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	podName := pods.Items[0].Name

	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace("default").SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Fatal(err)
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  s,
		Stdout: s,
		Stderr: s,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(err)
}

func Server() {
	ssh.Handle(ExecCmdExample)
	log.Println("starting ssh server on port 2222")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
