package rpc

import (
	"google.golang.org/grpc"
	"xmicro/internal/proto/pb"
)

func UserRpcRegister() {
	pb.RegisterUserServiceServer(grpc.NewServer(), &UserServer{})
}

func OrderRpcRegister(server *grpc.Server) {
	pb.RegisterOrderServiceServer(server, &OrderServer{})
}
