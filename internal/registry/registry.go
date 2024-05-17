package registry

import (
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
)

var (
	clientOnce sync.Once
	Client     naming_client.INamingClient
)

type Registrar interface {
	Register(appName string) error
	Deregister(appName string) error
	GetService(appName string) error
}

// 工厂类(抽象接口)
type RegistrarFactory interface {
	CreateRegistrar() (Registrar, error)
}
