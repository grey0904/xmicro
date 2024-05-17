package rpc

import (
	"context"
	"xmicro/internal/app/user/pb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServer) GetUserInfo(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	return &pb.GetUserInfoResponse{
		UserId:   "1",
		UserName: "ben",
		Email:    "benben@gmail.com",
		Age:      "22",
	}, nil
}
