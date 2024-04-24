package rpc

import (
	"context"
	"xmicro/internal/proto/pb"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
}

func (s *OrderServer) GetOrderByUserId(ctx context.Context, in *pb.GetOrderByUserIdRequest) (*pb.GetOrderByUserIdResponse, error) {
	return &pb.GetOrderByUserIdResponse{
		TotalAmount: "1000000",
		Status:      "1",
		OrderDate:   "20240424",
	}, nil
}
