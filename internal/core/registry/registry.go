package registry

type InstanceInfo struct {
	ID       string
	Address  string
	Port     int
	Metadata map[string]string
}

type Registry interface {
	Register(serviceName string) error
	Deregister(serviceName string) error
	Discover(serviceName string) ([]InstanceInfo, error)
}
