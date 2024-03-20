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
	"xmicro/internal/pb/im"
)

type imServer struct {
	im.UnimplementedImServiceServer
}

func (s *imServer) GetChannelsData(ctx context.Context, in *im.ChannelsDataRequest) (*im.ChannelsDataResponse, error) {
	return &im.ChannelsDataResponse{
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
	defer func() {
		if err := nacos.DeregisterFromNacos(); err != nil {
			log.Printf("Failed to deregister from Nacos: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":9971")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	im.RegisterImServiceServer(s, &imServer{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sig := <-sigChan
	log.Printf("Received signal: %v, initiating graceful shutdown", sig)
	s.GracefulStop() // Gracefully stop the server
}
