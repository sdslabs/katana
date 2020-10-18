package vmdeployerservice

import (
	"net"

	pb "github.com/sdslabs/katana/lib/foundry/protos/vmmanager"
	"google.golang.org/grpc"
)

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
