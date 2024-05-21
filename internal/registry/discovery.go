package registry

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gitlab.casinovip.tech/minigame_backend/om_struct/im_proto/im"
	"google.golang.org/grpc"
)

var (
	imServiceClients []im.ImServiceClient
	nextClientIndex  int32
	clientMu         sync.Mutex
)

// Setup initializes service discovery
func Setup() {
	client, err := NewNamingClient()
	if err != nil {
		log.Printf("Nacos NewNamingClient error: %v", err)
		return
	}

	err = client.Subscribe(&vo.SubscribeParam{
		ServiceName: "grpc:im",
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				log.Printf("nacos subscribe service error: %v", err)
				return
			}

			updateGrpcClients(services)
		},
	})
	if err != nil {
		log.Printf("Nacos Subscribe error: %v", err)
		return
	}

	err = updateServiceInfo(client)
	if err != nil {
		log.Printf("Nacos failed to get initial service info: %v", err)
		return
	}
}

func updateGrpcClients(services []model.Instance) {
	clientMu.Lock()
	defer clientMu.Unlock()

	imServiceClients = nil
	for _, service := range services {
		log.Printf("ServiceName: %s, IP: %s, Port: %d, Metadata: %s\n",
			service.ServiceName, service.Ip, service.Port, service.Metadata)

		//conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "192.168.0.151", 9998), grpc.WithInsecure()) // 本地调试
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", service.Ip, service.Port), grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to gRPC service: %v", err)
			continue
		}
		imServiceClients = append(imServiceClients, im.NewImServiceClient(conn))
	}

	if len(imServiceClients) == 0 {
		log.Printf("No available gRPC services")
	}
}

func updateServiceInfo(client naming_client.INamingClient) error {
	services, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "grpc:im",
	})
	if err != nil {
		return fmt.Errorf("nacos GetService error: %w", err)
	}

	updateGrpcClients(services)

	return nil
}

func GetNextClient() im.ImServiceClient {
	clientMu.Lock()
	defer clientMu.Unlock()

	if len(imServiceClients) == 0 {
		//panic(exception.New(exception.ImNoAvailableGrpcClients, errors.New("")))
	}

	index := atomic.AddInt32(&nextClientIndex, 1)
	return imServiceClients[index%int32(len(imServiceClients))]
}
