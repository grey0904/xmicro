package service

import (
	"xmicro/internal/app/user/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}
