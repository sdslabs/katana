package challengedeployerservice

import (
	"context"
	"fmt"
	"net"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	g "github.com/sdslabs/katana/configs"
	pb "github.com/sdslabs/katana/lib/foundry/protos/challengemanager"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

var (
	config       g.ChallengeDeployerConfig
	katanaConfig *g.KatanaCfg
	kubeclient   *kubernetes.Clientset
)

type server struct{}

func (s *server) Deploy(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	auth := &githttp.BasicAuth{
		Username: request.GetUname(),
		Password: request.GetToken(),
	}
	challName := request.GetChallname()
	repo := request.GetRepo()
	if err := clone(repo, challName, auth); err != nil {
		return &pb.Response{Success: false}, err
	}

	if err := broadcast(fmt.Sprintf("%s.zip", challName)); err != nil {
		return &pb.Response{Success: false}, err
	}

	return &pb.Response{Success: true}, nil
}

func configureServer(bindings pb.ChallengeFoundryServer) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterChallengeFoundryServer(srv, bindings)
	return srv
}

func Server() error {
	srv := configureServer(&server{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}

	return srv.Serve(lis)
}

func NewService() error {
	fmt.Println("Initiating Challenge Deployer service")
	config = g.ChallengeDeployerServiceConfig
	katanaConfig = g.KatanaConfig

	if err := getClient(katanaConfig.KubeConfig); err != nil {
		return err
	}
	return Server()
}
