package user

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xmicro/internal/app/user/pb"
	"xmicro/internal/app/user/service"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
	"xmicro/internal/core/repo"
	"xmicro/internal/registry"
	"xmicro/internal/utils/u_conv"
)

// RunV1 启动程序 启动grpc服务 启用http服务  启用日志 启用数据库
func RunV1(ctx context.Context) error {
	appName := config.LocalConf.AppName

	logs.Init()

	//注册 grpc service 需要数据库 mongo redis
	manager := repo.New()

	// 注册服务
	var fac registry.RegistrarFactory
	fac = new(registry.NacosFactory)
	var reg registry.Registrar
	reg, err := fac.CreateRegistrar()
	if err != nil {
		return fmt.Errorf("failed to CreateRegistrar: %v", err)
	}
	if err = reg.Register(appName); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	//启动grpc服务端
	server := grpc.NewServer()

	go func() {
		addr := config.Conf.Grpc.Host + ":" + u_conv.Uint64ToString(config.Conf.Grpc.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			logs.Fatal("user grpc server listen err:%v", err)
		}

		pb.RegisterUserServiceServer(server, &service.UserService{})

		if err = server.Serve(lis); err != nil {
			logs.Fatal("user grpc server run failed err:%v", err)
		}
	}()
	stop := func() {
		server.Stop()
		manager.Close()
		if err := reg.Deregister(appName); err != nil {
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

// RunV2 Etcd版：启动程序 启动grpc服务 启用http服务  启用日志 启用数据库
//func RunV2(ctx context.Context) error {
//	logs.Init()
//
//	//注册 grpc service 需要数据库 mongo redis
//	manager := repo.New()
//
//	//etcd注册中心 grpc服务注册到etcd中 客户端访问的时候 通过etcd获取grpc的地址
//	register := registry.NewRegister()
//
//	//启动grpc服务端
//	server := grpc.NewServer()
//
//	go func() {
//		addr := config.Conf.Grpc.Host + ":" + u_conv.Uint64ToString(config.Conf.Grpc.Port)
//		lis, err := net.Listen("tcp", addr)
//		if err != nil {
//			logs.Fatal("user grpc server listen err:%v", err)
//		}
//
//		err = register.Register(config.Conf.Etcd)
//		if err != nil {
//			logs.Fatal("user grpc server register etcd err:%v", err)
//		}
//
//		pb.RegisterUserServiceServer(server, service.NewAccountService(manager))
//
//		err = server.Serve(lis)
//		if err != nil {
//			logs.Fatal("user grpc server run failed err:%v", err)
//		}
//	}()
//	stop := func() {
//		server.Stop()
//		register.Close()
//		manager.Close()
//		//other
//		time.Sleep(3 * time.Second)
//		logs.Info("stop app finish")
//	}
//	//期望有一个优雅启停 遇到中断 退出 终止 挂断
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
//	for {
//		select {
//		case <-ctx.Done():
//			stop()
//			//time out
//			return nil
//		case s := <-c:
//			switch s {
//			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
//				stop()
//				logs.Info("user app quit")
//				return nil
//			case syscall.SIGHUP:
//				stop()
//				logs.Info("hang up!! user app quit")
//				return nil
//			default:
//				return nil
//			}
//		}
//	}
//}
