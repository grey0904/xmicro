package service

import (
	"context"
	"time"
	"xmicro/internal/app/user/dto"
	"xmicro/internal/app/user/pb"
	"xmicro/internal/common/constant"
	"xmicro/internal/common/result"
	"xmicro/internal/core/dao"
	"xmicro/internal/core/repo"
)

type UserService struct {
	userDao  *dao.UserDao
	redisDao *dao.RedisDao
}

func NewAccountService(manager *repo.Manager) *UserService {
	return &UserService{
		userDao:  dao.NewAccountDao(manager),
		redisDao: dao.NewRedisDao(manager),
	}
}

func (a *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//写注册的业务逻辑
	if req.LoginPlatform == constant.WeiXin {
		ac, err := a.wxRegister(req)
		if err != nil {
			return &pb.RegisterResponse{}, result.GrpcError(err)
		}
		return &pb.RegisterResponse{
			Uid: ac.Uid,
		}, nil
	}
	return &pb.RegisterResponse{}, nil
}

func (a *UserService) wxRegister(req *pb.RegisterRequest) (*dto.Account, *result.Error) {
	//1.封装一个account结构 将其存入数据库  mongo 分布式id objectID
	ac := &dto.Account{
		WxAccount:  req.Account,
		CreateTime: time.Now(),
	}
	//2.需要生成几个数字做为用户的唯一id  redis自增
	uid, err := a.redisDao.NextAccountId()
	if err != nil {
		return ac, result.SqlError
	}
	ac.Uid = uid
	err = a.userDao.SaveAccount(context.TODO(), ac)
	if err != nil {
		return ac, result.SqlError
	}
	return ac, nil
}
