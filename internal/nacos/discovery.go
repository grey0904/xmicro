package nacos

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"log"
	"sync"
	"sync/atomic"
	"xmicro/internal/proto/pb"
)

var (
	OrderServiceClients []pb.OrderServiceClient
	nextClientIndex     int32
	clientMu            sync.Mutex
)

func DiscoveryFromNacos(svcName string) {
	svcName = "grpc:" + svcName
	client, err := NewNamingClient()
	if err != nil {
		panic(err)
	}

	err = client.Subscribe(&vo.SubscribeParam{
		ServiceName: svcName,
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				panic(err)
			}

			updateGrpcClients(services)
		},
	})
	if err != nil {
		panic(err)
	}

	services, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: svcName,
	})
	if err != nil {
		panic(err)
	}

	updateGrpcClients(services)
}

func updateGrpcClients(services []model.Instance) {
	clientMu.Lock()
	defer clientMu.Unlock()

	OrderServiceClients = nil
	for _, service := range services {
		log.Printf("ServiceName: %s, IP: %s, Port: %d, Metadata: %s\n",
			service.ServiceName, service.Ip, service.Port, service.Metadata)

		//conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "192.168.0.151", 9998), grpc.WithInsecure()) // 本地调试
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", service.Ip, service.Port), grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to gRPC service: %v", err)
			continue
		}
		OrderServiceClients = append(OrderServiceClients, pb.NewOrderServiceClient(conn))
	}

	if len(OrderServiceClients) == 0 {
		log.Printf("No available gRPC service")
	}
}

func GetNextOrderClient() pb.OrderServiceClient {
	clientMu.Lock()
	defer clientMu.Unlock()

	if len(OrderServiceClients) == 0 {
		panic(errors.New("no available gRPC clients"))
	}

	index := atomic.AddInt32(&nextClientIndex, 1)
	return OrderServiceClients[index%int32(len(OrderServiceClients))]
}
