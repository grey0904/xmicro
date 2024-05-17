package discovery

import (
	"fmt"
	engine_client_grpc "gitlab.casinovip.tech/minigame_backend/c_engine/pkg/client/grpc"
	"gitlab.casinovip.tech/minigame_backend/c_engine/pkg/conf"
	engine_registry "gitlab.casinovip.tech/minigame_backend/c_engine/pkg/registry"
	"google.golang.org/grpc"
	"sync"
)

var (
	Client = &Clients{
		Lock:        new(sync.Mutex),
		Connections: make(map[string]*grpc.ClientConn)}
)

type Clients struct {
	Lock        *sync.Mutex
	Connections map[string]*grpc.ClientConn
}

func Init() error {

	for _, provider := range engine_registry.ProviderList() {

		var (
			tmp engine_client_grpc.Config
		)

		err := conf.UnmarshalKey(provider.ConfigKey, &tmp)

		if err != nil {
			return fmt.Errorf("grpc client config error: %w key:%s ", err, provider.ConfigKey)
		}

		cc, err := tmp.Build()
		if err != nil {
			return fmt.Errorf("grpc client build error %w key:%s ", err, provider.ConfigKey)
		}

		Client.Lock.Lock()
		defer Client.Lock.Unlock()
		Client.Connections[provider.ConfigKey] = cc
	}

	return nil
}

func (c *Clients) GetClient(service string) grpc.ClientConnInterface {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.Connections[service]
}

/*框架下需要关闭CLOSE 链接*/
func (c *Clients) Close() {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	for _, cc := range c.Connections {
		cc.Close()
	}
}
