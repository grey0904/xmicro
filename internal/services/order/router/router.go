package router

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"xmicro/internal/rpc"
	"xmicro/internal/services/order/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	orderCtl := new(controller.OrderController)

	r.POST("order/list", orderCtl.List)

	return r
}

func RunServer() {
	var err error
	gin.SetMode(gin.DebugMode)

	r := SetupRouter()

	logrus.Info("Server Start Success!")

	// 启动HTTP服务器
	go func() {
		err = r.Run(":8081")
		if err != nil {
			logrus.Error("启动HTTP服务器失败: " + err.Error())
			return
		}
	}()

	// 创建一个监听指定端口的 TCP 连接
	grpcListener, err := net.Listen("tcp", ":9997")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	// 创建一个 gRPC 服务器实例
	grpcServer := grpc.NewServer()

	// 注册您的 gRPC 服务实现
	rpc.OrderRpcRegister(grpcServer)

	// 启动 gRPC 服务器开始侦听并处理请求
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			logrus.Fatalf("failed to serve: %v", err)
		}
	}()

	// 等待终止信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
