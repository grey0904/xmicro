package app

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xmicro/internal/common/config"
	"xmicro/internal/common/discovery"
	"xmicro/internal/common/logs"
	"xmicro/internal/core/repo"
	"xmicro/internal/proto/pb"
)

//var (
//	Rd     *redis.Client
//	Db     *gorm.DB
//	Config *config.AppConfig
//)

// Run 启动程序 启动grpc服务 启用http服务  启用日志 启用数据库
func Run(ctx context.Context) error {

	logs.Init()
	manager := repo.New()

	//2. etcd注册中心 grpc服务注册到etcd中 客户端访问的时候 通过etcd获取grpc的地址
	register := discovery.NewRegister()
	//启动grpc服务端
	server := grpc.NewServer()
	//注册 grpc service 需要数据库 mongo redis

	go func() {
		lis, err := net.Listen("tcp", config.Conf.Grpc.Addr)
		if err != nil {
			logs.Fatal("user grpc server listen err:%v", err)
		}
		err = register.Register(config.Conf.Etcd)
		if err != nil {
			logs.Fatal("user grpc server register etcd err:%v", err)
		}
		pb.RegisterUserServiceServer(server, service.NewAccountService(manager))
		//阻塞操作
		err = server.Serve(lis)
		if err != nil {
			logs.Fatal("user grpc server run failed err:%v", err)
		}
	}()
	stop := func() {
		server.Stop()
		register.Close()
		manager.Close()
		//other
		time.Sleep(3 * time.Second)
		logs.Info("stop app finish")
	}
	//期望有一个优雅启停 遇到中断 退出 终止 挂断
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case <-ctx.Done():
			stop()
			//time out
			return nil
		case s := <-c:
			switch s {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				stop()
				logs.Info("user app quit")
				return nil
			case syscall.SIGHUP:
				stop()
				logs.Info("hang up!! user app quit")
				return nil
			default:
				return nil
			}
		}
	}
}
