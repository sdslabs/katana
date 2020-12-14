package challengedeployerservice

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/sdslabs/katana/configs"
	pb "github.com/sdslabs/katana/lib/foundry/protos/vmmanager"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

var config configs.ChallengeDeployerConfig
var kubeclient *kubernetes.Clientset

type server struct {
	pb.UnimplementedVMFoundryServer
}

func configureServer(bindings pb.VMFoundryServer) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterVMFoundryServer(srv, bindings)
	return srv
}

func Server() error {
	srv := configureServer(&server{})

	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic(err)
	}

	return srv.Serve(lis)
}

func test() {
	fmt.Println("Testing sageo")
	// auth := &githttp.BasicAuth{
	// 	Username: config.Username,
	// 	Password: config.AccessToken,
	// }
	// err := clone("https://github.com/Scar26/Im-in", "imin", auth)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	chal, err := os.Open("challenges/imin.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer chal.Close()
	params := make(map[string]string)
	params["challenge_name"] = "Im-in"
	params["key2"] = "val2"
	if err = sendFile(chal, params, "imin.zip", "http://localhost:1234/grab"); err != nil {
		log.Fatal(err)
	}
}

func NewService() error {
	fmt.Println("Starting challenge deployer")
	config = configs.ChallengeDeployerServiceConfig
	fmt.Println(configs.KatanaConfig)
	return nil
}
