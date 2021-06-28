package sshproviderservice

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gliderlabs/ssh"
	g "github.com/sdslabs/katana/configs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	kubeClientset *kubernetes.Clientset
	kubeConfig    *rest.Config
	execCmd       = []string{"/bin/bash"}
)

func sessionHandler(s ssh.Session) {
	kubeclient := kubeClientset.CoreV1().RESTClient()

	podName := s.User()

	req := kubeclient.Post().Resource("pods").Name(podName).Namespace(g.KatanaConfig.KubeNameSpace).SubResource("exec")

	option := &v1.PodExecOptions{
		Command: execCmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		fmt.Fprintf(s, "ERROR: %s", err.Error())
		s.Exit(1)
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  s,
		Stdout: s,
		Stderr: s,
	})

	if err != nil {
		fmt.Fprintf(s, "ERROR: %s", err.Error())
		s.Exit(1)
	}
}

func Server() {
	ssh.Handle(sessionHandler)
	log.Println("starting ssh server on port 2222")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func init() {
	var pathToCfg string
	if g.KatanaConfig.KubeConfig == "" {
		pathToCfg = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	} else {
		pathToCfg = g.KatanaConfig.KubeConfig
	}

	config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		log.Fatal(err)
	}
	kubeConfig = config

	kubeClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
}
