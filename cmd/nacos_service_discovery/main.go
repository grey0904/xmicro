package main

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"log"
	"xmicro/internal/nacos"
	"xmicro/internal/pb/im"
)

func main() {
	// 创建 Nacos 客户端
	client, err := nacos.CreateNacosClient()
	if err != nil {
		log.Fatalf("Failed to create Nacos client: %v", err)
	}

	serviceInfo, err := client.GetService(vo.GetServiceParam{
		ServiceName: "grpc:im",
	})
	if err != nil {
		log.Fatalf("Failed to get service info from Nacos: %v", err)
	}

	// 打印服务信息
	fmt.Println("Service Name:", serviceInfo.Name)
	fmt.Println("Hosts:")
	for _, host := range serviceInfo.Hosts {
		fmt.Printf("  %s:%d\n", host.Ip, host.Port)

		// 创建连接
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host.Ip, host.Port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to gRPC service: %v", err)
		}
		defer conn.Close()

		// 创建 gRPC 客户端
		c := im.NewImServiceClient(conn)

		// 调用 GetChannelsData 方法
		resp, err := c.GetChannelsData(context.Background(), &im.ChannelsDataRequest{ChannelIds: []string{"channel1", "channel2"}})
		if err != nil {
			log.Fatalf("Error calling GetChannelsData: %v", err)
		}

		// 处理响应
		fmt.Println("User count:", resp.UserCount)
		fmt.Println("Chat room link:", resp.ChatRoomLink)
	}
}
