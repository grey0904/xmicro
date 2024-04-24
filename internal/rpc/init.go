package rpc

import (
	"google.golang.org/grpc"
	"xmicro/internal/proto/pb"
)

func UserRpcRegister() {
	pb.RegisterUserServiceServer(grpc.NewServer(), &UserServer{})
}

func OrderRpcRegister() {
	pb.RegisterUserServiceServer(grpc.NewServer(), &UserServer{})
}
