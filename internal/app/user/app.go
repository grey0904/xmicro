package user

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xmicro/internal/app/user/pb"
	"xmicro/internal/app/user/service"
	"xmicro/internal/common/config/center"
	"xmicro/internal/common/logs"
	"xmicro/internal/common/registry"
	myRegistry "xmicro/internal/core/registry"
	"xmicro/internal/core/repo"
	"xmicro/internal/utils/u_conv"

	"google.golang.org/grpc"
)

// RunV1 Nacos版：启动程序 启动grpc服务 启用http服务  启用日志 启用数据库
func RunV1(ctx context.Context) error {
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

// RunV2 Etcd版：启动程序 启动grpc服务 启用http服务  启用日志 启用数据库
func RunV2(ctx context.Context) error {

	logs.Init()

	// 创建资源管理器
	manager := repo.New()
	defer manager.Close()

	// 创建并注册 ETCD
	register := registry.NewRegister()
	if err := register.Register(config.Conf.Etcd); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}
	defer register.Close()

	// 准备 gRPC 服务器
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, service.NewAccountService(manager))

	// 创建监听器
	addr := config.Conf.ServerRpc.Host + ":" + u_conv.Uint64ToString(config.Conf.ServerRpc.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer lis.Close()

	// 启动 gRPC 服务
	go func() {
		logs.Info("starting gRPC server on %s", addr)
		if err := server.Serve(lis); err != nil {
			logs.Error("gRPC server error: %v", err)
		}
	}()

	// 信号处理
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)

	// 等待退出信号
	select {
	case <-ctx.Done():
		logs.Info("shutting down by context cancellation")
	case s := <-c:
		logs.Info("received signal %v", s)
	}

	// 优雅关闭
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()

	// 等待优雅关闭或超时
	select {
	case <-done:
		logs.Info("graceful shutdown completed")
	case <-time.After(10 * time.Second):
		logs.Warn("graceful shutdown timed out, forcing stop")
		server.Stop()
	}

	return nil
}
