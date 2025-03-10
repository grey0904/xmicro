package user

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xmicro/internal/app/user/pb"
	"xmicro/internal/app/user/service"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
	myRegistry "xmicro/internal/core/registry"
	"xmicro/internal/core/repo"
	"xmicro/internal/utils/u_conv"

	"google.golang.org/grpc"
)

func Run(ctx context.Context) error {
	logs.Init()

	//注册 grpc service 需要数据库 mongo redis
	manager := repo.New()

	// 注册服务
	reg, err := myRegistry.InitRegistry()
	if err != nil {
		return err
	}

	if err = reg.Register(); err != nil {
		return err
	}

	//启动grpc服务端
	server := grpc.NewServer()

	go func() {
		addr := config.Conf.ServerRpc.Host + ":" + u_conv.Uint64ToString(config.Conf.ServerRpc.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			logs.Fatal("user grpc server listen err:%v", err)
		}

		pb.RegisterUserServiceServer(server, service.NewAccountService(manager))

		if err = server.Serve(lis); err != nil {
			logs.Fatal("user grpc server run failed err:%v", err)
		}
	}()
	stop := func() {
		server.Stop()
		manager.Close()
		if err := reg.Deregister(); err != nil {
			logs.Error("failed to deregister service: %v", err)
		}
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
