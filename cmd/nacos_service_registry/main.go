package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	err := nacos.RegistryToNacos()
	if err != nil {
		log.Fatalf("Failed to register gRPC service to Nacos: %v", err)
	}
	defer nacos.DeregisterFromNacos() // Defer deregistration

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":9971")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Register the server with the generated protobuf code
	pb.RegisterImServiceServer(s, &imServer{})

	// Prepare to catch system signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Serve in a goroutine so that it's non-blocking
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Block until a signal is received
	sig := <-sigChan
	log.Printf("Received signal: %v, initiating graceful shutdown", sig)
	s.GracefulStop() // Gracefully stop the server
}
