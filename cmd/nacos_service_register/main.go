package main

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"log"
	"net"
	"xmicro/internal/nacos"
	pb "xmicro/internal/pb"
)

type imServer struct {
	pb.UnimplementedImServiceServer
}

func (s *imServer) GetChannelsData(ctx context.Context, in *pb.ChannelsDataRequest) (*pb.ChannelsDataResponse, error) {
	return &pb.ChannelsDataResponse{
		UserCount:    1,
		ChatRoomLink: "www.123.com",
	}, nil
}

func main() {
	// Register gRPC server to Nacos
	err := registerToNacos()
	if err != nil {
		log.Fatalf("Failed to register gRPC service to Nacos: %v", err)
	}

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":9971")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Create a gRPC server object
	s := grpc.NewServer()
	// Register the server with the generated protobuf code
	pb.RegisterImServiceServer(s, &imServer{})
	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	defer func() {
		err := nacos.DeregisterFromNacos()
		if err != nil {
			log.Fatalf("failed to deregister from Nacos: %v", err)
		}
	}()
}

func registerToNacos() error {
	// Create Nacos client config
	clientConfig := constant.ClientConfig{
		NamespaceId:         "ccac2447-2d93-4401-b0db-1555471bd09f",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            "nacos",
		Password:            "nacos",
	}

	// Create Nacos server config
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "43.198.104.122",
			Port:   8848,
			Scheme: "http",
		},
	}

	// Create Nacos service client
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return err
	}

	// Register gRPC service
	serviceName := "chat" // Set your service name here
	instance := vo.RegisterInstanceParam{
		Ip:          "127.0.0.1", // Set your server IP here
		Port:        9971,        // Set your server port here
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	success, err := client.RegisterInstance(instance)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to register gRPC service to Nacos")
	}

	fmt.Println("Registered gRPC service to Nacos successfully")
	return nil
}
