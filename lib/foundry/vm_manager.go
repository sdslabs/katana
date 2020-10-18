package foundry

import (
	"context"

	pb "github.com/sdslabs/katana/lib/foundry/protos/vmmanager"
	"google.golang.org/grpc"
)

func DeployCluster(cpu, memory int32, instanceURL string) ([]byte, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewVMFoundryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.CreateCluster(ctx, &pb.ClusterRequest{
		Cpu:    cpu,
		Memory: memory,
	})
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func ClusterInfo(clusterID, instanceURL string) ([]byte, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewVMFoundryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.ClusterInfo(ctx, &pb.ClusterID{Id: clusterID})
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}
