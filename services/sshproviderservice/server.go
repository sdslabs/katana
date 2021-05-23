package sshproviderservice

import (
	"fmt"
	"log"

	"github.com/gliderlabs/ssh"
	g "github.com/sdslabs/katana/configs"
	"k8s.io/client-go/kubernetes"
)

var (
	config       g.ChallengeDeployerConfig
	katanaConfig *g.KatanaCfg
	kubeclient   *kubernetes.Clientset
)

func Server() {
	ssh.Handle(func(s ssh.Session) {
		s.Write([]byte(fmt.Sprintf("hello %s\n", s.User())))
	})
	log.Println("starting ssh server on port 2222")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
