package gate

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xmicro/internal/app/gate/router"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
)

// RunV1 Nacos版：启用日志 注册路由 启动HTTP
func RunV1(ctx context.Context) error {
	// TODO 需要改为NACOS

	//1.做一个日志库 info error fatal debug
	logs.Init()
	go func() {
		//gin 启动  注册一个路由
		r := router.RegisterRouter()
		//http接口
		if err := r.Run(fmt.Sprintf(":%d", config.Conf.ServerHttp.Port)); err != nil {
			logs.Fatal("gate gin run err:%v", err)
		}
	}()
	stop := func() {
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

// RunV2 Etcd版：启用日志 注册路由 启动HTTP
func RunV2(ctx context.Context) error {
	logs.Init()

	go func() {
		//gin 启动  注册一个路由
		r := router.RegisterRouter()
		//http接口
		if err := r.Run(fmt.Sprintf(":%d", config.Conf.ServerHttp.Port)); err != nil {
			logs.Fatal("gate gin run err:%v", err)
		}
	}()

	stop := func() {
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
