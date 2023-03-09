package sshproviderservice

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/gliderlabs/ssh"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	kubeClientset *kubernetes.Clientset
	kubeConfig    *rest.Config
	execCmd       = []string{"/bin/bash"}
)

func sessionHandler(s ssh.Session) {
	kubeclient := kubeClientset.CoreV1().RESTClient()

	Namespace := s.User()

	req := kubeclient.Post().Resource("pods").Name("katana-team-master-pod-0").Namespace(Namespace + "-ns").SubResource("exec")

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
		if err := s.Exit(1); err != nil {
			fmt.Printf("Failed to connect to server: %s", err)
		}
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  s,
		Stdout: s,
		Stderr: s,
	})

	if err != nil {
		fmt.Fprintf(s, "ERROR: %s", err.Error())
		if err := s.Exit(1); err != nil {
			fmt.Printf("Failed to stream data to remote server: %s", err)
		}
	}
}

func passwordHandler(s ssh.Context, password string) bool {
	team, err := mongo.FetchSingleTeam(s.User())
	if err != nil {
		return false
	}
	return utils.CompareHashWithPassword(team.Password, password)
}

func publicKeyHandler(s ssh.Context, key ssh.PublicKey) bool {
	team, err := mongo.FetchSingleTeam(s.User())
	if err != nil {
		return false
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(team.PublicKey))
	if err != nil {
		return false
	}

	return bytes.Equal(publicKey.Marshal(), key.Marshal())
}

func Server() *ssh.Server {
	return &ssh.Server{
		Addr:             net.JoinHostPort(g.SSHProviderConfig.Host, fmt.Sprintf("%d", g.SSHProviderConfig.Port)),
		Handler:          sessionHandler,
		PasswordHandler:  passwordHandler,
		PublicKeyHandler: publicKeyHandler,
	}
}

func init() {
	var err error
	kubeConfig, err = utils.GetKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeClientset, err = utils.GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}
}
