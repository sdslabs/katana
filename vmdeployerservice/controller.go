package vmdeployerservice

import (
	"context"

	pb "github.com/sdslabs/katana/lib/foundry/protos/vmmanager"
)

func (s *server) ClusterInfo(ctx context.Context, clusterID *pb.ClusterID) (*pb.ClusterResponse, error) {
	var str string = "hello this is so cool"
	return &pb.ClusterResponse{Data: []byte(str)}, nil
}
