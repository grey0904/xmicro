package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"xmicro/internal/core/repo"
	"xmicro/internal/proto/pb"
	"xmicro/internal/service/order/dto"
)

func OrderList(c *gin.Context, req dto.OrderListReq) (string, error) {

	//app.Db.Where("user_id = ?", 1).Find(&userOrders)

	return "", nil
}

type OrderService struct {
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(manager *repo.Manager) *OrderService {
	return &OrderService{}
}

func OrderRpcRegister(server *grpc.Server) {
	pb.RegisterOrderServiceServer(server, &OrderServer{})
}

func (s *OrderService) GetOrderByUserId(ctx context.Context, in *pb.GetOrderByUserIdRequest) (*pb.GetOrderByUserIdResponse, error) {
	return &pb.GetOrderByUserIdResponse{
		TotalAmount: "1000000",
		Status:      "1",
		OrderDate:   "20240424",
	}, nil
}
