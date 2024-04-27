package service

import (
	"context"
	"google.golang.org/grpc"
	"xmicro/internal/proto/pb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
}

func OrderRpcRegister(server *grpc.Server) {
	pb.RegisterOrderServiceServer(server, &OrderService{})
}

func (s *OrderService) GetOrderByUserId(ctx context.Context, in *pb.GetOrderByUserIdRequest) (*pb.GetOrderByUserIdResponse, error) {
	return &pb.GetOrderByUserIdResponse{
		TotalAmount: "1000000",
		Status:      "1",
		OrderDate:   "20240424",
	}, nil
}
